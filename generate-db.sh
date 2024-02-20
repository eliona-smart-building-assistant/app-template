go install github.com/volatiletech/sqlboiler/v4@latest
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest

go get github.com/volatiletech/sqlboiler/v4
go get github.com/volatiletech/null/v8

docker run --rm -d \
    --name "app_sql_boiler_code_generation" \
    -e "POSTGRES_PASSWORD=secret" \
    -p "60001:5432" \
    -v "${PWD}/conf/init.sql:/docker-entrypoint-initdb.d/init.sql" \
    debezium/postgres:12

sleep 5

sqlboiler psql \
    -c sqlboiler.toml \
    --wipe --no-tests

docker stop "app_sql_boiler_code_generation"

go mod tidy
