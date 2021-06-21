package database

import (
	"FiberFinanceAPI/util"
	"database/sql"
	"errors"
	"time"
)


// Connect connects to our database
func Connect(config util.Config) (*sql.DB,error) {
	conn, err := sql.Open(config.DBName, config.DBDriver); if err != nil{
		return nil, err
	}

	defer func(c *sql.DB) {
		err = c.Close()
	}(conn)

	conn.SetMaxOpenConns(12)
	// Check if database is running
	if err = waitForDB(conn, config); err != nil{
		return nil, err
	}
	return conn, err
}

// waitForDB checks is if the database is ready for connections or is up alive
func waitForDB(conn *sql.DB, config util.Config) error {
	ready := make(chan struct{})

	go func() {
		for {
			if err := conn.Ping(); err != nil{
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
