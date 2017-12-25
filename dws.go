package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "github.com/lib/pq"
)

func roothandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Main page<br />")
	var dbh = db_connect("pound", "posterkid", "dws")
	config_names, err := dbh.Query("SELECT id, name FROM config")
	if err != nil {
		fmt.Fprintf(w, "Error: " + err.Error() )
		return
	}
	fmt.Fprintf(w, "Got past database init<br />")
	defer config_names.Close()
	var id string
	var name string
	for config_names.Next() {
		err = config_names.Scan(&id, &name)
		fmt.Fprintf(w, "Row: " + string(id) + ": " + name + "<br />")
	}
}

func main() {
	port := getenv_with_default("DWS_PORT", "8000")
	http.HandleFunc("/", roothandler)
	http.ListenAndServe(":" + port, nil)
}

// #############
// POUNDDB

func db_connect(db_user string, db_pass string, db_db string) *sql.DB {
	// Connect to DB, return database handle
	db, err := sql.Open("postgres", "postgres://" + db_user + ":" + db_pass + "@localhost/" + db_db + "?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func get_dbresults() {
	// For the first record returned by the passed in query handle,
	// return the list as a hashmap
}

func get_dbcol() {
	// For the records returned by the passed in query handle,
	// return all their first tuple members as an array
}

func get_dbval() {
	// Given a query handle that should return at most one value,
	// return the single value directly
}

// ###############
// Init code

func getenv_with_default(key, fallback string) string {
	env_val := os.Getenv(key)
	if len(env_val) == 0 {
		return fallback
	}
	return env_val
}

// ###############
// application-database calls

func get_config_value(key string) {
	// Given a key, return its value in the config table
}

// Blog stuff
func get_blogentry(id int) {
	// Given an id in the blogentry table, return everything
	// about it needed to display it, including tags.
}

func identify_last_n_blogentries(count int, include_private bool) {
	// Returns blogentry(id) for the last up-to-$count blogentries
	// NOTE: It can return fewer than requested if there are not that many.
	// It will return the entries ordered so newer are earlier.
}

func identify_blogentries_with_tag(tag string) {
	// Returns blogentry(id) for all blogentries that have the
	// given tag
}

// Review stuff
func get_all_topics() {
	// Returns a hashmap of { review_topic => review_topic_id }
	// including all review topics
}

func get_all_targets_in_topicid(topicid int) {
	// Returns a hashmap of { review_target_name => review_target_ids}
	// for all review targets that are under the topicid
}

func identify_all_reviews_for_targetid(targetid int) {
	// Return reviewids for all reviews with the given target
}

func get_review(id int) {
	// Given an id in the review table, return everything about it
	// needed to display it.
}
