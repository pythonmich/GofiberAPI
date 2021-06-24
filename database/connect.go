package database

import (
	"FiberFinanceAPI/utils"
	"database/sql"
	"errors"
	logs "github.com/sirupsen/logrus"
	"reflect"
	"time"
)


// connect connects to our database
func connect(config utils.Config) (*sql.DB,error) {
	logs.WithField("func", "database/connect.go -> Connect()").Info()
	logs.WithFields(logs.Fields{
		"driver_name": config.DBName,
		"DBDriver": config.DBDriver,
	}).Info()
	conn, err := sql.Open(config.DBName, config.DBDriver); if err != nil{
		logs.WithError(err).Warn("cannot connect to database")
		return nil, err
	}
	logs.Info("connected to database")

	conn.SetMaxOpenConns(12)
	// Check if database is running
	if err = waitForDB(conn, config); err != nil{
		return conn, err
	}
	return conn, err
}
func NewConnection(config utils.Config) (*sql.DB,error) {
	return connect(config)
}
// waitForDB checks is if the database is ready for connections or is up alive
func waitForDB(conn *sql.DB, config utils.Config) error {
	logs.WithField("func", "database/connect.go -> waitForDB()").Info()
	logs.WithFields(logs.Fields{
		"conn": conn == nil,
		"timeout": config.DBTimeout,
		"type": reflect.TypeOf(config.DBTimeout),
	}).Info()
	ready := make(chan struct{})

	go func() {
		for {
			if err := conn.Ping(); err != nil{
				logs.WithError(err).Warn(err.Error())
				close(ready)
				return
			}
			time.Sleep(100*time.Millisecond)
		}
	}()

	select {
	case <-ready:
		return nil
	case <- time.After(config.DBTimeout):
		return errors.New("database not ready")
	}
}
