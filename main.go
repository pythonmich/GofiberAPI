package main

import (
	"FiberFinanceAPI/api"
	dbConn "FiberFinanceAPI/database"
	db "FiberFinanceAPI/database/sqlc"
	"FiberFinanceAPI/utils"
	_ "github.com/lib/pq"
)

func main() {
	logs := utils.NewLogger()
	config, err := utils.LoadConfig(".")
	if err != nil {
		logs.WithError(err).Fatal("unable to load config file")
	}
	logs.Debug("Connecting to database")
	newConn := dbConn.NewConn(config, logs)

	conn, err := newConn.Connect()
	if err != nil {
		logs.WithError(err).Warn("unable to connect database")
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			logs.WithError(err).Warn("unable to close connection to database")
		}
	}()

	logs.WithField("version", utils.GetVersion(config)).Debug("Starting server")
	logs.Debug("Connecting to server")

	// store returns an new Repo interface that takes in our newConn to ensure it implements our queries interface
	store := db.NewRepo(conn, logs)

	server, err := api.NewServer(config, logs, store)
	if err != nil {
		logs.WithError(err).Fatal("unable to start server")
	}
	logs.Debug("Running server")

	err = server.Run(config.ServerAddress)
	if err != nil {
		logs.WithError(err).Fatal("unable to run server")
	}

}
