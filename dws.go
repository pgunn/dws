package main

import (
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

// ####################################
// http dispatch functions are the main entrypoint into the code,
// and are executed by the webserver functions

func dispatch_root(w http.ResponseWriter, r *http.Request) {
	var dbh = db_connect()
	defer dbh.Close()

	var collector []string
	collector = append(collector, sthtml("Main page", true, false))
	collector = append(collector, "<ul>\n")
	collector = append(collector, "\t<li>" +  get_htlink(get_dispatch_path(dbh, "blogmain"),    "Blog",    true) + "</li>\n")
	collector = append(collector, "\t<li>" +  get_htlink(get_dispatch_path(dbh, "reviewsmain"), "Reviews", true) + "</li>\n")
	collector = append(collector, "\t<li>" +  get_htlink(get_dispatch_path(dbh, "cssmain"),     "CSS",     true) + "</li>\n")
	collector = append(collector, "</ul>\n")
	collector = append(collector, endhtml() )
	resp := strings.Join(collector, "")
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, resp)
}

func dispatch_blog_htmlview(w http.ResponseWriter, r *http.Request) {
	var dbh = db_connect()
	defer dbh.Close()

	var collector []string
	// display_blogmain() generates the framing content for a blogview, including
	// all the sidebars and topbar. If we want it to be pure-html (or at least restricted to
	//     looking up URL patterns), we need to prep the following data for it and pass it in:
	//   * all the topics the blog knows about (maybe filtered down to those that have associated entries)
	//   * the image for the blog
	//   * the calculated number of archive pages
	//   * the name of the blog owner (put this in the config table)
	//   * RSS/Atom enabledness?
	entries_per_archpage, _ := strconv.Atoi(get_config_value(dbh, "entries_per_archpage"))
	num_archivepages := get_num_archivepages(dbh, entries_per_archpage)

	collector = append(collector, sthtml("My blog", true, false))
	collector = append(collector, display_blogmain(dbh,
							get_config_value(dbh, "blogtitle"),
							"A blog by " + get_config_value(dbh, "owner"),
							get_config_value(dbh, "blogimg"),
							get_all_tags(dbh, true),
							num_archivepages,
							false)) // Retrieve URL from database, document image size
	collector = append(collector, "<div id=\"entrypart\">\n")
	var last_ten_entries = identify_last_n_blogentries(dbh, 10, false)
	for _, entryid := range last_ten_entries {
		var blogentry, tags, ok = get_blogentry(dbh, entryid)
		if !ok { // Something went wrong somehow in retrieving one of these
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			return
		}
		collector = append(collector, display_bnode(dbh, blogentry, tags, "entrylist"))
	}
	collector = append(collector, "</div><!-- entrypart -->\n")
	collector = append(collector, "</div><!-- centrearea -->\n")
	collector = append(collector, "<div id=\"footer\">\n")
	collector = append(collector, "Site served by DWS\n")
	collector = append(collector, "</div><!-- footer -->\n")
	collector = append(collector, endhtml() )
	w.Header().Set("Content-Type", "text/html") // Send HTTP headers as late as possible, ideally after errors might happen
	resp := strings.Join(collector, "")
	io.WriteString(w, resp)
}

func dispatch_blog_entry(w http.ResponseWriter, r *http.Request) {
	// We have a path (configured as a config entry named path_blogentry) but past that I want the public name of entries
	// in that virtual path to look like entry1388938708.html
	// so we're going to need to do a bit of string parsing to get the numeric part out.
	// TODO: Need to verify safety of the manipulation we're doing
	var dbh = db_connect()
	defer dbh.Close()
	var collector []string

	entries_per_archpage, _ := strconv.Atoi(get_config_value(dbh, "entries_per_archpage"))

	num_archivepages := get_num_archivepages(dbh, entries_per_archpage)

	blog_entryzeit_with_suffix := r.URL.Path[len(get_dispatch_path(dbh, "blogentry") + "entry"):]
	blog_entryzeit := strings.TrimSuffix(blog_entryzeit_with_suffix, ".html")
	beid, ok := get_beid_by_zeit(dbh, blog_entryzeit) // TODO: Error handling if the zeit does not exist
	if !ok {	// XXX We could return an error in this case, but is that useful?
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	blogentry, tags, ok := get_blogentry(dbh, beid)
	if !ok { // Retrieval failed. Really should not happen given that we just checked above.
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	collector = append(collector, sthtml("Blog: " + blogentry["title"], true, false))
	collector = append(collector, display_blogmain(dbh,
							"Blog: " + blogentry["title"],
							blogentry["title"], // Title is just the entry title
							get_config_value(dbh, "blogimg"),
							get_all_tags(dbh, true),
							num_archivepages,
							false)) // Retrieve URL from database, document image size
	collector = append(collector, "<div id=\"entrypart\">\n")
	collector = append(collector, display_bnode(dbh, blogentry, tags, "single"))
	collector = append(collector, "</div><!-- entrypart -->\n")
	collector = append(collector, "</div><!-- centrearea -->\n")
	collector = append(collector, "<div id=\"footer\">\n")
	collector = append(collector, "Site served by DWS\n")
	collector = append(collector, "</div><!-- footer -->\n")
	collector = append(collector, endhtml())
	w.Header().Set("Content-Type", "text/html") // Send HTTP headers as late as possible, ideally after errors might happen
	resp := strings.Join(collector, "")
	io.WriteString(w, resp)
}

func dispatch_blog_archive(w http.ResponseWriter, r *http.Request) {
	// Path is something like /blog/archive/page45.html
	var dbh = db_connect()
	defer dbh.Close()
	var collector []string

	page_requested_str := r.URL.Path[len(get_dispatch_path(dbh, "blogarchive") + "page"):]
	page_requested_str = strings.TrimSuffix(page_requested_str, ".html")
	page_requested, _ := strconv.Atoi(page_requested_str)

	entries_per_archpage, _ := strconv.Atoi(get_config_value(dbh, "entries_per_archpage"))

	num_archivepages := get_num_archivepages(dbh, entries_per_archpage)
	var archive_navigator []string

	if page_requested < 1 || page_requested > num_archivepages {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	if page_requested > 1 { // Title area should have a link to the prior archive page
		archive_navigator = append(archive_navigator, get_htlink(get_dispatch_path(dbh, "blogarchive") + "page" + strconv.Itoa(page_requested - 1) + ".html", "[Past]", true))
	}

	if page_requested < num_archivepages { // Title area should have a link to the next archive page
		archive_navigator = append(archive_navigator, get_htlink(get_dispatch_path(dbh, "blogarchive") + "page" + strconv.Itoa(page_requested + 1) + ".html", "[Future]", true))

	}

	collector = append(collector, sthtml("Blog Archive page " + strconv.Itoa(page_requested), true, false))
	collector = append(collector, display_blogmain(dbh,
							"Archives, page " + strconv.Itoa(page_requested),
							strings.Join(archive_navigator, ""),
							get_config_value(dbh, "blogimg"),
							nil,
							num_archivepages,
							false)) // Retrieve URL from database, document image size
	collector = append(collector, "<div id=\"entrypart\">\n")

	bentries := identify_blogentries_for_archive_page(dbh, page_requested, entries_per_archpage)
	for _, beid := range bentries {
		var blogentry, tags, ok = get_blogentry(dbh, beid)
		if !ok { // Internal error
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			return
		}
		collector = append(collector, display_bnode(dbh, blogentry, tags, "entrylist"))
	}

	collector = append(collector, "</div><!-- entrypart -->\n")
	collector = append(collector, "</div><!-- centrearea -->\n")
	collector = append(collector, "<div id=\"footer\">\n")
	collector = append(collector, "Site served by DWS\n")



	collector = append(collector, endhtml())
	w.Header().Set("Content-Type", "text/html") // Send HTTP headers as late as possible, ideally after errors might happen
	resp := strings.Join(collector, "")
	io.WriteString(w, resp)
}

func dispatch_blog_tagpage(w http.ResponseWriter, r *http.Request) {
	var dbh = db_connect()
	defer dbh.Close()
	var collector []string

	tag_safename := r.URL.Path[len(get_dispatch_path(dbh, "blogtag")):] // chop off the leading path.
	if len(tag_safename) < 1 {
		collector = append(collector, sthtml("Blog Tags", true, false))
		tags := get_all_tags(dbh, false)
		collector = append(collector, "<h1>All Tags</h1><br />\n")
		collector = append(collector, "<ul>\n")
		// Go really needs a sortrange or something like that.
		// for tagsafename, tagname := sort range tags {
		var tagsafelist []string // more legible than performant
		for tagsafename, _ := range tags {
			tagsafelist = append(tagsafelist, tagsafename)
		}
		sort.Strings(tagsafelist)
		for _, tagsafename := range tagsafelist {
			collector = append(collector, "\t<li>" + get_htlink(get_dispatch_path(dbh, "blogtag") + tagsafename, tags[tagsafename], false) + "</li>\n")
		}
		collector = append(collector, "</ul>\n")
	} else { // TODO: This page is ugly. Fix that.
		longname, ok := get_longname_for_safe_tag(dbh, tag_safename)
		if !ok {
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			return
		}
		tag_description := get_tag_description(dbh, tag_safename)
		bentries := identify_blogentries_with_tag(dbh, tag_safename)

		collector = append(collector, sthtml("Blog Tag - " + longname, true, false))
		collector = append(collector, "<h1>" + longname + "</h1><br />\n")
		collector = append(collector, tag_description + "<br />\n")
		collector = append(collector, "<ul>Entries\n")
		for _, beid := range bentries {
			bentry, _, ok := get_blogentry(dbh, beid)
			if !ok {
				http.Redirect(w, r, "/", http.StatusPermanentRedirect)
				return
			}
			collector = append(collector, "<li>" + get_htlink(path_to_blogentry(dbh, bentry["zeit"]), bentry["zeit"], true) + ": " + bentry["title"] + "</li>\n")
		}
		collector = append(collector, "</ul>\n")
	}
	collector = append(collector, endhtml() )
	w.Header().Set("Content-Type", "text/html") // Send HTTP headers as late as possible, ideally after errors might happen
	resp := strings.Join(collector, "")
	io.WriteString(w, resp)
}

func dispatch_blog_textview(w http.ResponseWriter, r *http.Request) {
	// Saving this because it's a good template for other views of the blog data
	var dbh = db_connect()
	defer dbh.Close()
	var resp = ""
	var last_ten_entries = identify_last_n_blogentries(dbh, 10, false)
	for _, entryid := range last_ten_entries {
		var blogentry, _, ok = get_blogentry(dbh, entryid)
		if !ok { // Something went wrong somehow in retrieving one of these
			http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			return
		}
		resp += "Begin blogentry\n"
		resp += "Title: " + blogentry["title"] + "\n"
		resp += "Posted: " + blogentry["zeit"] + "\n"
		resp += "Body\n-------------------\n" + blogentry["body"] + "\n--------------------\n"
		resp += "End blogentry\n\n"
	}
	w.Header().Set("Content-Type", "text/plain") // Send HTTP headers as late as possible, ideally after errors might happen
	io.WriteString(w, resp)
}

func dispatch_blog_rss(w http.ResponseWriter, r *http.Request) {
	var dbh = db_connect()
	defer dbh.Close()

	resp := do_blog_rss(dbh)
	w.Header().Set("Content-Type", "text/xml") // Send HTTP headers as late as possible, ideally after errors might happen
	io.WriteString(w, resp)
}

func dispatch_blog_atom(w http.ResponseWriter, r *http.Request) {
	var dbh = db_connect()
	defer dbh.Close()

	resp := do_blog_atom(dbh)
	w.Header().Set("Content-Type", "application/xhtml+xml") // Send HTTP headers as late as possible, ideally after errors might happen
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
	defer dbh.Close()

	resp := get_css(dbh, "")
	w.Header().Set("Content-Type", "text/css")
	w.WriteHeader(200)
	io.WriteString(w, resp)
}

func dispatch_reviews_frontpage(w http.ResponseWriter, r *http.Request) {
	// This page displays a list of links to review topics - the top level category
	// under which particular reviews are categorised. Think stuff like
	// "restaurants". It should say how many targets there are under each topic.
	var dbh = db_connect()
	defer dbh.Close()
	var collector []string

	collector = append(collector, sthtml("Reviews - Topics", true, false))
	collector = append(collector, "<ul>Review Topics</ul>\n")
	topics := get_all_topics(dbh)
	// Go really needs a sortrange or something like that.
	// for safename, name := sort range topics {
	var topsafelist []string
	for safename, _ := range topics {
		topsafelist = append(topsafelist, safename)
	}
	sort.Strings(topsafelist)
	for _, topsafename := range topsafelist {
		collector = append(collector, "\t<li>" + get_htlink(get_dispatch_path(dbh, "reviewstopic") + topsafename, topics[topsafename], true) + "</li>\n")
	}
	collector = append(collector, "</ul>\n")
	collector = append(collector, endhtml() )
	w.Header().Set("Content-Type", "text/html") // Send HTTP headers as late as possible, ideally after errors might happen
	resp := strings.Join(collector, "")
	io.WriteString(w, resp)
}

func dispatch_reviews_topical(w http.ResponseWriter, r *http.Request) {
	// This page displays a list of links to review targets under the given topic.
	// It takes a parameter and thus needs to parse its URL further.
	// The URL-pattern for these is /reviews/topic/$safename
	// (safename is a normalised version of the name that must be composed of boring characters)
	// It should say how many "thoughts" there are for each review target.
	// links go to /reviews/on/$target
	var dbh = db_connect()
	defer dbh.Close()
	var collector []string

	topic_safename := r.URL.Path[len(get_dispatch_path(dbh, "reviewstopic")):] // chop off the leading path.
	topics := get_all_topics(dbh)
	topic := topics[topic_safename]

	collector = append(collector, sthtml("Reviews - " + topic, true, false)) // todo: extend title to include topic name
	collector = append(collector, "<ul>" + topic + " Reviews:</ul>\n")
	targets := get_all_targets_in_topic(dbh, topic_safename)
	for safename, name := range targets {
		collector = append(collector, "\t<li>" + get_htlink(get_dispatch_path(dbh, "reviewstarget") + safename, name, true) + "</li>\n")
	}
	collector = append(collector, "</li>\n")
	collector = append(collector, endhtml() )
	w.Header().Set("Content-Type", "text/html") // Send HTTP headers as late as possible, ideally after errors might happen
	resp := strings.Join(collector, "")
	io.WriteString(w, resp)
}

func dispatch_reviews_target(w http.ResponseWriter, r *http.Request) {
	// This page displays all the reviews (aka thoughts) for the named review target.
	// It takes a parameter and thus needs to parse its URL further.
	// The URL-pattern for these is /reviews/on/$safename
	// (safename is a normalised version of the name that must be composed of boring characters)
	var dbh = db_connect()
	defer dbh.Close()
	var collector []string

	target_safename := r.URL.Path[len(get_dispatch_path(dbh, "reviewstarget")):] // chop off the leading path
	target, ok := get_longname_for_target(dbh, target_safename)

	if !ok {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	collector = append(collector, sthtml("Review: " + target, true, false)) // todo: extend title to include target name
	collector = append(collector, display_reviewmain())
	collector = append(collector, "<div id=\"reviewpart\">\n")

	reviewids := identify_all_reviews_for_target(dbh, target_safename)
	for _, reviewid := range reviewids {
		// XXX Right now we have title, zeit, body, and rating hooked up
		review := get_review(dbh, reviewid)
		collector = append(collector, display_rnode(review))
	}
	collector = append(collector, "</div><!-- reviewpart -->\n")
	collector = append(collector, "</div><!-- centrearea -->\n")
	collector = append(collector, "<div id=\"footer\">\n")
	collector = append(collector, "Site served by DWS\n")
	collector = append(collector, "</div><!-- footer -->\n")
	collector = append(collector, endhtml() )
	w.Header().Set("Content-Type", "text/html") // Send HTTP headers as late as possible, ideally after errors might happen
	resp := strings.Join(collector, "")
	io.WriteString(w, resp)
}

func getenv_with_default(key, fallback string) string {
	// Try to read something to the environment, with a fallback value
	env_val := os.Getenv(key)
	if len(env_val) == 0 {
		return fallback
	}
	return env_val
}

// Finally our main function
func main() {
	var dbh = db_connect()
	defer dbh.Close()

	print("Starting DWS\n")
	port := getenv_with_default("DWS_PORT", "8000")
	http.HandleFunc("/",			dispatch_root)
	http.HandleFunc(get_dispatch_path(dbh, "blogmain"),	dispatch_blog_htmlview)
	http.HandleFunc(get_dispatch_path(dbh, "blogentry"),	dispatch_blog_entry)
	http.HandleFunc(get_dispatch_path(dbh, "blogarchive"),	dispatch_blog_archive)
	http.HandleFunc(get_dispatch_path(dbh, "blogfeedrss"),	dispatch_blog_rss)
	http.HandleFunc(get_dispatch_path(dbh, "blogfeedatom"),	dispatch_blog_atom)
	http.HandleFunc(get_dispatch_path(dbh, "blogtag"),	dispatch_blog_tagpage)

	http.HandleFunc(get_dispatch_path(dbh, "reviewsmain"),	dispatch_reviews_frontpage)
	http.HandleFunc(get_dispatch_path(dbh, "reviewstopic"),	dispatch_reviews_topical)
	http.HandleFunc(get_dispatch_path(dbh, "reviewstarget"),dispatch_reviews_target)
	http.HandleFunc(get_dispatch_path(dbh, "cssmain"),	dispatch_css)
	http.ListenAndServe("localhost:" + port, nil)
}
