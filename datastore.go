package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Datastore struct {
		pool *pgxpool.Pool
	}

	LogEntry struct {
		Date time.Time
		Addr string
		Path string
	}
)

func NewDatastore(host string, port int, username, password, dbName, dbSSLMode string) (*Datastore, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", username, password, host, port, dbName, dbSSLMode)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	ds := &Datastore{pool: pool}
	ds.init()
	return ds, nil
}

func (ds *Datastore) init() {
	if _, err := ds.pool.Exec(context.Background(), `create table if not exists logs (date timestamp without time zone, addr varchar(64), path varchar(256));`); err != nil {
		slog.Error("unable to initdb", "error", err)
	}
}

func (ds *Datastore) Add(logEntry *LogEntry) error {
	if _, err := ds.pool.Exec(context.Background(), `insert into logs (date, addr, path) values ($1, $2, $3)`, logEntry.Date, logEntry.Addr, logEntry.Path); err != nil {
		return err
	}
	slog.Info("added entry", "timestamp", logEntry.Date.Format(time.RFC3339), "path", logEntry.Path)
	return nil
}

func (ds *Datastore) GetAll() ([]*LogEntry, error) {
	rows, err := ds.pool.Query(context.Background(), `select date, addr, path from logs`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*LogEntry
	for rows.Next() {
		var entry LogEntry
		if err := rows.Scan(&entry.Date, &entry.Addr, &entry.Path); err != nil {
			return nil, err
		}
		logs = append(logs, &entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	slog.Info("retrieved entries", "count", len(logs))
	return logs, nil
}

func (ds *Datastore) Close() {
	ds.pool.Close()
}
