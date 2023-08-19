package queries

import (
	"database/sql"
	"errors"
	"fmt"
)

// Query text of query.
type Query string

// Name number of query.
type Name int

// AddBotRecord query for adding record about bot.
// Query names.
const (
	_ = iota
)

var queriesSqlite3 = map[Name]Query{}

var queriesPostgres = map[Name]Query{}

// ErrNotFound occurs when query was not found.
var ErrNotFound = errors.New("the query was not found")

// ErrNilStatement occurs query statement is nil.
var ErrNilStatement = errors.New("query statement is nil")

var statements = make(map[Name]*sql.Stmt, 10)

// Prepare prepares all queries for db instance.
func Prepare(DB *sql.DB, vendor string) error {
	var queries map[Name]Query
	switch vendor {
	case "sqlite":
		queries = queriesSqlite3
	case "postgres":
		queries = queriesPostgres
	}

	for n, q := range queries {
		prep, err := DB.Prepare(string(q))
		if err != nil {
			return err
		}
		statements[n] = prep
	}
	return nil
}

// GetPreparedStatement returns *sql.Stmt by name of query.
func GetPreparedStatement(name int) (*sql.Stmt, error) {
	stmt, ok := statements[Name(name)]
	if !ok {
		return nil, ErrNotFound
	}

	if stmt == nil {
		return nil, ErrNilStatement
	}

	return stmt, nil
}

// Close closes all prepared statements.
func Close() error {
	for _, stmt := range statements {
		err := stmt.Close()
		if err != nil {
			return fmt.Errorf("error closing statement: %w", err)
		}
	}

	return nil
}
