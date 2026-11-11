CREATE TABLE config
(
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    url              TEXT    NOT NULL,
    api_key          TEXT    NOT NULL DEFAULT '',
    enable           INTEGER NOT NULL DEFAULT 1,
    refresh_interval INTEGER NOT NULL DEFAULT 60,
    request_timeout  INTEGER NOT NULL DEFAULT 30
);
