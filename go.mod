module template

go 1.20

require (
	github.com/eliona-smart-building-assistant/app-integration-tests v1.0.1
	github.com/eliona-smart-building-assistant/go-eliona v1.9.18
	github.com/eliona-smart-building-assistant/go-utils v1.0.41
	github.com/gorilla/mux v1.8.0
	github.com/volatiletech/sqlboiler/v4 v4.15.0
	gopkg.in/yaml.v3 v3.0.1
)

// Bugfix see: https://github.com/volatiletech/sqlboiler/blob/91c4f335dd886d95b03857aceaf17507c46f9ec5/README.md
// decimal library showing errors like: pq: encode: unknown type types.NullDecimal is a result of a too-new and broken version of the github.com/ericlargergren/decimal package, use the following version in your go.mod: github.com/ericlagergren/decimal v0.0.0-20181231230500-73749d4874d5
replace github.com/ericlagergren/decimal => github.com/ericlagergren/decimal v0.0.0-20181231230500-73749d4874d5

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eliona-smart-building-assistant/go-eliona-api-client/v2 v2.5.4 // indirect
	github.com/friendsofgo/errors v0.9.2 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.2 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/pgx/v4 v4.18.1 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/volatiletech/inflect v0.0.1 // indirect
	github.com/volatiletech/strmangle v0.0.5 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
)
