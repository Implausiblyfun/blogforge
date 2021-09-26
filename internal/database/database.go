package database

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MySqlConnection consumes a config and attempts to connect to the specified db.
func MySqlConnection(config interface{}) (string, *sqlx.DB, error) {
	mappedCfg, ok := config.(map[interface{}]interface{})
	if !ok {
		return "", nil, fmt.Errorf("db setup was in unexpected form, here it is: %v of %T", config, config)
	}

	host, ok := mappedCfg["host"]
	if !ok {
		return "", nil, fmt.Errorf("db must have a host denoted in dbcfg")
	}
	dbname, ok := mappedCfg["dbname"]
	if !ok {
		return "", nil, fmt.Errorf("db must have a dbname denoted in dbcfg")
	}
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")

	// inelegant retry set up but its easy implementation timing here.
	retryCount := 3
	var db *sqlx.DB
	var err error

	if err == nil && retryCount > 0 {
		// I didnt want to spend more time on setting up a good secret formatting in this repo for now
		// so lets leak creds for now
		db, err = sqlx.Connect("mysql", fmt.Sprintf("%s:%s@(%s:3306)/%s", user, pass, host, dbname))
		if err == nil {
			retryCount = 0
		} else {
			time.Sleep(time.Second)
		}
	}
	if err != nil {
		log.Fatalln(err)
	}

	friendlyName := mappedCfg["friendly"]
	if friendlyName == "" {
		friendlyName = dbname
	}
	return fmt.Sprintf("%v", friendlyName), db, nil
}
