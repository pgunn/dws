package main


import (
	"database/sql"
	"log"
	"strings"
)


func get_css(dbh *sql.DB, extra string) string {
	// Returns a string composed of the CSS for
	// everything DWS serves. If extra is not empty, integrates
	// the named theme into that CSS.
	var collector []string
	var ret string
	var all_css = make(map[string]map[string]map[string]string)

	dbq, err := dbh.Query("SELECT csstype, csselem, cssprop, cssval FROM themedata WHERE themeid=(SELECT id FROM theme WHERE name='BaseTheme')")
	defer dbq.Close()
	if err != nil {
		log.Print(err)
		return ret
	}
	for dbq.Next() { // Get everything we want from the database unless there's a theme
		var csstype, csselem, cssprop, cssval string
		dbq.Scan(&csstype, &csselem, &cssprop, &cssval)
		if        csstype == "CLASS" {
		} else if csstype == "ID"    {
		} else if csstype == "TAG"   {
		} else {continue} // Should probably output a warning
		// Need to fill in bits of the chained structure that are not present.
		if _, present := all_css[csstype]; !present {
			all_css[csstype] = make(map[string]map[string]string)
		}
		if _, present := all_css[csstype][csselem]; !present {
			all_css[csstype][csselem] = make(map[string]string)
		}
		all_css[csstype][csselem][cssprop] = cssval
	}
	// TODO: Code to retrieve/overlay a theme goes here
	// Next: prepare all that for display
	for csstype, _ := range all_css {
		var prefix string
		if        csstype == "CLASS" { prefix = "." // It's like the Go designers did a survey to find
		} else if csstype == "ID"    { prefix = "#" // the ugliest way to express this
		} else if csstype == "TAG"   { prefix = ""}
		for csselem, _ := range all_css[csstype] {
			collector = append(collector, prefix + csselem + "\n{\n")
			for cssprop, _ := range all_css[csstype][csselem] {
				var content = all_css[csstype][csselem][cssprop]
				if cssprop == "E" || cssprop == "B" { // These are for css properties that need verbatim handling
					if content != "" {
						collector = append(collector, content + "\n")
					}
				} else { // Most properties are key-value and fit nicely into this form instead
						collector = append(collector, cssprop + ": " + content + ";\n")
				}
			}
			collector = append(collector, "}\n\n")
		}
	}
	ret = strings.Join(collector, "")
	return ret
}

