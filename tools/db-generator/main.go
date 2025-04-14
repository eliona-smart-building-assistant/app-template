package main

import (
	"flag"
	"log"

	_ "github.com/lib/pq"

	"github.com/go-jet/jet/v2/generator/postgres"
)

func main() {
	dsn := flag.String("dsn", "", "Database DSN")
	schema := flag.String("schema", "", "Database schema")
	path := flag.String("path", "", "Destination directory")

	flag.Parse()

	if *dsn == "" || *schema == "" || *path == "" {
		log.Fatal("Missing required parameters. Please provide dsn, schema, and path.")
	}

	if err := postgres.GenerateDSN(*dsn, *schema, *path); err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}

	log.Println("Code generation completed successfully.")
}
