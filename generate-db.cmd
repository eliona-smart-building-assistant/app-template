@echo off
setlocal enabledelayedexpansion

go install github.com/eliona-smart-building-assistant/dev-utilities/cmd/db-generator@latest

:: Read the content of init.sql
set "INIT_SQL_CONTENT="
for /f "delims=" %%i in ('type "%CD%\db\init.sql"') do (
    set "line=%%i"
    set "INIT_SQL_CONTENT=!INIT_SQL_CONTENT!!line!!\n!"
)

:: Create init_wrapper.sql to run the script in a transaction
(
    echo BEGIN;
    echo !INIT_SQL_CONTENT!
    echo COMMIT;
) > .\db\init_wrapper.sql

:: Run PostgreSQL container
docker run -d ^
    --name "app_jet_code_generation" ^
    --platform "linux/amd64" ^
    -e "POSTGRES_PASSWORD=secret" ^
    -p "6001:5432" ^
    -v "%CD%\db\init_wrapper.sql:/docker-entrypoint-initdb.d/init_wrapper.sql" ^
    debezium/postgres:12

:: Wait for PostgreSQL to initialize
timeout /t 5 >nul

:: Run Go code generator
db-generator -dsn="postgres://postgres:secret@localhost:6001/postgres?sslmode=disable" -schema="app_schema_name" -path=".\db\generated"

:: Stop and clean up the container
docker stop "app_jet_code_generation" >nul 2>&1
docker logs "app_jet_code_generation" 2>&1 | find "ERROR" || (
    echo All good.
)
docker rm "app_jet_code_generation" >nul 2>&1

:: Clean up
del .\db\init_wrapper.sql

:: Run go mod tidy
go mod tidy
