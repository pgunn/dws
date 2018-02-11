package main

import (
	"database/sql"
	"strings"
	"strconv"
	"time"
	"github.com/gorilla/feeds"
)

// feeds.go
//
// This does everything relating to the RSS/Atom feeds for the blog, from database calls to formatting.

func do_blog_atom(dbh *sql.DB) string {
	last_ten_entries := identify_last_n_blogentries(dbh, 10, false) // Consider making the number configurable
	now := time.Now()

	feed := &feeds.Feed{
		Title:		get_config_value(dbh, "blogtitle"),
		Link:		&feeds.Link{Href: get_config_value(dbh, "blogstatic")},
		Description:	"Blog of " + get_config_value(dbh, "owner"),
		Author:		&feeds.Author{Name: get_config_value(dbh, "owner")},
		Created:	now,
		Items:		make([]*feeds.Item, 0),
	}

	for _, entryid := range last_ten_entries {

		var blogentry, _, _ = get_blogentry(dbh, entryid)
		blogentry_link := path_to_blogentry(dbh, blogentry["zeit"])
		zeit_int, _ := strconv.ParseInt(blogentry["zeit"], 10, 64)
		timeobj := time.Unix(zeit_int, 0)
		rendered, _ := do_markup(blogentry["body"], "rss", "entrylist")
		feedEntry := &feeds.Item{
			Title:		blogentry["title"],
			Link:		&feeds.Link{Href: blogentry_link},
			Description:	rendered,
			Author:		&feeds.Author{Name: get_config_value(dbh, "owner")},
			Created:	timeobj,
		}
		feed.Items = append(feed.Items, feedEntry)
	}
	ret, _ := feed.ToAtom()
	return ret
}

func do_blog_rss(dbh *sql.DB) string {
	var collector []string

	last_ten_entries := identify_last_n_blogentries(dbh, 10, false) // Consider making the number configurable
	collector = append(collector, do_rss_header(dbh))
	collector = append(collector, do_rss_sequence(dbh, last_ten_entries))

	return strings.Join(collector, "")

}


// -------------------------------

func do_rss_header(dbh *sql.DB) string {
	var collector []string

	collector = append(collector, "<rdf:RDF xmlns:rdf=\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\"\n")
	collector = append(collector, "xmlns:dc=\"http://purl.org/dc/elements/1.1/\"\n")
	collector = append(collector, "xmlns:sy=\"http://purl.org/rss/1.0/modules/syndication/\"\n")
	collector = append(collector, "xmlns:content=\"http://purl.org/rss/1.0/modules/content/\"\n")
	collector = append(collector, "xmlns=\"http://purl.org/rss/1.0/\">\n")
	collector = append(collector, "<channel rdf:about=\"" + get_config_value(dbh, "blogstatic") + "\">\n")
	collector = append(collector, "<title>" + get_config_value(dbh, "blogtitle") + "</title>\n")
	collector = append(collector, "<link>" + get_config_value(dbh, "blogstatic") + "</link>\n")
	collector = append(collector, "<description>Blog of " + get_config_value(dbh, "owner") + "</description>\n")
	collector = append(collector, "<language>en-us</language>\n") // TODO: make this a config value
	return strings.Join(collector, "")
}

func do_rss_sequence(dbh *sql.DB, entrylist []string) string {
	// RSS is a bit of an odd format and it wants some upfront metadata on entries first, followed by
	// payloads for those entries. So we get to iterate twice.
	var collector []string
	// First part is a resource list
	collector = append(collector, "<items>\n")
	collector = append(collector, "<rdf:Seq>\n")

	for _, entryid := range entrylist {
		var blogentry, _, _ = get_blogentry(dbh, entryid)
		blogentry_link := path_to_blogentry(dbh, blogentry["zeit"])
		collector = append(collector, "<rdf:li rdf:resource=\"" + blogentry_link + "\" />\n")
	}
	collector = append(collector, "</rdf:Seq>\n")
	collector = append(collector, "</items>\n")
	collector = append(collector, "</channel>\n")
	// Now for the second part, where we deliver payloads for preview.
	for _, entryid := range entrylist { // If we were doing a bigger list, caching would make sense. Or writing both string-sets in one iteration.
		var blogentry, _, _ = get_blogentry(dbh, entryid)
		blogentry_link := path_to_blogentry(dbh, blogentry["zeit"])

		rendered, _ := do_markup(blogentry["body"], "rss", "entrylist")
		brace_extractor := strings.NewReplacer("[", "_", "]", "_")
		rendered = brace_extractor.Replace(rendered) // brace markers break the parser I think
		collector = append(collector, "<item rdf:about=\"" + blogentry_link + "\">\n")
		collector = append(collector, "<title>" + blogentry["title"] + "</title>\n")
		collector = append(collector, "<link>" + blogentry_link + "</link>\n")
		collector = append(collector, "<description>" + blogentry["title"] + "</description>\n")
		collector = append(collector, "<dc:creator>" + get_config_value(dbh, "owner") + "</dc:creator>\n")
		// The designers of RSS had unfortunate taste in date strings
		zeit_int, _ := strconv.ParseInt(blogentry["zeit"], 10, 64)
		timestring := time.Unix(zeit_int, 0).Format("2006-01-02T15:04:05Z07:00") // RSS1 uses ISO8601 strings for time

		collector = append(collector, "<dc:date>" + timestring + "</dc:date>\n")
		collector = append(collector, "<content:encoded><![CDATA[" + rendered + "]]></content:encoded>\n")
		collector = append(collector, "</item>")
	}

	return strings.Join(collector, "")
}
