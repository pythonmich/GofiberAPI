package database

import (
	"FiberFinanceAPI/utils"
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

// connect connects to our database
func connect(config utils.Config, logs *utils.StandardLogger) (*sql.DB, error) {
	logs.WithField("func", "database/connect.go -> Connect()").Info()
	logs.WithFields(logrus.Fields{
		"driver_name": config.DBName,
		"DBDriver":    config.DBDriver,
	}).Debug()
	conn, err := sql.Open(config.DBName, config.DBDriver)
	if err != nil {
		logs.WithError(err).Warn("cannot connect to database")
		return nil, err
	}

	conn.SetMaxOpenConns(12)
	// Check if database is running
	if err = waitForDB(conn, config, logs); err != nil {
		return nil, err
	}
	logs.Info("connected to database")
	return conn, err
}
func NewConnection(config utils.Config, logs *utils.StandardLogger) (*sql.DB, error) {
	logs.WithField("func", "database/connect.go -> Connect()").Info("creating new connection")
	return connect(config, logs)
}

// waitForDB checks is if the database is ready for connections or is up alive
func waitForDB(conn *sql.DB, config utils.Config, logs *utils.StandardLogger) error {
	logs.WithField("func", "database/connect.go -> waitForDB()").Info()
	logs.WithFields(logrus.Fields{
		"conn is null": conn == nil,
		"timeout":      config.DBTimeout,
		"type":         reflect.TypeOf(config.DBTimeout),
	}).Debug()
	ready := make(chan struct{})

	go func() {
		for {
			err := conn.Ping()
			if err == nil {
				close(ready)
				return
			}
			logs.WithError(err).Fatal(err.Error())
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {
	case <-ready:
		return nil
	case <-time.After(config.DBTimeout):
		return errors.New("database not ready")
	}
}
