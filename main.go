package main

import (
	"FiberFinanceAPI/api"
	"FiberFinanceAPI/database"
	db "FiberFinanceAPI/database/sqlc"
	"FiberFinanceAPI/utils"
	_ "github.com/lib/pq"
	logs "github.com/sirupsen/logrus"
	"os"
)


func init()  {
	logs.SetFormatter(&logs.TextFormatter{})
	logs.SetOutput(os.Stdout)
	logs.SetLevel(logs.DebugLevel)
}
func main() {
	config, err := utils.LoadConfig("."); if err != nil{
		logs.WithError(err).Fatal("unable to load config file")
	}
	logs.Debug("Connecting to database")

	conn, err := database.NewConnection(config); if err != nil{
		logs.WithError(err).Warn("unable to connect database")
	}

	defer func() {
		err = conn.Close(); if err != nil {
			logs.WithError(err).Warn("unable to close connection to database")
		}
	}()

	logs.WithField("version", utils.GetVersion(config)).Debug("Starting server")
	logs.Debug("Connecting to server")
	// store returns an new Repo interface that takes in our conn to ensure it implements our queries interface
	store := db.NewRepo(conn)

	server ,err := api.NewServer(config, store); if err != nil{
		logs.WithError(err).Fatal("unable to start server")
	}
	logs.Debug("Running server")

	err = server.Run(config.ServerAddress); if err != nil{
		logs.WithError(err).Fatal("unable to run server")
	}

}
