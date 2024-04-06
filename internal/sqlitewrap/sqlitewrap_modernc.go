//go:build modernc
package sqlitewrap

import (
	_ "modernc.org/sqlite"
)

const Driver string = "sqlite"
