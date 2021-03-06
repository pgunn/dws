package main

import (
	"database/sql"
	"sort"
	"strconv"
	"strings"
	"time"
)

// #############
// bloghtml
//
// These are functions used to do layout for the blog portion of dws.
//
// These functions should not make database calls, and should just be simple
// string manipulation (or at least restricted to looking up global settings in the config
// table and URL-path stuff)
// TODO: Consider renaming these

func draw_bnode(dbh *sql.DB, bentrydata map[string]string, content string, tags map[string]string) string {
	// We'll have to extend the signature for both topics and the footer section
	var collector []string

	collector = append(collector, "<div class=\"jentry\">\n")
	collector = append(collector, "\t<div class=\"jehead\">\n")
	collector = append(collector, "\t\t<div class=\"jetitle\">" + bentrydata["title"] + "</div><!-- jetitle -->\n")
	// Time
	var zeit_int, _ = strconv.ParseInt(bentrydata["zeit"], 10, 64) // base 10, 64-bit output
	var timestring = time.Unix(zeit_int, 0).Format("2006-Jan-02 15:04:05 EST")
	collector = append(collector, "\t\t<div class=\"jeheadtime\">")
	collector = append(collector, "Date: " + timestring)
	collector = append(collector, "\t\t</div><!-- jeheadtime -->")

	if len(tags) > 0 {
		collector = append(collector, "\t<div class=\"jetagarea\">\n")
		collector = append(collector, "<div style=\"float:left;margin-right:1em;\">Tags: </div>")
		var tagparts []string
		for safetag, tag := range tags {
			tagparts = append(tagparts, " " + get_htlink(get_dispatch_path(dbh, "blogtag") + safetag, tag, false))
		}
		collector = append(collector, strings.Join(tagparts, ", "))
		collector = append(collector, "\t</div><!-- tagarea -->\n")
	}

	if music, present := bentrydata["music"]; present {
		collector = append(collector, "<div style=\"float:left;margin-right:1em;\">Music: </div>")
		collector = append(collector, music + "\n")
	}
	// TODO Code for jemisc, the extensible area for extra tabular data like music
	collector = append(collector, "</div><!-- jehead -->\n")
	collector = append(collector, "\t<div class=\"jbody\">\n")
	collector = append(collector, content + "\n")
	collector = append(collector, "\t</div><!-- jbody -->\n")
	collector = append(collector, "\t<div class=\"jetail\">\n")
	collector = append(collector, get_htlink(path_to_blogentry(dbh, bentrydata["zeit"]), "LINK/Expand", true) ) // Let people see just this entry
	// TODO Tail code here
	collector = append(collector, "\t</div><!-- jetail -->\n")
	collector = append(collector, "</div><!-- jentry -->")
	collector = append(collector, "<br /><br />\n\n")
	var ret = strings.Join(collector, "")
	return ret
}

func display_bnode(dbh *sql.DB, bentrydata map[string]string, tags map[string]string, context string) string {
	// Return HTML for a single blog entry. Called both for single-entry view
	// as well as showing a bunch on a page. This code is a lot less abstract than
	// the equivalent POUND code, since this blog engine doesn't do nearly as much.
	// title zeit body
	//if val, present := bentrydata["tags"]; present {
		// TODO: Old code turned this into a list of links to topic/tag pages
	//}

	// TODO: Footer section (also passed to bnode) with a link to just this entry

	// Render the markup language.
	var content, _ = do_markup(bentrydata["body"], "html", context)
	var ret = draw_bnode(dbh, bentrydata, content, tags)
	// draw_bnode() consumed the output of that and actually spat out the code.
	// in draw_bnode()
	return ret
}

func display_blogmain(dbh *sql.DB, title string, caption_extra string, blogimg string, topics map[string]string, num_archives int, do_feeds bool) string {
	var collector []string

	collector = append(collector, "<div id=\"toparea\">\n")
	collector = append(collector, "<div id=\"caption\">\n")
	collector = append(collector, "\t<div id=\"picarea\">\n")
	collector = append(collector, "\t\t<img src=\"" + blogimg + "\" />\n")
	collector = append(collector, "\t</div><!-- picarea -->\n")
	collector = append(collector, "\t<div id=\"topareatext\">\n")
	collector = append(collector, "\t\t<h1>" + title + "</h1>\n")
	collector = append(collector, "\t\t" + caption_extra)
	collector = append(collector, "\t</div><!-- topareatext -->\n")
	collector = append(collector, "</div><!-- caption -->\n")
	collector = append(collector, "</div><!-- toparea -->\n")

	collector = append(collector, "<div id=\"centrearea\">\n") // We leave this open
	collector = append(collector, "\t<div id=\"menupart\">\n")
	if num_archives > 0 {
		collector = append(collector, "\t\t<div id=\"archmenu\" class=\"gmenu\">\n")
		collector = append(collector, "\t\t\tArchives\n")
		collector = append(collector, "\t\t\t<div class=\"arentry\">" + get_htlink(path_to_archive_page(dbh, 1), "First Page", true) + "</div><!-- arentry -->\n")
		collector = append(collector, "\t\t\t<div class=\"arentry\">" + get_htlink(path_to_archive_page(dbh, num_archives), "Last Page", true) + "</div><!-- arentry -->\n")
		collector = append(collector, "\t\t</div><!-- archmenu -->\n")
	}
	if topics != nil {
		collector = append(collector, "\t\t<div id=\"topicmenu\" class=\"gmenu\">\n")
		collector = append(collector, "\t\t\tTopics\n")

		var topsafelist []string
		for topsafename, _ := range topics {
			topsafelist = append(topsafelist, topsafename)
		}
		sort.Strings(topsafelist)
		for _, topsafename := range topsafelist {
			collector = append(collector, "\t\t\t\t<div class=\"tmentry\">" + get_htlink(get_dispatch_path(dbh, "blogtag") + topsafename, topics[topsafename], false) + "</div><!-- tmentry -->\n")
		}
		collector = append(collector, "\t\t</div><!-- topicmenu -->")
	}
	collector = append(collector, "\t\t<br />\n")
	collector = append(collector, "\t</div><!-- menupart -->\n")

	var ret = strings.Join(collector, "")
	return ret
}
