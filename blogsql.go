package main

import (
	"database/sql"
)


// Blog stuff
func get_blogentry(dbh *sql.DB, id string) (map[string]string, map[string]string) {
	// Given an id in the blogentry table, return everything
	// about it needed to display it, including tags.
	var mymap = make(map[string]string)
	var tags = make(map[string]string)

	dbq, _ := dbh.Query("SELECT title, zeit, body FROM blogentry WHERE id=$1", id)
	for dbq.Next() {
		var title, zeit, body string
		dbq.Scan(&title, &zeit, &body)
		mymap["title"] = title
		mymap["zeit"] = zeit
		mymap["body"] = body
	}
	dbq, _ = dbh.Query("SELECT name, safename FROM tag WHERE id IN (SELECT tagid FROM blogentry_tags WHERE beid=$1)", id)
	for dbq.Next() {
		var tagname, safetagname string
		dbq.Scan(&tagname, &safetagname)
		tags[safetagname] = tagname
	}
	return mymap, tags
}

func get_beid_by_zeit(dbh *sql.DB, zeit string) string {
	dbq, _ := dbh.Query("SELECT id FROM blogentry where zeit=$1", zeit)
	var beid = ""
	for dbq.Next() {
		dbq.Scan(&beid)
	}
	return beid
}

func identify_last_n_blogentries(dbh *sql.DB, count int, include_private bool) []string {
	// Returns blogentry(id) for the last up-to-$count blogentries
	// NOTE: It can return fewer than requested if there are not that many.
	// It will return the entries ordered so newer are earlier.
	var ret []string
	var dbq *sql.Rows

	if include_private {
		dbq, _ = dbh.Query("SELECT id FROM blogentry ORDER BY zeit LIMIT $2", include_private, count)
	} else {
		dbq, _ = dbh.Query("SELECT id FROM blogentry WHERE private=false ORDER BY zeit LIMIT $1", count)
	}
	for dbq.Next() {
		var retval string
		dbq.Scan(&retval)
		ret = append(ret, retval)
	}
	return ret
}

func get_all_tags(dbh *sql.DB, include_empty bool) map[string]string {
	var ret = make(map[string]string)
	var dbq *sql.Rows

	if include_empty {
		dbq, _ = dbh.Query("SELECT name, safename FROM tag")
	} else {
		dbq, _ = dbh.Query("SELECT name, safename FROM tag WHERE id IN (SELECT tagid FROM blogentry_tags)")
	}
	for dbq.Next() {
		var name, safename string
		dbq.Scan(&name, &safename)
		ret[safename] = name
	}
	return ret
}

func get_longname_for_safe_tag(dbh *sql.DB, safename string) string {
	// Returns id for named tag
	// SELECT id FROM tag WHERE tagname=$tag
	var name = ""

	dbq, _ := dbh.Query("SELECT name FROM tag WHERE safename=$1", safename)
	for dbq.Next() {
		dbq.Scan(&name)
	}
	return name
}

func get_tag_description(dbh *sql.DB, safename string) string {
	var descrip = "No description yet"

	dbq, _ := dbh.Query("SELECT descrip FROM tag WHERE safename=$1", safename)
	for dbq.Next() {
		dbq.Scan(&descrip)
	}
	return descrip

}

func identify_blogentries_with_tag(dbh *sql.DB, safename string) []string {
	// Returns blogentry(id) for all blogentries that have the
	// given tag
	var ret []string

	dbq, _ := dbh.Query("SELECT id FROM blogentry WHERE id IN (SELECT beid FROM blogentry_tags WHERE tagid=(SELECT id FROM tag WHERE safename=$1))", safename)
	for dbq.Next() {
		var beid string
		dbq.Scan(&beid)
		ret = append(ret, beid)
	}
return ret
}

