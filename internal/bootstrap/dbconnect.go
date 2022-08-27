package bootstrap

import (
	"database/sql"
	"fmt"
)

func ConnectDB(pathDB string) (*sql.DB, error) {
	database, err := sql.Open("sqlite3", pathDB)
	if err != nil {
		return nil, fmt.Errorf("error connect database: %w", err)
	}
	return database, nil
}
