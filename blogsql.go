package main

import (
	"database/sql"
	"log"
)


// Blog stuff
func get_blogentry(dbh *sql.DB, id int) (map[string]string, map[string]string) {
	// Given an id in the blogentry table, return everything
	// about it needed to display it, including tags.
	// SELECT * FROM blogentry WHERE id=$id
	// and also
	// SELECT id as tagid, tagname FROM tag WHERE tagid IN (
	//	SELECT tagid FROM blogentry_tags WHERE beid=$id)
	// TODO: Actually do the tags part
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

func get_longname_for_safe_tag(dbh *sql.DB, safename string) string {
	// Returns id for named tag
	// SELECT id FROM tag WHERE tagname=$tag
	dbq, _ := dbh.Query("SELECT name FROM tag WHERE safename=$1", safename)
	var name = ""
	for dbq.Next() {
		dbq.Scan(&name)
	}
	return name
}

func identify_blogentries_with_tag(tagid int) {
	// Returns blogentry(id) for all blogentries that have the
	// given tag
	// SELECT id FROM blogentry WHERE id IN (
	//	SELECT beid FROM blogentry_tags WHERE tagid=$tagid)
}

