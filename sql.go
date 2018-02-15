package main


import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)

// sql
//
// This holds generic sql operations

func db_connect() *sql.DB {
	// Connect to DB, return database handle
	var db_user = getenv_with_default("DWS_USER", "pound")
	var db_pass = getenv_with_default("DWS_PASS", "posterkid")
	var db_db   = getenv_with_default("DWS_DB",   "dws")
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
	defer dbq.Close()
	if err != nil {
		return ""
	}
	for dbq.Next() {
		dbq.Scan(&ret)
	}
	return ret
}
