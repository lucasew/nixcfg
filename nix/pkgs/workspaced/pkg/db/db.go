package db

import (
	"context"
	"database/sql"
	"embed"
	"os"
	"path/filepath"
	"workspaced/pkg/host"
	"workspaced/pkg/db/sqlc"
	"workspaced/pkg/types"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

type DB struct {
	*sql.DB
	Queries *sqlc.Queries
}

func Open() (*DB, error) {
	dataDir, err := host.GetUserDataDir()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dataDir, "workspaced.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	dbConn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(dbConn); err != nil {
		return nil, err
	}

	return &DB{
		DB:      dbConn,
		Queries: sqlc.New(dbConn),
	}, nil
}

func runMigrations(db *sql.DB) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}

	source, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", source, "sqlite", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (db *DB) RecordHistory(ctx context.Context, event types.HistoryEvent) error {
	return db.Queries.RecordHistory(ctx, sqlc.RecordHistoryParams{
		Command:    event.Command,
		Cwd:        event.Cwd,
		Timestamp:  event.Timestamp,
		ExitCode:   int64(event.ExitCode),
		DurationMs: event.Duration,
	})
}

func (db *DB) BatchRecordHistory(ctx context.Context, events []types.HistoryEvent) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	q := db.Queries.WithTx(tx)
	for _, event := range events {
		err := q.RecordHistory(ctx, sqlc.RecordHistoryParams{
			Command:    event.Command,
			Cwd:        event.Cwd,
			Timestamp:  event.Timestamp,
			ExitCode:   int64(event.ExitCode),
			DurationMs: event.Duration,
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (db *DB) SearchHistory(ctx context.Context, query string, limit int) ([]types.HistoryEvent, error) {
	var rows []sqlc.History
	var err error

	limit64 := int64(limit)

	if query == "" {
		rows, err = db.Queries.GetHistory(ctx, limit64)
	} else {
		rows, err = db.Queries.SearchHistory(ctx, sqlc.SearchHistoryParams{
			Command: "%" + query + "%",
			Limit:   limit64,
		})
	}

	if err != nil {
		return nil, err
	}

	events := make([]types.HistoryEvent, len(rows))
	for i, row := range rows {
		events[i] = types.HistoryEvent{
			Command:   row.Command,
			Cwd:       row.Cwd,
			Timestamp: row.Timestamp,
			ExitCode:  int(row.ExitCode),
			Duration:  row.DurationMs,
		}
	}
	return events, nil
}
