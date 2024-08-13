go install github.com/volatiletech/sqlboiler/v4@latest
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest

go get github.com/volatiletech/sqlboiler/v4
go get github.com/volatiletech/null/v8

rem Read the content of init.sql
set "INIT_SQL_CONTENT="
for /f "delims=" %%i in ('type "%cd%\db\init.sql"') do set "INIT_SQL_CONTENT=!INIT_SQL_CONTENT!%%i\n"

rem Create init_wrapper.sql to run the script in a transaction. This is needed for
rem COMMIT AND CHAIN to work in the script.
(
    echo BEGIN;
    echo %INIT_SQL_CONTENT%
    echo COMMIT;
) > %cd%\db\init_wrapper.sql

docker run -d ^
    --name "app_sql_boiler_code_generation" ^
    -e "POSTGRES_PASSWORD=secret" ^
    -p "6001:5432" ^
    -v "%cd%"\db\init_wrapper.sql:/docker-entrypoint-initdb.d/init_wrapper.sql ^
    debezium/postgres:12  > NUL

rem Wait for PostgreSQL to initialize
timeout /t 5

sqlboiler psql ^
    -c db/sqlboiler.toml ^
    --wipe --no-tests

docker stop "app_sql_boiler_code_generation" > NUL

docker logs "app_sql_boiler_code_generation" 2>&1 | findstr "ERROR" || (
    echo All good.
)

docker rm "app_sql_boiler_code_generation" > NUL

del .\db\init_wrapper.sql

go mod tidy
