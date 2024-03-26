module template

go 1.21

require (
	github.com/eliona-smart-building-assistant/app-integration-tests v1.1.4
	github.com/eliona-smart-building-assistant/go-eliona v1.9.36
	github.com/eliona-smart-building-assistant/go-utils v1.1.0
	github.com/gorilla/mux v1.8.1
	github.com/volatiletech/sqlboiler/v4 v4.16.2
	gopkg.in/yaml.v3 v3.0.1
)

// Bugfix see: https://github.com/volatiletech/sqlboiler/blob/91c4f335dd886d95b03857aceaf17507c46f9ec5/README.md
// decimal library showing errors like: pq: encode: unknown type types.NullDecimal is a result of a too-new and broken version of the github.com/ericlargergren/decimal package, use the following version in your go.mod: github.com/ericlagergren/decimal v0.0.0-20181231230500-73749d4874d5
replace github.com/ericlagergren/decimal => github.com/ericlagergren/decimal v0.0.0-20181231230500-73749d4874d5

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eliona-smart-building-assistant/go-eliona-api-client/v2 v2.6.8 // indirect
	github.com/friendsofgo/errors v0.9.2 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgtype v1.14.3 // indirect
	github.com/jackc/pgx/v4 v4.18.3 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/volatiletech/inflect v0.0.1 // indirect
	github.com/volatiletech/strmangle v0.0.6 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/image v0.15.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
)
