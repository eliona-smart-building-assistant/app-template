go install github.com/eliona-smart-building-assistant/dev-utilities/cmd/db-generator@latest

# Read the content of init.sql
INIT_SQL_CONTENT=$(<"${PWD}/db/init.sql")

# Create init_wrapper.sql to run the script in a transaction. This is needed for
# COMMIT AND CHAIN to work in the script.
cat << EOF > ./db/init_wrapper.sql
BEGIN;

$INIT_SQL_CONTENT

COMMIT;
EOF

docker run -d \
    --name "app_jet_code_generation" \
    --platform "linux/amd64" \
    -e "POSTGRES_PASSWORD=secret" \
    -p "6001:5432" \
    -v "${PWD}/db/init_wrapper.sql:/docker-entrypoint-initdb.d/init_wrapper.sql" \
    debezium/postgres:12

# Wait for PostgreSQL to initialize
sleep 5

db-generator -dsn="postgres://postgres:secret@localhost:6001/postgres?sslmode=disable" -schema="app_schema_name" -path="./db/generated"

docker stop "app_jet_code_generation" > /dev/null

docker logs "app_jet_code_generation" 2>&1 | grep "ERROR" || {
    echo "All good."
}

docker rm "app_jet_code_generation" > /dev/null

rm ./db/init_wrapper.sql

go mod tidy
