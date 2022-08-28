package bootstrap

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDB(pathDB string) (*sql.DB, error) {
	database, err := sql.Open("sqlite3", pathDB)
	if err != nil {
		return nil, fmt.Errorf("error connect database: %w", err)
	}

	if err = database.Ping(); err != nil {
		return nil, fmt.Errorf("error ping database: %w", err)
	}

	return database, nil
}
