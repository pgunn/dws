package main

import (
	"database/sql"
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
	collector = append(collector, "\t<div class=\"jetitle\">Title: " + bentrydata["title"] + "</div>\n")
	// Time
	var zeit_int, _ = strconv.ParseInt(bentrydata["zeit"], 10, 64) // base 10, 64-bit output
	var timestring = time.Unix(zeit_int, 0).Format(time.RFC3339)
	collector = append(collector, "\t<div class=\"jeheadtime\">")
	collector = append(collector, "<div class=\"jeheadtimet\">") // I don't remember why we had two divs here
	collector = append(collector, "<div class=\"jeheadtimetext\">Date: " + timestring + "</div>\n")
	collector = append(collector, "</div><!-- jeheadtimet -->")
	collector = append(collector, "</div><!-- jeheadtime -->")

	if len(tags) > 0 {
		collector = append(collector, "\t<div class=\"jetagarea\">\n")
		collector = append(collector, "Tags: ")
		for safetag, tag := range tags {
			collector = append(collector, " " + get_htlink(get_dispatch_path(dbh, "blogtag") + safetag, tag, false))
		}
		collector = append(collector, "\t</div><!-- tagarea -->\n")
	}
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

func display_bnode(dbh *sql.DB, bentrydata map[string]string, tags map[string]string) string {
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
	var ret = draw_bnode(dbh, bentrydata, content, tags)
	// draw_bnode() consumed the output of that and actually spat out the code.
	// in draw_bnode()
	return ret
}

func display_blogmain(title string, owner string, blogimg string, topics []string, num_archives int, do_feeds bool) string {
	var collector []string
	caption_extra := "A blog by " + owner

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
