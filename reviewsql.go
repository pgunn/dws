package main

import (
	"database/sql"
)

// reviewsql.go
//
// This holds SQL operations relating to dws review functionality

func get_all_topics(dbh *sql.DB) map[string]string {
	// Returns a hashmap of { review_topic_safename => review_topic_name }
	// including all review topics
	var ret map[string]string

	dbq, err := dbh.Query("SELECT safename, name FROM review_topic ORDER BY name")
	if err != nil {
		return ret
	}
	for dbq.Next() {
		var safename string
		var name string
		dbq.Scan(&safename, &name)
		ret[safename] = name
	}
	return ret
}


func get_all_targets_in_topic(dbh *sql.DB, topicsafename string) map [string]string {
	// Returns a hashmap of { review_target_safename => review_target_name}
	// for all review targets that are under the topicid
	var ret map[string]string

	dbq, err := dbh.Query("SELECT name, safename FROM review_target WHERE topic IN (SELECT id FROM review_topic WHERE topic.safename=$1)", topicsafename)
	if err != nil {
		return ret
	}
	for dbq.Next() {
		var name string
		var safename string
		dbq.Scan(&name, &safename)
		ret[safename] = name
	}
	return ret
}


func identify_all_reviews_for_target(dbh *sql.DB, targetsafename string) []string {
	// Return reviewids for all reviews with the given target
	// TODO: Flag to filter on hidden
	var ret []string

	dbq, err := dbh.Query("SELECT id FROM review WHERE target IN (SELECT id FROM review_target WHERE review_target.safename=$1) ORDER BY zeit", targetsafename)
	if err != nil {
		return ret
	}
	for dbq.Next() {
		var id string
		dbq.Scan(&id)
		ret = append(ret, id)
	}
	return ret
}

func get_review(dbh *sql.DB, id int) map[string]string {
	// Given an id in the review table, return everything about it
	// needed to display it
	// SELECT * FROM review WHERE id=$id
	var ret map[string]string

	dbq, err := dbh.Query("SELECT title, zeit, body, rating FROM review WHERE id=$1", id)
	if err != nil {
		return ret
	}
	for dbq.Next() { // Should only return one tuple
		var title, zeit, body, rating string
		dbq.Scan(&title, &zeit, &body, &rating)
		ret["title"] = title
		ret["zeit"] = zeit
		ret["body"] = body
		ret["rating"] = rating
	}
	return ret
}

