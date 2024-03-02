package main

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/koron/funddb/internal/ammufg"
	"github.com/koron/funddb/internal/fidelity"
	_ "github.com/mattn/go-sqlite3"
)

type FundData struct {
	Scheme string
	ID     string
	Date   time.Time
	Price  int64
}

type FundDataSource interface {
	Scheme() string
	ID() string
	Date() time.Time
	Price() int64
}

func fill(src FundDataSource) FundData {
	return FundData{
		Scheme: src.Scheme(),
		ID:     src.ID(),
		Date:   src.Date(),
		Price:  src.Price(),
	}
}

func get(ctx context.Context, name string) (FundDataSource, error) {
	parts := strings.SplitN(name, ":", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid name, required format \"{scheme}:{id}\": %s", name)
	}
	scheme, id := parts[0], parts[1]
	switch scheme {
	case "fidelity":
		return fidelity.Get(ctx, id)
	case "ammufg":
		return ammufg.Get(ctx, ammufg.CodeTypeFund, id)
	default:
		return nil, fmt.Errorf("unknown scheme: %s", scheme)
	}
}

func prepareTable(ctx context.Context, db *sql.DB) error {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS fund (
	scheme TEXT,
	id TEXT,
	date TEXT,
	price INTEGER)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS fund_scheme_idx ON fund (scheme)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS fund_scheme_id_idx ON fund (scheme, id)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS fund_scheme_id_date_idx ON fund (scheme, id, date)`)
	if err != nil {
		return err
	}
	return nil
}

func insertFundData(tx *sql.Tx, d FundData) error {
	_, err := tx.Exec(`INSERT OR REPLACE INTO fund (scheme, id, date, price) VALUES (?, ?, ?, ?)`, d.Scheme, d.ID, d.Date, d.Price)
	return err
}

func run(ctx context.Context, dbfile string, items []string) error {
	if len(items) == 0 {
		return errors.New("no items")
	}
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return err
	}
	defer db.Close()

	err = prepareTable(ctx, db)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	commit := true
	defer func() {
		if !commit {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	for _, item := range items {
		d, err := get(ctx, item)
		if err != nil {
			commit = false
			return fmt.Errorf("failed to get %q: %w", item, err)
		}
		err = insertFundData(tx, fill(d))
		if err != nil {
			commit = false
			return err
		}
	}
	return nil
}

func loadItems(name string) ([]string, error) {
	var out []string
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		out = append(out, sc.Text())
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func main1() {
	dbfile := flag.String("d", "fund.db", `database file`)
	itemfile := flag.String("i", "", `items file`)
	flag.Parse()

	args := flag.Args()
	if s := *itemfile; s != "" {
		items, err := loadItems(s)
		if err != nil {
			log.Fatal(err)
		}
		if len(items) > 0 {
			args = append(args, items...)
		}
	}

	err := run(context.Background(), *dbfile, args)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	main1()
}
