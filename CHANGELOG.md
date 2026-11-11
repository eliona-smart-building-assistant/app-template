# Changelog

All notable changes to this project will be documented in this file.

## [v2.0.0] - 2026-11-11

### Changed

- rebuild the template on API v2 only: no platform database access, own data in SQLite [@al]
- go-eliona v2, API client v3, oapi-codegen strict server on the standard library router, Go 1.27 [@al]
- assets belong to the tenant of the API key; no project selection in the configuration [@al]

### Added

- app contract: API v2 only, local data under DATA_DIR, declared egress hosts, probes on port 5000 [@al]
- metadata.json declares category external, API key role and egress hosts [@al]
- example worker that maps vendor devices to assets by global asset identifier and writes their data [@al]
- fake vendor for local development, GitHub Actions CI, developer guide [@al]

### Removed

- own database schema, `init.sql`, `INIT_CONNECTION_STRING` and Azure DevOps pipelines; the old state is tag `legacy-v14` [@al]
