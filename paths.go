package main

import (
	"database/sql"
	"strconv"
)

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
	return get_config_value(dbh, "blogstatic") + get_dispatch_path(dbh, "blogentry") + "entry" + zeit + ".html"
}

func path_to_archive_page(dbh *sql.DB, page int) string {
	return get_config_value(dbh, "blogstatic") + get_dispatch_path(dbh, "blogarchive") + "page" + strconv.Itoa(page) + ".html"
}
