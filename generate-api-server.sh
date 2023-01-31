go install golang.org/x/tools/cmd/goimports@latest

docker run --rm \
    -v "${PWD}:/local" \
    openapitools/openapi-generator-cli:v6.2.1 generate \
    -g go-server \
    --git-user-id eliona-smart-building-assistant \
    --git-repo-id python-eliona-api-client \
    -i /local/openapi.yaml \
    -o /local/apiserver \
    --additional-properties="packageName=apiserver,sourceFolder=,outputAsLibrary=true"

docker run --rm \
    -v "${PWD}:/local" \
    openapitools/openapi-generator-cli:v6.2.1 generate \
    -g openapi \
    -i /local/openapi.yaml \
    -o /local/apiserver \
    --additional-properties=outputFile=openapi.json

goimports -w ./apiserver