# App template for BuildingPro Suites

Start here to build your own app for BuildingPro Suites. The app talks to the platform only through the public
[API v2](https://github.com/eliona-smart-building-assistant/eliona-api) with its own API key. It needs no access
to the platform database and keeps its own data in an embedded SQLite file.

> Looking for the old template with its own database schema (`init.sql`, `INIT_CONNECTION_STRING`)?
> It is kept as tag [`legacy-v14`](https://github.com/eliona-smart-building-assistant/app-template/tree/legacy-v14).
> New apps should not use it.

## What the template does

The example app polls a vendor API, creates one asset per vendor device and writes its temperature as data:

```
vendor API  --GET /devices-->  app  --API v2 (asset types, assets, data)-->  BuildingPro Suites
                                |
                                +-- SQLite (/data/app.db): configurations, vendor ID -> asset ID
```

- **Configuration API** (`/v1/configs`): BuildingPro Suites renders the configuration form from the
  `Configuration` schema in [openapi.yaml](openapi.yaml).
- **Asset type** from [resources/asset-types](resources/asset-types), created on start.
- **Worker** polls every enabled configuration and creates or updates assets and data.
- **fake-vendor** (`cmd/fake-vendor`) stands in for the vendor API during development.

## Quickstart

Requirements: Go 1.27, Docker.

1. Create your repository from this template (GitHub button "Use this template") and replace the placeholders,
   see [Rename the app](docs/DEVELOPMENT.md#rename-the-app).

2. Start the platform mock (API v2 on `localhost:3000`, API key `secret`):

   ```sh
   git clone https://github.com/eliona-smart-building-assistant/eliona-mock
   docker compose -f eliona-mock/docker-compose.yml up -d
   ```

3. Start the fake vendor and the app:

   ```sh
   go run ./cmd/fake-vendor &
   mkdir -p data
   API_ENDPOINT=http://localhost:3000/v2 API_TOKEN=secret API_SERVER_PORT=3001 DATA_DIR=./data \
     go run ./cmd/app-name
   ```

4. Add a configuration:

   ```sh
   curl -X POST localhost:3001/v1/configs -H 'Content-Type: application/json' \
     -d '{"url":"http://localhost:8080","refreshInterval":10,"projectIDs":["1"]}'
   ```

   Within one refresh interval the app creates three assets of type `app_name_device` in project 1.

## Environment

| Variable | Required | Default | Description |
|---|---|---|---|
| `API_ENDPOINT` | yes | | API v2 base URL, e.g. `https://<name>/api/v2` |
| `API_TOKEN` | yes | | API key of the app; it belongs to one tenant |
| `API_SERVER_PORT` | no | `3000` | Port of the app's own API |
| `DATA_DIR` | no | `/data` | Directory of the SQLite file; mount a volume here |
| `LOG_LEVEL` | no | `info` | `debug`, `info`, `warn` or `error` |

## Run as container

```sh
docker build --build-arg GIT_COMMIT=$(git rev-parse HEAD) -t app-name .
docker run -d -p 3001:3000 -v app-name-data:/data \
  -e API_ENDPOINT=https://<name>/api/v2 -e API_TOKEN=<key> app-name
```

Run exactly one instance per API key; SQLite does not support several writers.

## Documentation

- [Developer guide](docs/DEVELOPMENT.md): structure, rename, configuration schema, data, tests, release
- [API v2 reference](https://api.eliona.io/?https://raw.githubusercontent.com/eliona-smart-building-assistant/eliona-api/main/openapi.yaml)
- [Platform mock](https://github.com/eliona-smart-building-assistant/eliona-mock)
