# HTTP Layer

Handlers in this project **do not bind, validate, or format responses**. They
take a typed request, call a service, and return a typed result:

```go
func (s *Server) login(ctx *gin.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	return s.authService.Login(req)
}
```

Everything around that line is done once, in [`handler.go`](./handler.go), by a
generic wrapper. If you are coming from NestJS, the pieces map directly:

| NestJS | Here | File |
| --- | --- | --- |
| Global `ValidationPipe` | `Handle[Req, Res]` binds and validates before your code runs | [`handler.go`](./handler.go) |
| `HttpException` | `apperror.Unauthorized("...")`, returned by a service | [`internal/apperror`](../apperror/apperror.go) |
| Global exception filter | `utils.HandleError` — the one place an error becomes a status | [`internal/utils/exception.go`](../utils/exception.go) |
| Response interceptor | `utils.DataResponse` — the success envelope | [`internal/utils/response.go`](../utils/response.go) |
| `@Body()` / `@Query()` / `@Param()` | One DTO, bound from all three sources | [`internal/dto`](../dto) |

> There is no `ShouldBindJSON` in any handler, and no `http.Status*` either.
> If you are writing one, you are working against the pattern — see [Rules](#5-rules).

---

## 1. Adding an endpoint

Four files, in this order. Example: `POST /api/v1/products`.

### Step 1 — Define the request and response DTOs

[`internal/dto/product.go`](../dto/product.go). Validation lives in the
`binding` tag; the wrapper enforces it before your handler is entered.

```go
type CreateProductRequest struct {
	CategoryId  uint    `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	SKU         string  `json:"sku" binding:"required"`
}
```

> **Never accept or return a `models.*` struct on the wire.** Models carry
> columns the client must not set (`Password`, `Role`, `DeletedAt`). DTOs are
> the contract; models are storage.

### Step 2 — Write the service method

[`internal/services`](../services). Services hold the business logic, own the
database, and **describe their failures with `apperror`** (see [section 4](#4-errors)).
They must not import `gin`.

```go
func (s *ProductService) Create(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	var category models.Category
	if err := s.db.First(&category, req.CategoryId).Error; err != nil {
		return nil, apperror.NotFound("Category not found")
	}
	// ...
}
```

### Step 3 — Write the handler

`internal/server/<area>_handlers.go`. It should be one to three lines. Anything
longer usually belongs in the service.

```go
func (s *Server) createProduct(ctx *gin.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	return s.productService.Create(req)
}
```

The signature is fixed: `func(*gin.Context, *Req) (Res, error)`. Keep `ctx` even
when unused — the wrapper passes it so you can reach the authenticated user
(see [section 6](#6-reading-the-authenticated-user)).

### Step 4 — Register the route

[`server.go`](./server.go). The success status and message live here, where
`@HttpCode()` would be in Nest:

```go
products := api.Group("/products")
products.POST("", s.authMiddleware(), Handle(http.StatusCreated, "Product created", s.createProduct))
```

If the service is new, add it to the `Server` struct and build it once in `New`
— never per request:

```go
type Server struct {
	config         *config.Config
	db             *gorm.DB
	logger         *zerolog.Logger
	authService    *services.AuthService
	productService *services.ProductService   // <- add
}
```

### That's the whole loop

```bash
curl -X POST localhost:8090/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"a@b.com","password":"secret123"}'
```

```json
{ "success": true, "message": "Login successful", "data": { "user": {...}, "access_token": "..." }, "error": "" }
```

Send it garbage and you get every problem at once, without writing a line for it:

```bash
curl -X POST localhost:8090/api/v1/auth/login \
  -H 'Content-Type: application/json' -d '{"email":"nope","password":""}'
```

```json
{
  "success": false,
  "message": "Validation failed",
  "data": null,
  "error": "",
  "errors": [
    { "field": "email", "message": "email must be a valid email address" },
    { "field": "password", "message": "password is required" }
  ]
}
```

---

## 2. What the wrapper does

`Handle(status, message, handler)` runs four steps. Your handler is step three.

```
request
   |
   |  1. bind    uri params -> query string -> body       (bindRequest)
   |  2. validate the completed struct, once              (binding.Validator)
   |         |-- invalid --> 400 + per-field messages     (utils.ValidationErrorResponse)
   |
   |  3. your handler(ctx, req) (Res, error)
   |         |-- error ----> status from the error        (utils.HandleError)
   |
   v  4. 'status' + 'message' + data                      (utils.DataResponse)
```

**One DTO can span all three sources.** Tag the fields by where they come from,
and a request like `PATCH /products/42?notify=true` binds into a single struct:

| Source | Tag | Example |
| --- | --- | --- |
| Path param | `uri` | ``ID uint `uri:"id" binding:"required"` `` |
| Query string | `form` | ``Page int `form:"page,default=1"` `` |
| JSON body | `json` | ``Name string `json:"name" binding:"required"` `` |

Validation runs **once, at the end** — after every source has been read. This is
why the intermediate binds deliberately swallow validation errors: a body field
is not "missing" just because the path params were bound first.

### Two wrappers

| Wrapper | For | Handler signature |
| --- | --- | --- |
| `Handle[Req, Res]` | Routes with input — a body, query, or path params | `func(*gin.Context, *Req) (Res, error)` |
| `HandleEmpty[Res]` | Routes with no input at all | `func(*gin.Context) (Res, error)` |

A `GET /products?page=2` still uses `Handle` — the query string is input.

---

## 3. Validation tags

Common tags. The full list is in the
[validator docs](https://pkg.go.dev/github.com/go-playground/validator/v10#hdr-Baked_In_Validators).

| Tag | Means | Message produced |
| --- | --- | --- |
| `required` | Must be present and non-zero | `email is required` |
| `email` | Valid email address | `email must be a valid email address` |
| `min=8` / `max=64` | Length (strings, slices) or value (numbers) | `password must be at least 8 characters` |
| `gt=0` / `gte=0` | Greater than / at least | `price must be greater than 0` |
| `oneof=pending paid` | One of a fixed set | `status must be one of: pending, paid` |
| `len=6` | Exact length | `otp must be exactly 6 characters` |
| `eqfield=Password` | Matches another field | `confirm_password must match Password` |
| `omitempty` | Skip the other rules when empty | — |

Messages are generated in [`internal/utils/validation.go`](../utils/validation.go).
**Add a `case` there when you start using a tag it does not know** — the fallback
(`phone failed the "e164" rule`) is correct but not something to ship to a client.

Field names in errors come from the `json` tag, then `form`, then `uri`. An
untagged field reports its Go name, which is a good signal you forgot the tag.

> ### The `required` trap
>
> `required` rejects the **zero value**, not just a missing key. So
> ``IsActive bool `binding:"required"` `` can never be `false`, and
> ``Stock int `binding:"required"` `` can never be `0`.
>
> For a field where zero is legitimate, use a pointer and check for `nil`:
>
> ```go
> IsActive *bool `json:"is_active" binding:"required"`   // false is allowed, absent is not
> Stock    int   `json:"stock" binding:"min=0"`          // 0 is allowed
> ```

---

## 4. Errors

A handler returns an error; it never picks a status code. `utils.HandleError`
maps it:

| Returned by the service | Response |
| --- | --- |
| `apperror.BadRequest(msg)` | 400 with `msg` |
| `apperror.Unauthorized(msg)` | 401 with `msg` |
| `apperror.Forbidden(msg)` | 403 with `msg` |
| `apperror.NotFound(msg)` | 404 with `msg` |
| `apperror.Conflict(msg)` | 409 with `msg` |
| `gorm.ErrRecordNotFound` (unwrapped) | 404 `Resource not found` |
| **anything else** | 500 `Internal server error` |

The rule of thumb:

- **Expected failure** — the client did something the business rules reject
  (wrong password, duplicate email, unknown id). Return an `apperror`. The
  message is shown to the client, so write it for them.
- **Unexpected failure** — the database is down, a token cannot be signed, a
  bug. Return the raw `error`. It becomes a 500, and **the detail is stripped in
  release mode** (`GIN_MODE=release`) so internals never leak. You see the full
  error while `GIN_MODE=debug`.

Need a status with context? `apperror.Wrap(status, message, err)` keeps the
cause for `errors.Is`/`errors.As` while showing only `message` to the client.

---

## 5. Rules

1. **No `ShouldBindJSON`, `ctx.JSON`, or `http.Status*` inside a handler.**
   Binding is the wrapper's job, the status is the route's or the error's.
2. **Every route goes through `Handle` or `HandleEmpty`.** A raw
   `gin.HandlerFunc` skips validation and the error filter, and its responses
   will not match the envelope. (`healthCheck` is the one deliberate exception.)
3. **Services return `apperror` for anything the client caused.** A plain error
   is a 500 — that is the contract, not an accident.
4. **DTOs on the wire, models in the database.** Map between them in the service.
5. **Services never import `gin`.** They take DTOs and return DTOs, which is what
   makes them testable without a request.
6. **Construct services once in `New`,** not per request.

---

## 6. Reading the authenticated user

[`authMiddleware`](./middlewares.go) validates the `Authorization: Bearer <token>`
header and puts the claims on the context. `adminMiddleware` additionally
requires the `admin` role, and must come after it.

```go
me := api.Group("/me", s.authMiddleware())
me.GET("", HandleEmpty(http.StatusOK, "Profile", s.profile))

admin := api.Group("/admin", s.authMiddleware(), s.adminMiddleware())
```

Read the claims in the handler through the context:

```go
func (s *Server) profile(ctx *gin.Context) (*dto.UserResponse, error) {
	userID := ctx.GetUint("user_id")     // also: user_email, user_role (strings)
	return s.userService.Profile(userID)
}
```

---

## 7. Response shapes

Every response is a [`utils.Response`](../utils/response.go). Clients can rely on
`success` alone to branch.

```jsonc
// success — status from the route
{ "success": true, "message": "Product created", "data": { ... }, "error": "" }

// apperror — status from the error
{ "success": false, "message": "Invalid email or password", "data": null, "error": "" }

// validation failure — always 400
{ "success": false, "message": "Validation failed", "data": null, "error": "",
  "errors": [ { "field": "email", "message": "email is required" } ] }

// malformed JSON — always 400, no per-field detail is possible
{ "success": false, "message": "Invalid request data", "data": null, "error": "unexpected EOF" }
```

For list endpoints there is `utils.PaginatedSuccessResponse`, which adds a
`meta` block (`page`, `limit`, `total`, `total_pages`). It **writes the response
itself**, so it cannot be used from inside a `Handle` handler — the wrapper would
then write a second body. Either carry the meta in the response DTO you return:

```go
type ProductListResponse struct {
	Items []ProductResponse    `json:"items"`
	Meta  utils.PaginationMeta `json:"meta"`
}
```

or register that one route as a plain `gin.HandlerFunc` and call
`utils.PaginatedSuccessResponse` yourself. Prefer the first — it keeps the route
inside the validation pipeline.

---

## 8. Troubleshooting

### A valid `false` or `0` is rejected as "required"

`required` rejects zero values. See [the `required` trap](#3-validation-tags).

### The error says `FirstName` instead of `first_name`

The field has no `json`/`form`/`uri` tag, or `utils.RegisterValidationTagNames()`
was not called. It runs at the top of `SetupRoutes` — keep it there.

### I get a 500 where I expect a 400/404

The service returned a plain `error`. Wrap it: `apperror.NotFound("...")`.
In `GIN_MODE=debug` the response body shows the underlying error, which usually
names the line.

### `type func(...) does not match inferred type ...` when registering a route

The handler signature is wrong. It must be exactly
`func(*gin.Context, *Req) (Res, error)` — note the **pointer** to the DTO — or
`func(*gin.Context) (Res, error)` for `HandleEmpty`.

### Query parameters arrive empty

Query fields bind from the `form` tag, not `json`. `?page=2` needs
``Page int `form:"page"` ``.

### A path param is always zero

Path fields bind from the `uri` tag, and the name must match the route
placeholder: `/products/:id` needs `uri:"id"`.

### The body is ignored on a `GET`

Bodies are only read for POST, PUT, PATCH and DELETE. Send those parameters in
the query string instead.
