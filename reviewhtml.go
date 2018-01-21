package main

import (
	"strconv"
	"strings"
	"time"
)

// #############
// reviewhtml
//
// These are functions used to do layout for the review portion of dws.
// initially this is just a slightly-tweaked version of bloghtml, but I expect
// later to differentiate them more. If I find that they're really not different
// enough, I'll revive the "node" concept and pull that out to a separate file and make
// it responsible for the main content portion.
//
// These functions should not make database calls, and should just be simple
// string manipulation (or at least restricted to looking up global settings in the config
// table and URL-path stuff)
// TODO: Consider renaming these

func draw_rnode(revdata map[string]string, content string) string {
	var collector []string

	collector = append(collector, "<div class=\"reventry\">\n")
	collector = append(collector, "\t<div class=\"revhead\">\n")
	collector = append(collector, "\t<div class=\"revtitle\">" + revdata["title"] + "</div>\n")
	// Time
	var zeit_int, _ = strconv.ParseInt(revdata["zeit"], 10, 64) // base 10, 64-bit output
	var timestring = time.Unix(zeit_int, 0).Format(time.RFC3339)
	collector = append(collector, "\t<div class=\"revheadtime\">")
	collector = append(collector, "<div class=\"revheadtimet\">") // I don't remember why we had two divs here
	collector = append(collector, "<div class=\"revheadtimetext\">" + timestring + "</div>")
	collector = append(collector, "</div><!-- revheadtimet -->\n")
	collector = append(collector, "</div><!-- revheadtime -->\n")
	if _, present := revdata["rating"]; present { // Rating is not mandatory, omit if not present.
		collector = append(collector, "\t<div class=\"revheadrating\">\n")
		collector = append(collector, "\tRating: " + revdata["rating"] + "\n")
		collector = append(collector, "\t</div><!-- revheadrating -->\n")
	}
	// TODO Code for jemisc, the extensible area for extra tabular data like music
	collector = append(collector, "</div><!-- revhead -->\n")
	collector = append(collector, "\t<div class=\"revbody\">\n")
	collector = append(collector, "<p>" + content + "</p>\n")
	collector = append(collector, "\t</div><!-- revbody -->\n")
	collector = append(collector, "\t<div class=\"revtail\">\n")
	// TODO Tail code here
	collector = append(collector, "\t</div><!-- revtail -->\n")
	collector = append(collector, "</div><!-- reventry -->")
	collector = append(collector, "<br /><br />\n\n")
	var ret = strings.Join(collector, "")
	return ret
}

func display_rnode(revdata map[string]string) string {
	// Return HTML for a single review. 
	var content, _ = do_markup(revdata["body"], "reventryv1")
	var ret = draw_rnode(revdata, content)
	// draw_rnode() consumed the output of that and actually spat out the code.
	// in draw_rnode()
	return ret
}

func display_reviewmain() string {
	var collector []string

	collector = append(collector, "<div id=\"centrearea\">\n") // We leave this open
	var ret = strings.Join(collector, "")
	return ret
}
