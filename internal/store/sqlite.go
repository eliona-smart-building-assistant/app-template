// Copyright (c) 2026 ABB Schweiz AG
// Authors: Adam Lange, et al.

package store

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	_ "modernc.org/sqlite" // pure Go driver, no CGO needed
)

//go:embed migrations/*.sql
var migrations embed.FS

// SQLite stores the data in one file. Run a single replica: SQLite does not support writers in several pods.
type SQLite struct {
	db *sql.DB
}

var _ Store = (*SQLite)(nil)

// OpenSQLite opens (or creates) dir/app.db and applies all pending migrations.
func OpenSQLite(ctx context.Context, dir string) (*SQLite, error) {
	dsn := "file:" + filepath.Join(dir, "app.db") + "?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	if err := migrate(ctx, db); err != nil {
		db.Close()
		return nil, err
	}
	return &SQLite{db: db}, nil
}

// migrate applies migrations/NNNN_*.sql in order, each in its own transaction, and records the version.
func migrate(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_version (version INTEGER PRIMARY KEY)`); err != nil {
		return fmt.Errorf("create schema_version: %w", err)
	}
	var current int
	if err := db.QueryRowContext(ctx, `SELECT COALESCE(MAX(version), 0) FROM schema_version`).Scan(&current); err != nil {
		return fmt.Errorf("read schema_version: %w", err)
	}

	files, err := fs.Glob(migrations, "migrations/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(files)
	for _, f := range files {
		version, err := strconv.Atoi(strings.SplitN(filepath.Base(f), "_", 2)[0])
		if err != nil {
			return fmt.Errorf("migration %s: name must start with a number", f)
		}
		if version <= current {
			continue
		}
		script, err := migrations.ReadFile(f)
		if err != nil {
			return err
		}
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, string(script)); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %s: %w", f, err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (?)`, version); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %s: %w", f, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("migration %s: %w", f, err)
		}
	}
	return nil
}

func (s *SQLite) Close() error { return s.db.Close() }

const configColumns = `id, url, api_key, enable, refresh_interval, request_timeout, project_ids`

type scanner interface{ Scan(dest ...any) error }

func scanConfig(row scanner) (Config, error) {
	var c Config
	var projects string
	if err := row.Scan(&c.ID, &c.URL, &c.APIKey, &c.Enable, &c.RefreshInterval, &c.RequestTimeout, &projects); err != nil {
		return Config{}, err
	}
	if err := json.Unmarshal([]byte(projects), &c.ProjectIDs); err != nil {
		return Config{}, fmt.Errorf("config %d: project_ids: %w", c.ID, err)
	}
	return c, nil
}

func (s *SQLite) ListConfigs(ctx context.Context) ([]Config, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+configColumns+` FROM config ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	configs := []Config{}
	for rows.Next() {
		c, err := scanConfig(rows)
		if err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}
	return configs, rows.Err()
}

func (s *SQLite) GetConfig(ctx context.Context, id int64) (Config, error) {
	c, err := scanConfig(s.db.QueryRowContext(ctx, `SELECT `+configColumns+` FROM config WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return Config{}, ErrNotFound
	}
	return c, err
}

func (s *SQLite) CreateConfig(ctx context.Context, c Config) (Config, error) {
	projects, err := json.Marshal(nonNil(c.ProjectIDs))
	if err != nil {
		return Config{}, err
	}
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO config (url, api_key, enable, refresh_interval, request_timeout, project_ids) VALUES (?, ?, ?, ?, ?, ?)`,
		c.URL, c.APIKey, c.Enable, c.RefreshInterval, c.RequestTimeout, string(projects))
	if err != nil {
		return Config{}, err
	}
	if c.ID, err = res.LastInsertId(); err != nil {
		return Config{}, err
	}
	c.ProjectIDs = nonNil(c.ProjectIDs)
	return c, nil
}

func (s *SQLite) UpdateConfig(ctx context.Context, c Config) (Config, error) {
	projects, err := json.Marshal(nonNil(c.ProjectIDs))
	if err != nil {
		return Config{}, err
	}
	res, err := s.db.ExecContext(ctx,
		`UPDATE config SET url = ?, api_key = ?, enable = ?, refresh_interval = ?, request_timeout = ?, project_ids = ? WHERE id = ?`,
		c.URL, c.APIKey, c.Enable, c.RefreshInterval, c.RequestTimeout, string(projects), c.ID)
	if err != nil {
		return Config{}, err
	}
	if n, err := res.RowsAffected(); err != nil {
		return Config{}, err
	} else if n == 0 {
		return Config{}, ErrNotFound
	}
	c.ProjectIDs = nonNil(c.ProjectIDs)
	return c, nil
}

func (s *SQLite) DeleteConfig(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM config WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrNotFound
	}
	return nil
}

func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
