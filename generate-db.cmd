go install github.com/volatiletech/sqlboiler/v4@latest
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest

go get github.com/volatiletech/sqlboiler/v4
go get github.com/volatiletech/null/v8

docker run -d ^
    --name "app_sql_boiler_code_generation" ^
    -e "POSTGRES_PASSWORD=secret" ^
    -p "60001:5432" ^
    -v "%cd%"\conf\init.sql:/docker-entrypoint-initdb.d/init.sql ^
    debezium/postgres:12  > NUL

timeout /t 5

sqlboiler psql ^
    -c sqlboiler.toml ^
    --wipe --no-tests

docker stop "app_sql_boiler_code_generation" > NUL

docker logs "app_sql_boiler_code_generation" 2>&1 | findstr "ERROR" || (
    echo All good.
)

docker rm "app_sql_boiler_code_generation" > NUL

go mod tidy
