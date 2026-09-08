package database

import (
	"regexp"
)

type DbInterface string

const (
	Postgres DbInterface = "postgres"
	MySQL    DbInterface = "mysql"
)

var Currentdb DbInterface = Postgres

var pgPlaceholder = regexp.MustCompile(`\$[0-9]+`)

// NormalizeQuery converts Postgres ($n) into MySQL '?'
func NormalizeQuery(d DbInterface, q string) string {
	if d != MySQL {
		return q
	}
	// Replace all $n with ?
	return pgPlaceholder.ReplaceAllString(q, "?")
}
