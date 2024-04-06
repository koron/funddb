//go:build !modernc
package sqlitewrap

import (
	_ "github.com/mattn/go-sqlite3"
)

const Driver string = "sqlite3"
