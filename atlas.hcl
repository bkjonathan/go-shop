# Atlas config — generates versioned SQL migrations from the GORM models.
#
#   make db-diff name=add_wishlist   # inspect models -> write db/migrations/*.sql
#   make migrate-up                  # apply with golang-migrate
#
# `dev` is a throwaway database Atlas starts in Docker to normalise the schema.
# It is created and destroyed per command and never touches your real data, but
# Docker must be running. Pin it to the same major version as production.

variable "dev_url" {
  type    = string
  default = getenv("ATLAS_DEV_URL")
}

data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-C", "db/loader",
    ".",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = var.dev_url != "" ? var.dev_url : "docker://postgres/12/dev?search_path=public"

  migration {
    dir = "file://db/migrations"
    # Emit <version>_<name>.up.sql / .down.sql so the golang-migrate CLI you
    # already use in the Makefile can apply them unchanged.
    format = golang-migrate
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
