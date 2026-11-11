# Developer guide

## App contract

Every store app (category `external`) must keep these rules. The acceptance check tests them, and the platform
enforces the network and runtime rules for running apps.

| Area | Rule |
|---|---|
| Platform access | API v2 only, with the API key the app is installed with (`API_ENDPOINT`, `API_TOKEN`). No database: no connection string, no identity, no SQL. |
| Own data | Only under `DATA_DIR` (default `/data`, a volume of this instance). Keep it small: configuration and sync state. |
| Network | Egress only to API v2 and to the hosts listed in `metadata.json` `egress`. Everything else is blocked. |
| Tenancy | One instance per installation: the API key belongs to one tenant. |
| API | The app's own API under `/v1` on `API_SERVER_PORT` (default 3000), described by `openapi.yaml`, served as `/v1/version/openapi.json`. |
| Probes | `GET /livez` and `GET /readyz` on port 5000. Ready means the app can serve its API. |
| Runtime | Runs as non-root with a read-only file system except `DATA_DIR`; stops cleanly on SIGTERM. |
| Image | Fixed version tag (never only `latest`), no critical vulnerabilities. |
| Secrets | Fields marked `writeOnly` (like a vendor API key) are stored but never returned by the API. |

`metadata.json` declares the app:

| Field | Meaning |
|---|---|
| `category` | always `external` for store apps (connections built by ABB are `app`) |
| `elionaMinVersion` | lowest platform version the app was tested against |
| `apiKeyRole` | role of the API key the app needs: `standard`, `integration` or `read-only` |
| `egress` | host names the app may reach besides API v2 |

## Structure

```
cmd/app-name/          main: environment, store, asset types, API server, worker
cmd/fake-vendor/       vendor stand-in for local development
internal/api/          implementation of the app's API (openapi.yaml)
internal/apiserver/    generated from openapi.yaml - do not edit
internal/config/       environment variables
internal/platform/     API v2 calls (asset types, assets, data)
internal/store/        Store interface, SQLite implementation, migrations/ (configurations only)
internal/vendor/       client of the vendor API - replace with the real one
internal/worker/       polling loop: vendor -> assets and data
resources/asset-types/ asset types created on start
files.go               embeds metadata.json, openapi.yaml and resources into the binary
metadata.json          app registration: name, category, versions, API key role, egress hosts, texts
openapi.yaml           the app's API and the configuration form
```

## Rename the app

Replace these strings in all files and rename the matching paths:

| Placeholder | Replace with | Example |
|---|---|---|
| `app-name` | repository name, lower case with `-` (also `cmd/app-name/` and the Go module) | `acme-sensors` |
| `app_name` | the same with `_` (asset type name and file) | `acme_sensors` |
| `App Name` | display name | `ACME Sensors` |

```sh
grep -rl --exclude-dir=.git -e app-name -e app_name -e 'App Name' . | xargs sed -i 's/app-name/acme-sensors/g; s/app_name/acme_sensors/g; s/App Name/ACME Sensors/g'
git mv cmd/app-name cmd/acme-sensors
git mv resources/asset-types/app_name_device.json resources/asset-types/acme_sensors_device.json
```

Then update `metadata.json` (texts, `elionaMinVersion`, `documentationUrl`) and replace `icon`.

## Configuration schema

The app page in BuildingPro Suites renders its form from `components.schemas.Configuration` in `openapi.yaml`.
Use these conventions so the form works without app-specific frontend code:

| Schema | Form field |
|---|---|
| `type: string` | text |
| `format: password` + `writeOnly: true` | secret, never shown again |
| `format: uri` | URL with validation |
| `type: integer` + `default`, `minimum` | number |
| `type: boolean` | toggle |
| `enum` | select |
| `x-eliona-bind: <resource>` | picker of the tenant's resources, e.g. asset types |
| `readOnly: true` | hidden in the form |

Validate the same rules in `internal/api` as well; the form is a convenience, the API is the contract.

After changing `openapi.yaml`, regenerate the server code:

```sh
go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.5.0 -config codegen/oapi-codegen.yaml openapi.yaml
```

Path parameters must be valid Go identifiers (`configId`, not `config-id`); the server uses the standard
library router.

## Data in BuildingPro Suites

- **Asset types**: one JSON file per type in `resources/asset-types`. They are created or updated on every start,
  so changing a file and restarting updates the type. The client rejects unknown fields: attributes take a
  numeric `unitId`, not a `unit` text as in older apps.
- **Assets** belong to the tenant of the app's API key; the app chooses no tenant or project. They are
  created with a stable global asset identifier (`<asset type>_<config ID>_<vendor ID>`). The upsert
  by this identifier is idempotent, so the app needs no stored mapping: it caches the returned asset IDs in memory
  and upserts each asset once again after a restart.
- **Data**: written with subtype `input` for measured values. Use `info` for static properties, `status` for
  device state and `output` for setpoints the platform may write back.

## Own data (store)

- Add a table: create `internal/store/migrations/0002_<what>.sql`. Migrations run in file order on start, each
  in its own transaction, and are recorded in `schema_version`. Never change a released migration; add a new one.
- Add methods to the `Store` interface and implement them in `sqlite.go`. The rest of the app only uses the
  interface, so a different storage (for example a platform storage API) only needs a new implementation.
- The `/data` volume holds the configurations including vendor credentials. It is removed with the
  installation.

## Tests

```sh
go vet ./...
go test ./...
```

- `internal/store`: migrations and the store against a temporary SQLite file.
- `internal/api`: the API through HTTP, including defaults, validation and that secrets are never returned.
- `internal/worker`: one sync against an `httptest` vendor and a fake platform.

For an end-to-end run against the platform mock see the [Quickstart](../README.md#quickstart).

## Release

1. Add the changes to `CHANGELOG.md` under the new version.
2. Tag the commit on `main`: `git tag -a v1.2.3 -m "<what changed>"` and push the tag.
3. Publish the container image with that exact tag. Never publish `latest` as the only tag; installations pin a
   version.
