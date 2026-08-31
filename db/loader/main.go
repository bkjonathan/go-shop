// Command loader prints the DDL that GORM would generate for every model in
// internal/models. Atlas runs it as an "external schema" provider: the output
// is the *desired* state, which Atlas diffs against db/migrations to produce a
// new versioned migration (the equivalent of `prisma migrate dev`).
//
// Register every model here — a struct that is missing is a table Atlas will
// think you meant to drop.
package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/bkjonathan/go-shop/internal/models"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&models.User{},
		&models.RefreshToken{},
		&models.Cart{},
		&models.CartItem{},
		&models.Category{},
		&models.Product{},
		&models.ProductImage{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	if _, err := io.WriteString(os.Stdout, stmts); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write schema: %v\n", err)
		os.Exit(1)
	}
}
