package main

import (
	"FiberFinanceAPI/api"
	"FiberFinanceAPI/database"
	"FiberFinanceAPI/util"
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
	config, err := util.LoadConfig("."); if err != nil{
		logs.WithError(err).Fatal("unable to load config file")
	}

	logs.Debug("Connecting to database")

	conn, err := database.Connect(config); if err != nil{
		logs.WithError(err).Fatal("unable to close database")
	}

	logs.WithField("version", util.GetVersion(config)).Debug("Starting server")
	logs.Debug("Connecting to server")

	server ,err := api.NewServer(config, conn, logs.New()); if err != nil{
		logs.WithError(err).Fatal("unable to start server")
	}
	logs.Debug("Running server")

	err = server.Run(config.ServerAddress); if err != nil{
		logs.WithError(err).Fatal("unable to run server")
	}

}
