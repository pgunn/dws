package main

import (
	"database/sql"
	"fmt"
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
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Main page<br />")
	var dbh = db_connect()
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

func dispatch_blog_htmlview(w http.ResponseWriter, r *http.Request) {
	// FIXME: To HTML! Port the below functions from POUNDBLOGHTML
	var dbh = db_connect()
	var resp = ""
	// display_blogmain() + display_entrywrapper()
	var last_ten_entries = identify_last_n_blogentries(dbh, 10, false)
	for _, entryid := range last_ten_entries {
		// display_bnode()
		entryid_i , _ := strconv.Atoi(entryid) // XXX Consider having last_ten_entries be []integer
		var blogentry = get_blogentry(dbh, entryid_i )
		resp += "Begin blogentry\n"
		resp += "Title: " + blogentry["title"] + "\n"
		resp += "Posted: " + blogentry["zeit"] + "\n"
		resp += "Body\n-------------------\n" + blogentry["body"] + "\n--------------------\n"
		resp += "End blogentry\n\n"
	}
	// close_entrywrapper()
	// display_footer()
	w.Header().Set("Content-Type", "text/plain") // Send HTTP headers as late as possible, ideally after errors might happen
	fmt.Fprintf(w, resp)
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
	fmt.Fprintf(w, resp)
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
	w.Header().Set("Content-Type", "text/css")
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

func get_dbresults(res *sql.Rows) {
	// For the first record returned by the passed in query handle,
	// return the list as a hashmap
}

func get_dbcol(res *sql.Rows) []string {
	// For the records returned by the passed in query handle,
	// return all their first tuple members as an array
	// XXX Right now the database/sql package requires the cardinality
	// of the tuple to be matched by the arguments to scan. I can work around
	// this by calling ColumnTypes and building junk variables, but this prototyping
	// will just accept these undesirable conventions instead.
	// Also: see if there's a 2d-array interface to Rows.
	// Or .. or we could decide these abstractions are not helpful for this interface
	var ret []string
	for res.Next() {
		var resval string
		_ = res.Scan(&resval)
		ret = append(ret, resval)
	}
	return ret
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
	// SELECT value FROM config WHERE name=$key
}

// CSS stuff

func get_css(dbh *sql.DB, extra string) string {
	// Returns a string composed of the CSS for
	// everything DWS serves. If extra is not nil, integrates
	// the named theme into that CSS.
	var collector []string
	var ret string
	var all_css = make(map[string]map[string]map[string]string)

	dbq, err := dbh.Query("SELECT csstype, csselem, cssval FROM themedata WHERE themeid=(SELECT id FROM theme WHERE name='BaseTheme')")
	if err != nil {
		log.Print(err)
		return ret
	}
	for dbq.Next() { // Get everything we want from the database unless there's a theme
		var csstype, csselem, cssprop, cssval string
		dbq.Scan(&csstype, &csselem, &cssprop, &cssval)
		if      csstype == "CLASS" {
		} else if csstype == "ID"    {
		} else if csstype == "TAG"   {
		} else {continue} // Should probably output a warning

		all_css[csstype][csselem][cssprop] = cssval
	}
	// TODO: Code to retrieve/overlay a theme goes here
	// Next: prepare all that for display
	for csstype, _ := range all_css {
		var prefix string
		if csstype == "CLASS" { prefix = "." // It's like the Go designers did a survey to find
		} else if csstype == "ID" { prefix = "#" // the ugliest way to express this
		} else if csstype == "TAG" {prefix = ""}
		for csselem, _ := range all_css[csstype] {
			collector = append(collector, prefix + csselem + "\n{\n")
			for cssprop, _ := range all_css[csstype][csselem] {
				var content = all_css[csstype][csselem][cssprop]
				if cssprop == "E" || cssprop == "B" { // These are for css properties that need verbatim handling
					if content != "" {
						collector = append(collector, content + "\n")
					}
				} else { // Most properties are key-value and fit nicely into this form instead
						collector = append(collector, cssprop + ": " + content + ";\n")
				}
			}
			collector = append(collector, "}\n\n")
		}
	}
	ret = strings.Join(collector, "")
	return ret
}

// Blog stuff
func get_blogentry(dbh *sql.DB, id int) map[string]string {
	// Given an id in the blogentry table, return everything
	// about it needed to display it, including tags.
	// SELECT * FROM blogentry WHERE id=$id
	// and also
	// SELECT id as tagid, tagname FROM tag WHERE tagid IN (
	//	SELECT tagid FROM blogentry_tags WHERE beid=$id)
	var mymap = make(map[string]string)

	dbq, err := dbh.Query("SELECT title, zeit, body FROM blogentry WHERE id=$1", id)
	if err != nil {
		mymap["title"] = "Fake title for blogentry " + strconv.Itoa(id)
		mymap["zeit"] = "1514238175"
		mymap["body"] = "Fake blogentry body"
		log.Print(err)
		return mymap
	}
	for dbq.Next() {
		var title, zeit, body string
		dbq.Scan(&title, &zeit, &body)
		mymap["title"] = title
		mymap["zeit"] = zeit
		mymap["body"] = body
	}
	return mymap
}

func identify_last_n_blogentries(dbh *sql.DB, count int, include_private bool) []string {
	// Returns blogentry(id) for the last up-to-$count blogentries
	// NOTE: It can return fewer than requested if there are not that many.
	// It will return the entries ordered so newer are earlier.
	// SELECT id FROM blogentry WHERE hidden=$include_private
	// 	ORDER BY zeit DESC LIMIT $count
	var ret []string
	dbq, err := dbh.Query("SELECT id FROM blogentry WHERE private=$1 ORDER BY zeit LIMIT $2", include_private, count)
	if err != nil {
		ret = append(ret, "1")
		ret = append(ret, "2")
		log.Print(err)
		return ret
	}
	for dbq.Next() {
		var retval string
		err = dbq.Scan(&retval)
		ret = append(ret, retval)
	}
	return ret
}

func tagid_for_tag(tag string) {
	// Returns id for named tag
	// SELECT id FROM tag WHERE tagname=$tag
}

func identify_blogentries_with_tag(tagid int) {
	// Returns blogentry(id) for all blogentries that have the
	// given tag
	// SELECT id FROM blogentry WHERE id IN (
	//	SELECT beid FROM blogentry_tags WHERE tagid=$tagid)
}

// Review stuff
func get_all_topics() {
	// Returns a hashmap of { review_topic => review_topic_id }
	// including all review topics
	// SELECT id, name FROM review_topic
}

func get_all_targets_in_topicid(topicid int) {
	// Returns a hashmap of { review_target_name => review_target_ids}
	// for all review targets that are under the topicid
	// SELECT id, name FROM review_target WHERE topic=$topicid
}

func identify_all_reviews_for_targetid(targetid int) {
	// Return reviewids for all reviews with the given target
	// SELECT id FROM review WHERE target=$targetid
}

func get_review(id int) {
	// Given an id in the review table, return everything about it
	// needed to display it
	// SELECT * FROM review WHERE id=$id
	// and also
	// SELECT name FROM review_target WHERE id=(target from previous query)
	// and possibly the same for review_topic too, not sure.
}

// Finally our main function
func main() {
	port := getenv_with_default("DWS_PORT", "8000")
	http.HandleFunc("/",		dispatch_root)
	http.HandleFunc("/blog",	dispatch_blog_htmlview)
	http.HandleFunc("/site.css",	dispatch_css)
	http.ListenAndServe(":" + port, nil)
}
