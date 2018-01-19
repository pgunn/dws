package main

import "strings"
//import "regexp"

// #############
// This code implements the markup language I used in my classic blog. It's similar to
// Mediawiki's markup. I liked it, so we implement something similar here
//
// See for the old code:
// https://github.com/pgunn/pound/blob/master/mod_perl/MyApache/POUND/POUNDMarkup.pm

func do_markup(data string, rendermode string) string {
	// For the old wikimarkup rendermode:
	//	extract_attrs() to pull attributes out
	//	linelevel_markup() iterates over lines handling
	//		wiki-style lists and doing paragraphs
	//	elevel_markup() steps over everything and handles emphasis markers
	//		and links
	//	Implement cuts if we're in a rendermode where we're showing many entries
	//	Pack it all back up and return it
	// FIXME
	return data
}

func extract_attrs(input string) (string, map[string]string) {
	// Accepts input, returns a revised version of input as well as
	// a hashmap containing attributes.
	// Anywhere in the doc where there's a string that looks like
	// [!SOMETHING] it is an attribute. If it has a : then it is
	// a key-value pair and is packed into the hashmap as such. If
	// it does not, it's just a key and it's packed into the hashmap
	// with a value of 1. The revised string does not retain these tags.
	var collector []string
	attrs := make(map[string]string)

	// TODO
	resp := strings.Join(collector, "")
	return resp, attrs
}


