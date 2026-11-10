CREATE TABLE config
(
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    url              TEXT    NOT NULL,
    api_key          TEXT    NOT NULL DEFAULT '',
    enable           INTEGER NOT NULL DEFAULT 1,
    refresh_interval INTEGER NOT NULL DEFAULT 60,
    request_timeout  INTEGER NOT NULL DEFAULT 30,
    project_ids      TEXT    NOT NULL DEFAULT '[]' -- JSON array
);

CREATE TABLE asset
(
    config_id  INTEGER NOT NULL REFERENCES config (id) ON DELETE CASCADE,
    project_id TEXT    NOT NULL,
    vendor_id  TEXT    NOT NULL,
    asset_id   INTEGER NOT NULL,
    PRIMARY KEY (config_id, project_id, vendor_id)
);
