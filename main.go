package main

import (
	"webserver/api"
	"webserver/internal/bootstrap"
	"webserver/internal/domain/service"
	"webserver/logger"

	"webserver/internal/repository"

	"github.com/spf13/pflag"
)

type Scope struct {
	port string
	db   string
}

//go run main.go -d <путь до базы данных sqlite3> -p <порт(по дефолту 8000)> -l (логгер(по дефолту выключен))

func main() {
	var scope Scope

	logger := logger.NewLogger()

	pflag.StringVarP(&scope.port, "port", "p", "8000", "flag for port, default 8000")
	pflag.StringVarP(&scope.db, "database", "d", "", "database path")
	pflag.Parse()
	if scope.db == "" {
		logger.Fatal("empty database path")
	}

	database, err := bootstrap.ConnectDB(scope.db)
	if err != nil {
		logger.Fatal("database error:", err)
	}

	logger.Print("connected to database")

	defer func() {
		err = database.Close()
		if err != nil {
			logger.Print("error close database", err)
		}
	}()

	store := repository.NewStore(database)

	srvc := service.NewService(store)

	srv := api.NewServer(srvc, logger)
	logger.Print("starting server...")

	err = srv.Start(scope.port)
	if err != nil {
		logger.Fatal("server error:", err)
	}
}
