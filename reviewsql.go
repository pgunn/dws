package main

import (
	"database/sql"
)

// reviewsql.go
//
// This holds SQL operations relating to dws review functionality

func get_all_topics(dbh *sql.DB) map[string]string {
	// Returns a hashmap of { review_topic => review_topic_id }
	// including all review topics
	var ret map[string]string

	dbq, err := dbh.Query("SELECT id, name FROM review_topic")
	if err != nil {
		return ret
	}
	for dbq.Next() {
		var topic string
		var id string
		dbq.Scan(&id, &topic)
		ret[topic] = id
	}
	return ret
}


func get_all_targets_in_topicid(dbh *sql.DB, topicid int) map [string]string {
	// Returns a hashmap of { review_target_name => review_target_ids}
	// for all review targets that are under the topicid
	var ret map[string]string

	dbq, err := dbh.Query("SELECT id, name FROM review_target WHERE topic=$1", topicid)
	if err != nil {
		return ret
	}
	for dbq.Next() {
		var name string
		var id string
		dbq.Scan(&id, &name)
		ret[name] = id
	}
	return ret
}


func identify_all_reviews_for_targetid(dbh *sql.DB, targetid int) []string {
	// Return reviewids for all reviews with the given target
	// TODO: Flag to filter on hidden
	var ret []string

	dbq, err := dbh.Query("SELECT id FROM review WHERE target=$1 ORDER BY zeit", targetid)
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

