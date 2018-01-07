package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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
	// display_entrywrapper()
	var last_ten_entries = identify_last_n_blogentries(dbh, 10, false)
	for _, entryid := range last_ten_entries {
		// display_bnode()
		entryid_i , _ := strconv.Atoi(entryid) // XXX Consider having last_ten_entries be []integer
		var blogentry = get_blogentry(dbh, entryid_i )
		collector = append(collector, display_bnode(blogentry))
	}
	collector = append(collector, "</div><!-- entrypart -->\n")
	collector = append(collector, "</div><!-- centrearea -->\n") // TODO: Make sure we're closing divs in the right order
	// display_footer()
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
// POUNDBLOGHTML/POUNDHTML
// These functions should not make database calls, and should just be simple
// string manipulation (or at least restricted to looking up global settings in the config
// table and URL-path stuff)
// TODO: Consider renaming these
func display_blogmain(title string, owner string, blogimg string, topics []string, num_archives int, do_feeds bool) string {
	var collector []string
	caption_extra := "A blog by " + owner

	collector = append(collector, "<div id=\"toparea\">\n")
	collector = append(collector, "<div id=\"caption\">\n")
	collector = append(collector, "\t<div id=\"picarea\">\n")
	collector = append(collector, "\t\t<img src=\"" + blogimg + "\" />\n")
	collector = append(collector, "\t</div><!-- picarea -->\n")
	collector = append(collector, "\t<div id=\"picareatext\">\n")
	collector = append(collector, "\t\t<h1>" + title + "</h1>\n")
	collector = append(collector, "\t\t<h1>" + caption_extra + "</h1>")
	collector = append(collector, "\t</div><!-- picareatext -->\n")
	collector = append(collector, "</div><!-- caption -->\n")
	collector = append(collector, "</div><!-- toparea -->\n")

	// TODO tmentry div and its contents (topics)
	collector = append(collector, "<div id=\"centrearea\">\n") // We leave this open
	collector = append(collector, "\t<div id=\"menupart\">\n")
	collector = append(collector, "\t\t<div id=\"archmenu\" class=\"gmenu\">\n")
	collector = append(collector, "\t\t\tArchives\n")
	// TODO archives
	collector = append(collector, "\t\t</div><!-- archmenu -->\n")
	collector = append(collector, "\t\t<div id=\"topicmenu\" class=\"gmenu\">\n")
	collector = append(collector, "\t\t\tTopics\n")
	// TODO Topics
	collector = append(collector, "\t\t</div><!-- topicmenu -->")
	collector = append(collector, "\t\t<br />\n")
	collector = append(collector, "\t</div><!-- menupart -->\n")

	var ret = strings.Join(collector, "")
	return ret
}

func display_entrywrapper() string {
	return ""
}

func display_bnode(bentrydata map[string]string) string {
	// Return HTML for a single blog entry. Called both for single-entry view
	// as well as showing a bunch on a page. This code is a lot less abstract than
	// the equivalent POUND code, since this blog engine doesn't do nearly as much.
	// title zeit body
	//if val, present := bentrydata["tags"]; present {
		// TODO: Old code turned this into a list of links to topic/tag pages
	//}

	// TODO: Footer section (also passed to bnode) with a link to just this entry

	// Render the markup language.
	var content = do_markup(bentrydata["body"], "blogentryv1")
	var ret = draw_bnode(bentrydata, content)
	// draw_bnode() consumed the output of that and actually spat out the code.
	// in draw_bnode()
	return ret
}

func close_entrywrapper() string {
	return ""
}

func display_footer() string {
	return ""
}

func draw_bnode(bentrydata map[string]string, content string) string {
	// We'll have to extend the signature for both topics and the footer section
	var collector []string

	collector = append(collector, "<div class=\"jentry\">\n")
	collector = append(collector, "\t<div class=\"jehead\">\n")
	collector = append(collector, "\t<div class=\"jetitle\">" + bentrydata["title"] + "</div>\n")
	// Time
	var zeit_int, _ = strconv.ParseInt(bentrydata["zeit"], 10, 64) // base 10, 64-bit output
	var timestring = time.Unix(zeit_int, 0).Format(time.RFC3339)
	collector = append(collector, "\t<div class=\"jeheadtime\">")
	collector = append(collector, "<div class=\"jeheadtimet\">") // I don't remember why we had two divs here
	collector = append(collector, "<div class=\"jeheadtimetext\">" + timestring + "</div>")
	collector = append(collector, "</div><!-- jeheadtimet -->")
	collector = append(collector, "</div><!-- jeheadtime -->")
	// TODO Code for jemisc, the extensible area for extra tabular data like music
	collector = append(collector, "</div><!-- jehead -->\n")
	collector = append(collector, "\t<div class=\"jbody\">\n")
	collector = append(collector, "<p>" + content + "</p>\n")
	collector = append(collector, "\t</div><!-- jbody -->\n")
	collector = append(collector, "\t<div class=\"jetail\">\n")
	// TODO Tail code here
	collector = append(collector, "\t</div><!-- jetail -->\n")
	collector = append(collector, "</div><!-- jentry -->")
	collector = append(collector, "<br /><br />\n\n")
	var ret = strings.Join(collector, "")
	return ret
}

func sthtml(title string, public bool, do_feeds bool) string {
	var collector []string

	if !public {
		collector = append(collector, "<META name=\"ROBOTS\" context=\"NOINDEX\">\n")
	}
	collector = append(collector, "<!DOCTYPE html PUBLIC \"-//W3C/DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n")
	collector = append(collector, "<html xmlns=\"http://www.w3.org/1999/xhtml\">\n")
	collector = append(collector, "<head>\n")
	collector = append(collector, "<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" />\n")
	collector = append(collector, "<style type=\"text/css\">\n")
	collector = append(collector, "@import url(\"/site.css\");\n") // make this a configurable path?
	collector = append(collector, "</style>\n")
	collector = append(collector, "<title>" + title + "</title>\n")
	if do_feeds {
		var rss_url  = ""
		var atom_url = ""
		collector = append(collector, "<link rel=\"alternate\" type=\"application/rss+xml\" title=\"RSS\" href=\"" + rss_url +  "\" />\n")
		collector = append(collector, "<link rel=\"alternate\" type=\"application/atom+xml\" title=\"Atom\" href=\"" + atom_url +  "\" />\n")
	}
	collector = append(collector, "</head>\n")
	collector = append(collector, "<body>\n")
	var ret = strings.Join(collector, "")
	return ret
}

func endhtml() string {
	return "</body></html>"
}

// #############
// POUNDMarkup
// This code implements the markup language I used in my classic blog. It's similar to
// Mediawiki's markup. I liked it, so we implement something similar here

func do_markup(data string, rendermode string) string {
	// FIXME
	return data
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
	// everything DWS serves. If extra is not empty, integrates
	// the named theme into that CSS.
	var collector []string
	var ret string
	var all_css = make(map[string]map[string]map[string]string)

	dbq, err := dbh.Query("SELECT csstype, csselem, cssprop, cssval FROM themedata WHERE themeid=(SELECT id FROM theme WHERE name='BaseTheme')")
	if err != nil {
		log.Print(err)
		return ret
	}
	for dbq.Next() { // Get everything we want from the database unless there's a theme
		var csstype, csselem, cssprop, cssval string
		dbq.Scan(&csstype, &csselem, &cssprop, &cssval)
		if        csstype == "CLASS" {
		} else if csstype == "ID"    {
		} else if csstype == "TAG"   {
		} else {continue} // Should probably output a warning
		// Need to fill in bits of the chained structure that are not present.
		if _, present := all_css[csstype]; !present {
			all_css[csstype] = make(map[string]map[string]string)
		}
		if _, present := all_css[csstype][csselem]; !present {
			all_css[csstype][csselem] = make(map[string]string)
		}
		all_css[csstype][csselem][cssprop] = cssval
	}
	// TODO: Code to retrieve/overlay a theme goes here
	// Next: prepare all that for display
	for csstype, _ := range all_css {
		var prefix string
		if        csstype == "CLASS" { prefix = "." // It's like the Go designers did a survey to find
		} else if csstype == "ID"    { prefix = "#" // the ugliest way to express this
		} else if csstype == "TAG"   { prefix = ""}
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
	// TODO: Actually do the tags part
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
