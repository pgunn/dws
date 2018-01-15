package main

import (
	"strings"
)

// html.go
// Generic html-related code not related to what part of dws is being run
// Try not to do heavy database lifting in these.

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

func get_htlink(target string, content string, follow_ok bool) string {
	// Generic html link making
	var collector []string

	collector = append(collector, "<a href=\"" + target + "\"")
	if ! follow_ok {
		collector = append(collector, " rel=\"nofollow\"")
	}
	collector = append(collector, ">" + content + "</a>")
	var ret = strings.Join(collector, "")
	return ret
}
