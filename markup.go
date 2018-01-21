package main

import (
	"strings"
	"strconv"
	"regexp"
)

// #############
// This code implements the markup language I used in my classic blog. It's similar to
// Mediawiki's markup. I liked it, so we implement something similar here
//
// See for the old code:
// https://github.com/pgunn/pound/blob/master/mod_perl/MyApache/POUND/POUNDMarkup.pm
//
// It'd be amazing if Go's regex library had a ReplaceAllStringFuncCaptured, which
// might pass to the func a []string of the capture group matches rather than
// the entire match. It does not, so the functions we use need to clean some things out.

func do_markup(data string, rendermode string) (string, map[string]string) {
	// For the old wikimarkup rendermode:
	//	extract_attrs() to pull attributes out
	//	linelevel_markup() iterates over lines handling
	//		wiki-style lists and doing paragraphs
	//	elevel_markup() steps over everything and handles emphasis markers
	//		and links
	//	Implement cuts if we're in a rendermode where we're showing many entries
	//	Pack it all back up and return it
	// FIXME
	var attrs map[string]string
	data, attrs = extract_attrs(data)
	data = linelevel_markup(data)
	data = elevel_markup(data)
	// TODO cuts work here

	return data, attrs
}

// extract_attrs is self-contained
func extract_attrs(input string) (string, map[string]string) {
	// Accepts input, returns a revised version of input as well as
	// a hashmap containing attributes.
	// Anywhere in the doc where there's a string that looks like
	// [!SOMETHING] it is an attribute. If it has a : then it is
	// a key-value pair and is packed into the hashmap as such. If
	// it does not, it's just a key and it's packed into the hashmap
	// with a value of 1. The revised string does not retain these tags.
	//
	// In Perl, the main loop for this was the much prettier
	// while(s/\[!(.*?)\]//)
	//	{
	//	# Stash in attrs
	//	}
	var taglist []string
	attrs := make(map[string]string)

	extractpr := regexp.MustCompile(`\[\!(.*?)\]`)
        saver := func(matched string) string { taglist = append(taglist, matched) ; return "" ; }

	output := extractpr.ReplaceAllStringFunc(input, saver)

	for _, unparsed := range taglist {
		parts := strings.SplitN(unparsed, ":", 2)
		if len(parts) == 1 { // non-k-v tag
			attrs[parts[0]] = "1"
		} else { // k-v tag, so remember its value
			attrs[parts[0]] = parts[1]
		}
	}
	return output, attrs
}

// Line-level markup functions
func linelevel_markup(input string) string {
	// This turns wikimarkup lists and paragraph formation into something
	// native to HTML.
	//
	// We might eventually want to replace this with an AST
	// TODO: Either make this aware of other format targets,
	// or change when it is used.
	var ret []string
	var parser = make(map[string]string)

	// Poor man's state machine
	parser["listsig"] = "" // Block-level information
	parser["lastblank"] = "1" // This is used to coalesce multiple blank lines into 1
	// end of state machine

	// regex_blockstart := regexp.MustCompile("^[*#: ]+")

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		//if parser['listsig'] != '' || regex_blockstart.MatchString(line) {
		//	handle_block()
		//}
		ret = append(ret, line + "\n")
	}
	return strings.Join(ret, "")
}


func handle_block() {

}

func build_sig_from_line() {

}

func markup_for_lspart_item() {

}

func markup_for_lspart() {

}

func lspart_closeall() {

}

func lspart_openall() {

}

func lspart_diff() {

}

// Element-level markup functions
func elevel_markup(input string) string {
	// This is the main entry point for element-level markup.
	// We presently handle exactly two types of markup here,
	//   1) Doing 2-4 single quote markers in a row starts or ends an emphasis markup
	//   2) Opening and closing brackets mark the start and the end of an external link, respectively.
	//      External links are either just URL, or pipe-delimited URL|Text
	//
	// We may eventually support inner links of various sorts again; the (presently reserved)
	// syntax for that would be using double-brackets and then an inner-link type, a colon, and
	// the payload.
	//
	// We must run these regexes in order.
	// Someday consider doing an AST instead so we might validate and reject invalid code more;
	// badly-crafted markup here can generate illegal html.
	rex_emph_four  := regexp.MustCompile(`''''(.*?)''''`)
	rex_emph_three := regexp.MustCompile(`'''(.*?)'''`)
	rex_emph_two   := regexp.MustCompile(`''(.*?)''`)
	rex_extlink    := regexp.MustCompile(`\[(.*?)\]`)

	four_parser  := func(matched string) string { return markup_emphasis(matched, 4) }
	three_parser := func(matched string) string { return markup_emphasis(matched, 3) }
	two_parser   := func(matched string) string { return markup_emphasis(matched, 2) }
	ext_parser   := func(matched string) string { return markup_ext_link(matched) }

	input =  rex_emph_four.ReplaceAllStringFunc(input, four_parser)
	input = rex_emph_three.ReplaceAllStringFunc(input, three_parser)
	input =   rex_emph_two.ReplaceAllStringFunc(input, two_parser)
	input =    rex_extlink.ReplaceAllStringFunc(input, ext_parser)

	return input
}

func markup_emphasis(input string, level int) string {
	// Wikimarkup emphasis marks are done with a set of single quotes, the number of which
	// indicate the style of emphasis. We use CSS to actually paint the text accordingly.
	// This function is a little more flexible than its caller; to match MediaWiki's markup
	// we only support 2-4 single quotes in the markup, but being more flexible here is
	// easier to code and harmless.
	emphasis_extractor := strings.NewReplacer(strings.Repeat("'", level), "")
	input = emphasis_extractor.Replace(input)
	return "<span class=\"markup" + strconv.Itoa(level) + "\">" + input + "</span>"
}

func markup_ext_link(input string) string {
	// Wikimarkup external links look like [http://www.foo.com|The foo website]
	// This function gets one of those strings with the enclosing [] already removed,
	// and returns a link.
	// It is legal not to have the pipe, in which case the link text is just "LINK"
	enclosure_extractor := strings.NewReplacer("[", "", "]", "")
	input = enclosure_extractor.Replace(input)
	parts := strings.Split(input, "|")
	var linktext string
	if len(parts) == 1 {
		linktext = "LINK"
	} else {
		linktext = parts[1]
	}
	return get_htlink(parts[0], linktext, true)
}

