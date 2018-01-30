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

func do_markup(data string, render_target string, display_context string) (string, map[string]string) {
	// 	display_context currently includes "single", "entrylist", and "review". We handle cuts differently based on that.
	//	extract_attrs() to pull attributes out
	//	linelevel_markup() iterates over lines handling
	//		wiki-style lists and doing paragraphs
	//	elevel_markup() steps over everything and handles emphasis markers
	//		and links
	//	Implement cuts if we're in a rendermode where we're showing many entries
	//	Pack it all back up and return it
	var attrs map[string]string
	data, attrs = extract_attrs(data)
	data = linelevel_markup(data)
	data = elevel_markup(data)
	if display_context == "entrylist" {
		rex_hidecuts := regexp.MustCompile(`<cut>.*?</cut>`)
		cut_hat := func(matched string) string { return "<br /><b>(Expand post to view behind cut - " + strconv.Itoa(len(matched)) + " characters)</b><br />\n" }
		data = rex_hidecuts.ReplaceAllStringFunc(data, cut_hat)
	}

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
	// native to HTML. We return a string of marked-up text.
	//
	// We might eventually want to replace this with an AST
	// TODO: Either make this aware of other format targets,
	// or change when it is used.
	var collector []string
	var parser = make(map[string]string)

	// Poor man's state machine
	parser["lastsig"] = "" // the current set of sigils
	parser["lastblank"] = "0" // This is used to coalesce multiple blank lines into 1
	//parser["firstline"] = "1" // A bit of extra smarts around the first line
	// end of state machine

	rex_blockstart := regexp.MustCompile(`^[*#: ]+`)
	rex_whitespace := regexp.MustCompile(`^\s*$`)

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if parser["lastsig"] != "" || rex_blockstart.MatchString(line) {
			// If our last line had a sigil or this one does, call the block handler
			// to do any needed wrapping
			outparse := handle_block(line, parser)
			collector = append(collector, outparse...)
		} else if parser["lastblank"] == "1" {
			// if the last line was blank,
			// 	if this line is blank too, ignore this line; we want to collapse successive blanks.
			//	otherwise reset lastblank to zero, start a paragraph marker, and add the line
			if rex_whitespace.MatchString(line) {
				continue
			} else {
				parser["lastblank"] = "0"
				collector = append(collector, "<p>")
				collector = append(collector, line)
			}
		} else if rex_whitespace.MatchString(line) { // Blank line that is not after another
			parser["lastblank"] = "1"
			collector = append(collector, "</p>") // Use this kind of thing to delimit paragraphs
		} else { // Needs no special handling
			collector = append(collector, line)
		}
	}
	// We're done with all the lines in the document at this point. 
	if parser["lastsig"] != "" {
		// If we have sigils remaining, these are lists we want to close,
		// so we should pass a blank line to the block handler to have it close all those tags.
		fakeline := ""
		outparse := handle_block(fakeline, parser)
		collector = append(collector, outparse...)
	}
	return strings.Join(collector, "")
}


func handle_block(line string, parser map[string]string) []string {
	// Part of line-level markup. Called if we are either entering, in the middle of, or
	// exiting some kind of block-level markup. Is responsible for munging the string
	// and the state variables appropriately. This is a port of some old and confusing
	// code from my prior blog engine.
	//
	// Input:
	//	input string
	// Modifies:
	//	parser (state machine)
	// Output:
	//	output lines generated from parse
	// XXX Be careful; in this port we're chainging the data flow significantly
	//my ($lineref, $lsigref, $outlineref) = @_;
	var collector []string

	line, sig := build_sig_from_line(line)
	sigparts := strings.Split(sig, "")
	if sig == parser["lastsig"] { // No signature change, so just format this as a "list" element with no enclosure changes
		sig_final_element := sigparts[len(sigparts) - 1] // Is this the most idiomatic way to do this in Go?
		opener, closer := markup_for_lspart_item(sig_final_element)
		collector = append(collector, opener + line + closer) // e.g. <li>Line</li>
	} else if sig != "" { // there was a signature change but we're still in a block
		sig_final_element := sigparts[len(sigparts) - 1] // Is this the most idiomatic way to do this in Go?
		opener, closer := markup_for_lspart_item(sig_final_element)
		collector = append(collector, lspart_diff(sig, parser["lastsig"]) + opener + line + closer)
		parser["lastsig"] = sig
	} else {	// We're transitioning entirely out of a block, so close up shop.
			// This code is not strictly needed as a seperate block but is
			// included for clarity.
		collector = append(collector, lspart_diff("", parser["lastsig"]) + line)
		parser["lastsig"] = ""
	}
	return collector
}

func build_sig_from_line(line string) (string, string) {
	// Given an input string, figure out the markup sigil portion of it.
	// Return a version of that input with the sigil removed, and separately the sigil portion
	var sigil_list []string
	var n int

	Outer:
	for n = 0; n < len(line) ; n++ {
		switch line[n] {
			case '*', '#', ':':
				sigil_list = append(sigil_list, string(line[n]))
			case ' ': // disallow nesting other structures inside preformatted text
				sigil_list = append(sigil_list, string(line[n]))
				break Outer
			default: // we hit something that's not a sigil, stop looking
				break Outer
		}
	}
	cleaned_string := line[n:] // rest of string
	return cleaned_string, strings.Join(sigil_list, "")
}

// lspart line-level markup handling. These are used to transition between
// nested list structures in the markup, opening and closing list structures
// appropriately. As an example, imagine the transitions between
// **# (unordered list wraps unordered list wraps ordered list)
// and
// * (unordered list)
// We need to close the ordered list, then the unordered list, just leaving the
// top-level list

func lspart_diff(entering string, exiting string) string {
	// Given the previous line's markup prefices and this line's prefices,
	// generate the needed transitions.
	newparts := strings.Split(entering, "")
	oldparts := strings.Split(exiting, "")
	var collector []string

	for len(newparts) > 0 || len(oldparts) > 0 {
		if len(oldparts) > 0 {
			oldbit := oldparts[0] // equivalent of Perl's shift()
			oldparts = oldparts[1:]
			if len(newparts) > 0 { // Differential!
				newbit := newparts[0]
				newparts = newparts[1:]
				if newbit != oldbit { // Close all the rest from oldparts, then open all the rest from newparts
					collector = append(collector, lspart_closeall(oldparts))
					collector = append(collector, lspart_openall(newparts))
					break
					}
				// Otherwise, just keep going
			} else { // Just close off the oldparts
				oldparts = append([]string{oldbit}, oldparts...) // Prepend. Equivalent of perl's unshift()
				collector = append(collector, lspart_closeall(oldparts))
				break
			}
		} else { // Just open newparts
			collector = append(collector, lspart_openall(newparts))
			break
		}
	}
	return strings.Join(collector, "")
}

func markup_for_lspart_item(listsig string) (string, string) {
	// Return the immediate markup to start/end an item in a list given the
	// listitem sigil
	if listsig == "#" {
		return "\t<li>", "</li>\n";
	} else if listsig == "*" {
		return "\t<li>", "</li>\n";
	} else if listsig == ":" {
		return "\t<dd>", "</dd>\n";
	} else if listsig == " " {
		return "", "" // Predefined blocks have nothing like this
	}
	// else {
	//	errorpage(text => "Internal error in markup_for_lspart_item: [listsig] is not valid key!\n");
	// }
	return "", ""
}

func markup_for_lspart(listsig string) (string, string) {
	// Return the markup to begin or end a list/markup block area given the
	// listitem sigil
	if listsig == "#" {
		return "\n<ol>\n", "</ol>\n"
	} else if listsig == "*" {
		return "\n<ul>\n", "</ul>\n"
	} else if listsig == ":" {
		return "\n<dl>\n", "</dl>\n"
	} else if listsig == " " {
		return "\n<pre>\n", "</pre>\n"
	}
//      else {
//		errorpage(text => "Internal error in get_markup_for_last_listpart: [listsig] is not valid key!\n");
//	}
	return "", ""
}

func lspart_openall(input []string) string {
	// Generate the list of list openers for all the sigils we've gathered
	var collector []string

	for _, sp := range input {
		retpart, _ := markup_for_lspart(sp)
		collector = append(collector, retpart)
	}
	return strings.Join(collector, "")
}


func lspart_closeall(input []string) string {
	// Generate the list of list closers for all the sigils we've gathered
	var collector []string
	for n := len(input) - 1; n >= 0; n-- { // iterate backwards. Wishing Go had a more elegant phrasing for this.
		sp := input[n]
		_, retpart := markup_for_lspart(sp)
		collector = append(collector, retpart)
	}
	return strings.Join(collector, "")
}

//  Element-level markup functions
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

