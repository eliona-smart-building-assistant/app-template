package main

import (
	"flag"
	"log"

	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/go-jet/jet/v2/generator/template"
	postgres2 "github.com/go-jet/jet/v2/postgres"
	"github.com/lib/pq"
)

func main() {
	dsn := flag.String("dsn", "", "Database DSN")
	schema := flag.String("schema", "", "Database schema")
	path := flag.String("path", "", "Destination directory")

	flag.Parse()

	if *dsn == "" || *schema == "" || *path == "" {
		log.Fatal("Missing required parameters. Please provide dsn, schema, and path.")
	}

	err := postgres.GenerateDSN(
		*dsn,
		*schema,
		*path,
		template.Default(postgres2.Dialect).
			UseSchema(func(schema metadata.Schema) template.Schema {
				return template.DefaultSchema(schema).
					UseModel(template.DefaultModel().
						UseTable(func(table metadata.Table) template.TableModel {
							return template.DefaultTableModel(table).
								UseField(func(column metadata.Column) template.TableModelField {
									defaultTableModelField := template.DefaultTableModelField(column)

									// Check if the column is of type ARRAY and its element type is text or varchar
									if column.DataType.Kind == metadata.ArrayType && column.DataType.Name == "text[]" {
										defaultTableModelField.Type = template.NewType(pq.StringArray{})
									}

									return defaultTableModelField
								})
						}),
					)
			}),
	)

	if err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}

	log.Println("Code generation completed successfully.")
}
