package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	_ "github.com/lib/pq"
)

// ####################################
// http dispatch functions are the main entrypoint into the code,
// and are executed by the webserver functions

func dispatch_root(w http.ResponseWriter, r *http.Request) {
	var collector []string
	collector = append(collector, sthtml("Main page", true, false))
	collector = append(collector, "<ul>\n")
	collector = append(collector, "\t<li><a href=\"blog\">Blog</a></li>\n")
	collector = append(collector, "\t<li><a href=\"reviews\">Reviews</a></li>\n")
	collector = append(collector, "\t<li><a href=\"/site.css\">CSS</a></li>\n")
	collector = append(collector, "</ul>\n")
	collector = append(collector, endhtml() )
	resp := strings.Join(collector, "")
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, resp)
}

func dispatch_blog_htmlview(w http.ResponseWriter, r *http.Request) {
	// FIXME: To HTML! Port the below functions from POUNDBLOGHTML
	var dbh = db_connect()
	var collector []string
	// display_blogmain() generates the framing content for a blogview, including
	// all the sidebars and topbar. If we want it to be pure-html (or at least restricted to
	//     looking up URL patterns), we need to prep the following data for it and pass it in:
	//   * all the topics the blog knows about (maybe filtered down to those that have associated entries)
	//   * the image for the blog
	//   * the calculated number of archive pages
	//   * the name of the blog owner (put this in the config table)
	//   * RSS/Atom enabledness?
	collector = append(collector, sthtml("My blog", true, false)) // FIXME
	collector = append(collector, display_blogmain("My Blog Title", "My Name", "http://127.0.0.1/cat.jpg", nil, 40, false)) // FIXME
	collector = append(collector, "<div id=\"entrypart\">\n")
	var last_ten_entries = identify_last_n_blogentries(dbh, 10, false)
	for _, entryid := range last_ten_entries {
		// display_bnode()
		entryid_i , _ := strconv.Atoi(entryid) // XXX Consider having last_ten_entries be []integer
		var blogentry = get_blogentry(dbh, entryid_i )
		collector = append(collector, display_bnode(blogentry))
	}
	collector = append(collector, "</div><!-- entrypart -->\n")
	collector = append(collector, "</div><!-- centrearea -->\n") // TODO: Make sure we're closing divs in the right order
	collector = append(collector, "<div id=\"footer\">\n")
	collector = append(collector, "Site served by DWS\n")
	collector = append(collector, "</div><!-- footer -->\n") // TODO: Make sure we're closing divs in the right order
	collector = append(collector, endhtml() ) // FIXME
	w.Header().Set("Content-Type", "text/html") // Send HTTP headers as late as possible, ideally after errors might happen
	resp := strings.Join(collector, "")
	io.WriteString(w, resp)
}

func dispatch_blog_textview(w http.ResponseWriter, r *http.Request) {
	// Saving this because it's a good template for other views of the blog data
	var dbh = db_connect()
	var resp = ""
	var last_ten_entries = identify_last_n_blogentries(dbh, 10, false)
	for _, entryid := range last_ten_entries {
		entryid_i , _ := strconv.Atoi(entryid) // XXX Consider having last_ten_entries be []integer
		var blogentry = get_blogentry(dbh, entryid_i )
		resp += "Begin blogentry\n"
		resp += "Title: " + blogentry["title"] + "\n"
		resp += "Posted: " + blogentry["zeit"] + "\n"
		resp += "Body\n-------------------\n" + blogentry["body"] + "\n--------------------\n"
		resp += "End blogentry\n\n"
	}
	w.Header().Set("Content-Type", "text/plain") // Send HTTP headers as late as possible, ideally after errors might happen
	io.WriteString(w, resp)
}

func dispatch_css(w http.ResponseWriter, r *http.Request) {
	// Send the CSS needed for DWS
	// All the CSS is in the "themedata" and "theme" tables in the database.
	// 1) The "BaseTheme" is always loaded, and if the user sends a cookie specifying
	//    an additional theme we'll load that too.
	// 2) If necessary we'll reconcile the themes, and transform meta-elements into
	//    their final form
	// 3) We send it onwards!
	// TODO: Most of step two. We haven't built the cookie-handling logic yet and are
	//       chasing MVP so we'll only load and transform the BaseTheme
	// Also note: Just for ease of debugging we're going to try to output things ordered by
	//            css type (CLASS vs ID vs TAG), then
	//            element
	// DISCLAIMER: I am not particularly proficient in CSS. The fact that this works doesn't mean
	//             it's a good place for others to learn CSS. I fumbled through it the first time
	//             with a lot of trial-and-error and am now cargo-culting my own code. Caveat emptor.
	// This is loosely based off of:
	//   https://github.com/pgunn/pound/blob/master/mod_perl/MyApache/POUND/POUNDCSS.pm
	var dbh = db_connect()
	resp := get_css(dbh, "")
	w.Header().Set("Content-Type", "text/css")
	w.WriteHeader(200)
	io.WriteString(w, resp)
}




// #############
// POUNDDB

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

// Review stuff
func get_all_topics(dbh *sql.DB) map[string]string {
	// Returns a hashmap of { review_topic => review_topic_id }
	// including all review topics
	var ret map[string]string

	dbq, err := dbh.Query("SELECT id, name FROM review_topic")
	if err != nil {
		return ret
	}
	for dbq.Next() {
		var topic string
		var id string
		dbq.Scan(&id, &topic)
		ret[topic] = id
	}
	return ret
}

func get_all_targets_in_topicid(dbh *sql.DB, topicid int) map [string]string {
	// Returns a hashmap of { review_target_name => review_target_ids}
	// for all review targets that are under the topicid
	var ret map[string]string

	dbq, err := dbh.Query("SELECT id, name FROM review_target WHERE topic=$1", topicid)
	if err != nil {
		return ret
	}
	for dbq.Next() {
		var name string
		var id string
		dbq.Scan(&id, &name)
		ret[name] = id
	}
	return ret
}

func identify_all_reviews_for_targetid(dbh *sql.DB, targetid int) []string {
	// Return reviewids for all reviews with the given target
	// TODO: Flag to filter on hidden
	var ret []string

	dbq, err := dbh.Query("SELECT id FROM review WHERE target=$1 ORDER BY zeit", targetid)
	if err != nil {
		return ret
	}
	for dbq.Next() {
		var id string
		dbq.Scan(&id)
		ret = append(ret, id)
	}
	return ret
}

func get_review(dbh *sql.DB, id int) map[string]string {
	// Given an id in the review table, return everything about it
	// needed to display it
	// SELECT * FROM review WHERE id=$id
	var ret map[string]string

	dbq, err := dbh.Query("SELECT title, zeit, body, rating FROM review WHERE id=$1", id)
	if err != nil {
		return ret
	}
	for dbq.Next() { // Should only return one tuple
		var title, zeit, body, rating string
		dbq.Scan(&title, &zeit, &body, &rating)
		ret["title"] = title
		ret["zeit"] = zeit
		ret["body"] = body
		ret["rating"] = rating
	}
	return ret
}

// Finally our main function
func main() {
	port := getenv_with_default("DWS_PORT", "8000")
	http.HandleFunc("/",		dispatch_root)
	http.HandleFunc("/blog",	dispatch_blog_htmlview)
	http.HandleFunc("/site.css",	dispatch_css)
	http.ListenAndServe(":" + port, nil)
}
