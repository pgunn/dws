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
	if err != nil {
		return ""
	}
	for dbq.Next() {
		dbq.Scan(&ret)
	}
	return ret
}

// Path stuff. Maybe should go in its own file. Also, we should do more of these.

func get_dispatch_path(dbh *sql.DB, feature string) string {
	// convenience wrapper around get_config_value that reasons about
	// paths to various dispatch paths. Later extend this to add a
	// base prefix. Maybe do some error handling too.
	// XXX: Any failures here should abort a request, not kill the server.
	path := get_config_value(dbh, "path_" + feature)
	return path
}

func path_to_blogentry(dbh *sql.DB, zeit string) string {
	return get_dispatch_path(dbh, "blogentry") + "entry" + zeit + ".html"
}
