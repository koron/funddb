package dataobj

import (
	"fmt"

	"github.com/koron/funddb/internal/sqlitewrap"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

var initStatements []string = []string{
	`CREATE TABLE IF NOT EXISTS funds (
		id       TEXT PRIMARY KEY NOT NULL,
		name     TEXT NOT NULL,
		url      TEXT NOT NULL,
		fetch_id TEXT NULL)`,
	`CREATE UNIQUE INDEX IF NOT EXISTS UQE_funds_name ON funds (name)`,
	`CREATE UNIQUE INDEX IF NOT EXISTS UQE_funds_url ON funds (url)`,
	`CREATE UNIQUE INDEX IF NOT EXISTS UQE_funds_fetch_id ON funds (fetch_id)`,

	`CREATE TABLE IF NOT EXISTS prices (
		id    TEXT    NOT NULL,
		date  TEXT    NOT NULL,
		value INTEGER NOT NULL,
		net_assets INTEGER NULL,
		PRIMARY KEY (id, date),
		FOREIGN KEY (id) REFERENCES funds (id) ON DELETE CASCADE)`,
	`CREATE INDEX IF NOT EXISTS IDX_prices_id ON prices (id)`,
	`CREATE INDEX IF NOT EXISTS IDX_prices_date ON prices (date)`,
	`CREATE UNIQUE INDEX IF NOT EXISTS UQE_prices_id_date ON prices (id, date)`,
}

func NewEngine(dbname string) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine(sqlitewrap.Driver, dbname)
	if err != nil {
		return nil, err
	}
	engine.SetColumnMapper(names.GonicMapper{})
	return engine, nil
}

func InitSchema(engine *xorm.Engine, verbose bool) error {
	if verbose {
		engine.ShowSQL(true)
		defer engine.ShowSQL(false)
	}
	engine.SetColumnMapper(names.GonicMapper{})
	for i, s := range initStatements {
		_, err := engine.Exec(s)
		if err != nil {
			return fmt.Errorf("init statement #%d failed: %w", i, err)
		}
	}
	return nil
}
