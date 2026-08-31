# Database Migrations

Schema changes are **generated from the GORM models**, not written by hand.

You edit a struct in [`internal/models`](../internal/models), run one command, and
the SQL for the change is written to [`db/migrations`](./migrations). This is the
same workflow as `prisma migrate dev` or `drizzle-kit generate` in Node.

Two tools do the work:

| Tool | Role |
| --- | --- |
| [Atlas](https://atlasgo.io) | Reads the models, diffs them against the migration history, **writes** the SQL |
| [golang-migrate](https://github.com/golang-migrate/migrate) | **Applies** that SQL to the database |

> `AutoMigrate` is not used anywhere in this project. Every schema change is a
> reviewable SQL file in git.

---

## 1. One-time setup

You need Go, Docker, and the two CLIs above.

```bash
# Atlas — generates the migrations
brew install ariga/tap/atlas

# golang-migrate — applies them (with the Postgres driver)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

> `go install` puts binaries in `$(go env GOPATH)/bin`, which is **not** on the
> default macOS PATH. The Makefile adds it automatically, so `make` targets work
> either way. To call `migrate` directly, add this to your `~/.zshrc`:
> ```bash
> export PATH="$PATH:$(go env GOPATH)/bin"
> ```

Start Postgres and bring the schema up to date:

```bash
make docker-up     # postgres + localstack
make migrate-up    # apply all pending migrations
```

Verify:

```bash
make migrate-status   # prints the current version, e.g. 20260831171939
```

---

## 2. Changing the schema

### Step 1 — Edit the model

```go
// internal/models/product.go
type Product struct {
	BaseEntity

	CategoryID  uint    `json:"category_id" gorm:"not null"`
	Name        string  `json:"name" gorm:"not null"`
	SKU         string  `json:"sku" gorm:"uniqueIndex;not null"`  // <- new column
	Price       float64 `json:"price" gorm:"not null"`
}
```

### Step 2 — Register it, **if you added a new model**

New *field* on an existing struct? Skip this step.

New *struct*? Add it to [`db/loader/main.go`](./loader/main.go):

```go
stmts, err := gormschema.New("postgres").Load(
	&models.User{},
	&models.Product{},
	&models.Wishlist{},   // <- add it here
)
```

> **This is the easiest thing to get wrong.** Atlas treats that list as the
> complete definition of the schema. A model missing from it is not "skipped" —
> it is read as *"this table should not exist"*, and Atlas will generate a
> `DROP TABLE` for it.

### Step 3 — Generate the SQL

```bash
make db-diff name=add_product_sku
```

Name it after the change, in `snake_case` — it becomes part of the filename.
Two files appear in [`db/migrations`](./migrations):

```
20260901093012_add_product_sku.up.sql     # forward
20260901093012_add_product_sku.down.sql   # rollback
```

Atlas writes only the **delta**, never the whole schema:

```sql
-- 20260901093012_add_product_sku.up.sql
ALTER TABLE "products" ADD COLUMN "sku" text NOT NULL;
CREATE UNIQUE INDEX "idx_products_sku" ON "products" ("sku");
```

If nothing changed, it prints `The migration directory is synced with the desired
state` and writes nothing. Running it repeatedly is safe.

### Step 4 — Read the SQL before you apply it

This is not a formality. Generated SQL is correct but not always *safe* on a
table that already holds rows. Look for:

- **`ADD COLUMN ... NOT NULL` with no default** — fails outright if the table has rows.
- **`CREATE UNIQUE INDEX`** — fails if existing rows would collide.
- **`DROP COLUMN` / `DROP TABLE`** — silent data loss. If you did not intend it,
  something is wrong: usually a model missing from Step 2.
- **A column rename** — Atlas sees a drop plus an add, not a rename. Edit the
  file by hand into `ALTER TABLE ... RENAME COLUMN ...`, then run `make db-hash`
  (see [Editing a generated migration](#editing-a-generated-migration)).

You may edit these files. They are yours once generated.

### Step 5 — Apply

```bash
make migrate-up
```

### Step 6 — Commit everything together

```bash
git add internal/models/ db/migrations/ db/loader/
```

Both `.sql` files **and** `atlas.sum` must be committed. A migration without its
checksum entry breaks `db-diff` for everyone else.

---

## 3. Command reference

| Command | Does |
| --- | --- |
| `make db-diff name=xxx` | Generate a migration from model changes |
| `make db-inspect` | Print the DDL the models currently describe (no files written) |
| `make db-hash` | Recompute `atlas.sum` after hand-editing a migration |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Roll back **one** migration |
| `make migrate-reset` | Roll back **every** migration (drops all tables) |
| `make migrate-status` | Print the currently applied version |

Point any of them at another database with `DB_URL`:

```bash
make migrate-up DB_URL="postgresql://user:pass@staging-host:5432/shop?sslmode=require"
```

---

## 4. How it works

```
internal/models/*.go
        |
        |  db/loader  (a separate Go module)
        v
   desired DDL  ------>  atlas migrate diff  ------>  db/migrations/*.sql
                              ^                              |
                              |                              | golang-migrate
                       db/migrations                         v
                     (current state)                     Postgres
```

**Atlas needs a scratch database.** To diff reliably it replays the migration
directory into a throwaway Postgres that it starts in Docker and destroys
immediately after. This is why Docker must be running for `db-diff`. **Your real
database is never touched by `db-diff`** — only `migrate-up` writes to it.

The scratch database is pinned to Postgres 12 in [`atlas.hcl`](../atlas.hcl) to
match `docker-compose.yml`. Keep those in sync, or the generated SQL may not
match what production accepts. Override for a one-off:

```bash
ATLAS_DEV_URL="docker://postgres/16/dev?search_path=public" make db-diff name=xxx
```

**`db/loader` is its own Go module** ([`go.mod`](./loader/go.mod)) on purpose.
`atlas-provider-gorm` pulls in the MySQL and SQL Server drivers plus the Azure
and GCP SDKs. Isolating it keeps those out of the application's dependency graph
and out of the built binary. Nothing in `cmd/` or `internal/` imports it.

---

## 5. Rules

1. **Never edit a migration that has been applied anywhere but your laptop.**
   Once a file is merged, it is history. Fix a mistake with a *new* migration.
2. **Never hand-write a migration for a change the models can express.** Generate
   it, so the models and the schema stay the single source of truth. Hand-written
   SQL is for what GORM cannot describe: data backfills, triggers, custom types.
3. **Review the generated SQL every time.** See Step 4.
4. **Commit `atlas.sum` with the `.sql` files.**

---

## 6. Troubleshooting

### `checksum mismatch` / `atlas.sum` error

You edited, added, or deleted a migration file by hand. Recompute:

```bash
make db-hash
```

### `Dirty database version N. Fix and force version.`

A migration failed halfway. golang-migrate refuses to continue until you resolve
it. Inspect the database, finish or undo the partial change by hand, then tell
golang-migrate which version you are actually on:

```bash
migrate -path db/migrations -database "$DB_URL" force N
```

Use `N-1` if you fully undid it, `N` if you completed it by hand. Then re-run
`make migrate-up`.

### `db-diff` wants to DROP a table you still use

A model is missing from [`db/loader/main.go`](./loader/main.go). See Step 2.
Delete the generated files, fix the loader, run `make db-hash`, and diff again.

### `Cannot connect to the Docker daemon`

Docker Desktop is not running. `db-diff` needs it; `migrate-up` does not.

### `migrate: command not found`

`$(go env GOPATH)/bin` is not on your PATH. See [One-time setup](#1-one-time-setup).

### Throwing away a migration you generated but have not applied

```bash
rm db/migrations/<version>_<name>.{up,down}.sql
make db-hash
```

If you already applied it, run `make migrate-down` **first**.

### Starting over locally

```bash
make migrate-reset && make migrate-up
```

### `Atlas Pro can lint your new migration...`

An ad, printed after every successful diff. Ignore it. Everything this project
uses is free and no account is required.

---

## 7. Editing a generated migration

Atlas checksums the directory. After changing any `.sql` file by hand:

```bash
make db-hash
```

Skip it and the next `db-diff` fails with a checksum mismatch.

Note that hand-editing changes only what the *database* gets — the models still
define the desired state. If your edit makes the schema diverge from the structs,
the next `db-diff` will try to "correct" the difference. Keep hand edits to
things the models cannot express (backfills, triggers), or change the models to
match.
