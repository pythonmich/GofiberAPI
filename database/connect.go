package database

import (
	"FiberFinanceAPI/utils"
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

type ConnInter interface {
	Connect() (*sql.DB, error)
	waitForDB() error
}

type connect struct {
	conn   *sql.DB
	config utils.Config
	logs   *utils.StandardLogger
}

// NewConn creates a new connect instance
func NewConn(config utils.Config, logs *utils.StandardLogger) ConnInter {
	logs.WithField("func", "database/connect.go -> NewConn()").Debug("creating new connect")
	return &connect{
		config: config,
		logs:   logs,
	}
}

// Connect connects to our database
func (c *connect) Connect() (*sql.DB, error) {
	c.logs.WithField("func", "database/connect.go -> Connect()").Debug()
	c.logs.WithFields(logrus.Fields{
		"driver_name": c.config.DBName,
		"DBDriver":    c.config.DBDriver,
	}).Debug()
	conn, err := sql.Open(c.config.DBName, c.config.DBDriver)
	if err != nil {
		c.logs.WithError(err).Warn("cannot connect to database")
		return nil, err
	}
	c.conn = conn
	c.conn.SetMaxOpenConns(12)
	// Check if database is running and ready for connect
	if err = c.waitForDB(); err != nil {
		return nil, err
	}
	c.logs.Info("connected to database")
	return c.conn, err
}

// waitForDB checks is if the database is ready for connections or is up alive
func (c *connect) waitForDB() error {
	c.logs.WithField("func", "database/connect.go -> waitForDB()").Debug()
	c.logs.WithFields(logrus.Fields{
		"connect is null": c.conn == nil,
		"timeout":         c.config.DBTimeout,
		"type":            reflect.TypeOf(c.config.DBTimeout),
	}).Debug()
	ready := make(chan struct{})

	go func() {
		for {
			err := c.conn.Ping()
			if err == nil {
				close(ready)
				return
			}
			c.logs.WithError(err).Fatal(err.Error())
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {
	case <-ready:
		return nil
	case <-time.After(c.config.DBTimeout):
		return errors.New("database not ready")
	}
}
