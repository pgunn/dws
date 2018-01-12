package main


import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)

// sql
//
// This package holds generic sql operations

func db_connect() *sql.DB {
	// Connect to DB, return database handle
	var db_user = "pound" // TODO: Use env vars
	var db_pass = "posterkid"
	var db_db   =  "dws"
	db, err := sql.Open("postgres", "postgres://" + db_user + ":" + db_pass + "@localhost/" + db_db + "?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}


func get_config_value(dbh *sql.DB, key string) string {
	// Given a key, return its value in the config table
	ret := ""
	dbq, err := dbh.Query("SELECT value FROM config WHERE name=$1", key)
	if err != nil {
		return ""
	}
	for dbq.Next() {
		dbq.Scan(&ret)
	}
	return ret
}
