-- This sets up tables and does necessary inserts needed for the software to function.

CREATE TABLE blogentry
	(
	id SERIAL PRIMARY KEY NOT NULL,
	zeit BIGINT UNIQUE NOT NULL, -- External epochtime timestamp
	format TEXT, -- this is a string describing a rendering policy
	title TEXT,
	body TEXT,
	music TEXT, -- I like to say what music I'm listening to
	private BOOLEAN DEFAULT FALSE
	);

CREATE TABLE tag
	(
	id SERIAL PRIMARY KEY NOT NULL,
	name TEXT UNIQUE NOT NULL,
	safename TEXT UNIQUE NOT NULL,
	descrip TEXT -- Not sure if we'll use this
	);

CREATE TABLE blogentry_tags
	(
	beid INTEGER NOT NULL REFERENCES blogentry(id),
	tagid INTEGER NOT NULL REFERENCES tag(id)
	);

-- We might eventually make this hierarchical
CREATE TABLE review_topic
	(
	id SERIAL PRIMARY KEY NOT NULL,
	name TEXT UNIQUE NOT NULL,
	safename TEXT UNIQUE NOT NULL
	);

CREATE TABLE review_target
	(
	id SERIAL PRIMARY KEY NOT NULL,
	name TEXT NOT NULL, -- Unique is tempting
	safename TEXT UNIQUE NOT NULL,
	topic INTEGER REFERENCES review_topic(id)
	);

CREATE TABLE review
	(
	id SERIAL PRIMARY KEY NOT NULL,
	zeit BIGINT NOT NULL, -- Epochtime
	title TEXT,
	body TEXT,
	rating TEXT,
	target INTEGER REFERENCES review_target(id),
	hidden BOOLEAN DEFAULT(FALSE)
	);

CREATE TABLE config
	(
	id SERIAL PRIMARY KEY NOT NULL,
	name TEXT UNIQUE NOT NULL,
	value TEXT,
	avalues VARCHAR(10), -- allowable values, a type and possibly a range.
	description TEXT
	);

-- CSS theming stuff
CREATE TABLE theme
	(
	id SERIAL PRIMARY KEY NOT NULL,
	name TEXT UNIQUE NOT NULL,
	value TEXT,
	description TEXT
	);

CREATE TABLE themedata
	(
	id SERIAL PRIMARY KEY,
	themeid INTEGER REFERENCES theme NOT NULL,
	csstype VARCHAR(10) NOT NULL, -- ID, CLASS, or TAG
	csselem VARCHAR(30) NOT NULL, -- name of what we're theming
	cssprop VARCHAR(30) NOT NULL, -- the property we're tweaking
	cssval TEXT NOT NULL -- the value we're setting
	);

-- These are copied in from POUND as sample config values. We'll replace
-- them with things appropriate for DWS as we move forward.
INSERT INTO config(name, value, avalues, description) VALUES ('xmlfeed', 10, 'i[1-100]', 'How many entries to show in the default XML feeds');
INSERT INTO config(name, value, avalues, description) VALUES ('entries_per_archpage', 10, 'i[1-30]', 'How many entries are part of a blog archive page');

INSERT INTO config(name, value, avalues, description) VALUES ('blogstatic', 'http://localhost', 't[URL]', 'Base URL (includes http part) for the server');
INSERT INTO config(name, value, avalues, description) VALUES ('main_blogname', 'dachte', 't', 'Shortname of the "main" blog (if any)');

-- These are paths internal to the app
INSERT INTO config(name, value, avalues, description) VALUES('path_blogmain',		'/blog/',		't', 'Path to the main blog page');
INSERT INTO config(name, value, avalues, description) VALUES('path_blogentry',		'/blog/entries/',	't', 'Path to blog page for a single entry');
INSERT INTO config(name, value, avalues, description) VALUES('path_blogtag',		'/blog/tags/',		't', 'Path to blog page for a particular tag');
INSERT INTO config(name, value, avalues, description) VALUES('path_reviewsmain',	'/reviews/',		't', 'Path to the main reviews page');
INSERT INTO config(name, value, avalues, description) VALUES('path_reviewstopic',	'/reviews/topic/',	't', 'Path to a review topic');
INSERT INTO config(name, value, avalues, description) VALUES('path_reviewstarget',	'/reviews/on/',		't', 'Path to a review target');
INSERT INTO config(name, value, avalues, description) VALUES('path_cssmain',		'/site.css',		't', 'Path to the CSS');

-- We should store CSS in the database. Here are two themes from my previous blogging software.
INSERT INTO theme(name, description) VALUES('BaseTheme',   'All other themes use this as a basis');
-- CSS defaults
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'TAG',   'body',          'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'TAG',   'body',          'background',            '#aaaccc');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'TAG',   'body',          'font-family',           '"Verdana", sans-serif');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'TAG',   'body',          'margin',                '0 0 0 0');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'TAG',   'body',          'font-size',             '10pt');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'entrypart',     'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'entrypart',     'width',                 '70%');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'entrypart',     'margin-right',          '12%');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'entrypart',     'float',                 'right');

-- for now, entrypart mirrors reviewpart
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'reviewpart',    'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'reviewpart',    'width',                 '70%');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'reviewpart',    'margin-right',          '12%');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'reviewpart',    'float',                 'right');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'menupart',      'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'menupart',      'float',                 'right');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'menupart',      'width',                 '15%');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'menupart',      'margin-bottom',         '1em');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'gmenu',         'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'gmenu',         'background',            'lightgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'gmenu',         'color',                 'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'gmenu',         '-moz-border-radius',    '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'gmenu',         '-webkit-border-radius', '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'gmenu',         'border-radius',         '10px 10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'gmenu',         'border',                '2px solid #000');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'gmenu',         'margin-right',          '8px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'gmenu',         'padding',               '8px');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jehead',        'border',                '3px solid green');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jehead',        '-moz-border-radius',    '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jehead',        '-webkit-border-radius', '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jehead',        'border-radius',         '10px 10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jehead',        'padding',               '2px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jehead',        'background',            'rgb(150,150,140)');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revhead',       'border',                '3px solid green');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revhead',       '-moz-border-radius',    '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revhead',       '-webkit-border-radius', '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revhead',       'border-radius',         '10px 10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revhead',       'padding',               '2px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revhead',       'background',            'rgb(150,150,140)');


INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'arentry',       'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'arentry',       'B',                     'left: 0;');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'arentry',       'color',                 'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'arentry',       'background',            'lightgrey');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'amentry',       'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'amentry',       'color',                 'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'amentry',       'background',            'lightgrey');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'tmentry',       'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'tmentry',       'B',                     'left: 0;');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'tmentry',       'color',                 'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'tmentry',       'background',            'lightgrey');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'accountmenu',   'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'archmenu',      'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'topicmenu',     'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'amentry',       'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'arentry',       'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'tmentry',       'E',                     '');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'nontop',        'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'nontop',        'Position',              'relative');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'caption',       'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'caption',       'background-color',      'lightgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'caption',       '-moz-border-radius',    '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'caption',       '-webkit-border-radius', '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'caption',       'border-radius',         '10px 10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'caption',       'border',                '2px solid #000');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'caption',       'margin',                '5px 5px 5px 5px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'caption',       'height',                '140px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'caption',       'padding',               '10px');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'picarea',       'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'picarea',       'width',                 '101px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'picarea',       'height',                '130px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'picarea',       'Position',              'absolute');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'toparea',       'height',                '200px');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'topareatext',   'left',                  '130px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'topareatext',   'Position',              'absolute');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'footer',        'margin',                '0px 0px 0px 0px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'footer',        'border-top',            '1px solid darkgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'footer',        'width',                 '100%');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'footer',        'height',                '3em');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'footer',        'background',            'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'footer',        'color',                 'lightgrey');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'centrearea',    'Position',              'relative');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'centrearea',    'overflow',              'auto');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'centrearea',    'width',                 '100%');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jentry',        'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jentry',        'background',            'grey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jentry',        'color',                 'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jentry',        '-moz-border-radius',    '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jentry',        '-webkit-border-radius', '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jentry',        'border-radius',         '10px 10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jentry',        'border',                '4px groove grey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jentry',        'padding',               '4px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jentry',        'display',               'inline-block');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jentry',        'min-width',             '60%');

-- for now, mirrors jentry
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'reventry',      'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'reventry',      'background',            'grey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'reventry',      'color',                 'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'reventry',      '-moz-border-radius',    '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'reventry',      '-webkit-border-radius', '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'reventry',      'border-radius',         '10px 10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'reventry',      'border',                '4px groove grey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'reventry',      'padding',               '4px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'reventry',      'display',               'inline-block');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'reventry',      'min-width',             '60%');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetimetext',    'margin-left',           '1em');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetimetext',    'margin-top',            '0.5em');

-- for now mirrors jetimetext
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revtimetext',   'margin-left',           '1em');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revtimetext',   'margin-top',            '0.5em');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jeheadtime',    'border',                '1px solid darkgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jeheadtime',    'height',                '38px');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revheadtime',   'border',                '1px solid darkgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revheadtime',   'height',                '38px');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jeheadtimet',   'border-left',           '1px solid darkgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jeheadtimet',   'float',                 'left');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jeheadtimet',   'height',                '100%');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revheadtimet',  'border-left',           '1px solid darkgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revheadtimet',  'float',                 'left');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revheadtimet',  'height',                '100%');


INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jbody',         'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jbody',         'background',            'lightgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jbody',         'color',                 'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jbody',         'font-family',           'Monospace');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jbody',         '-moz-border-radius',    '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jbody',         '-webkit-border-radius', '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jbody',         'border-radius',         '10px 10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jbody',         'border',                '2px solid white');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jbody',         'padding',               '2px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jbody',         'clear',                 'left');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jbody',         'margin-bottom',         '.3em');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revbody',       'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revbody',       'background',            'lightgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revbody',       'color',                 'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revbody',       'font-family',           'Monospace');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revbody',       '-moz-border-radius',    '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revbody',       '-webkit-border-radius', '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revbody',       'border-radius',         '10px 10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revbody',       'border',                '2px solid white');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revbody',       'padding',               '2px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revbody',       'clear',                 'left');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revbody',       'margin-bottom',         '.3em');


INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jemisc',        'border',                '1px solid lightgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jemisc',        'margin-top',            '0px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jemisc',        'height',                '1.4em');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetopic',       'border-left',           '1px solid lightgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetopic',       'float',                 'left');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetitle',       'border',                '0px solid grey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetitle',       'margin',                '3px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetitle',       'font-family',           'Serif');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetitle',       'text-decoration',       'underline');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revtitle',      'border',                '0px solid grey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revtitle',      'margin',                '3px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revtitle',      'font-family',           'Serif');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'revtitle',      'text-decoration',       'underline');


INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'quoted',        'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'quoted',        'color',                 'green');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'quoted',        'font-family',           'monospace');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'markup2',       'font-style',            'italic');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'markup3',       'font-weight',           'bold');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'markup4',       'font-style',            'italic');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'markup4',       'font-weight',           'bold');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'logo',          'Position',              'absolute');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'logo',          'B',                     'left: 105px;');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'logo',          'width',                 '101px');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'nonlogoheader', 'Position',              'absolute');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'nonlogoheader', 'B',                     'left: 105px;');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'nonlogoheader', 'height',                '130px');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'linkarea',      'Position',              'relative');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'headerlinks',   'color',                 'white');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'headerlinks',   'Position',              'absolute');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'headerlinks',   'B',                     'left: 0;');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'headermisc',    'color',                 'purple');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'headermisc',    'Position',              'absolute');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'privatetext',   'color',                 'red');

-- POUND had two features we might eventually add back in:
-- webpaths configured in the database, and uploading of files.
-- I left those out for now because even if we do them,
-- we might do them differently.

-- The latter should be revised if we ever do privilege separation;
-- we might reasonably make a second database user with writing privs
-- that we only use to do new posts, assuming ordinary users never
-- do things that write to the database.

GRANT ALL ON ALL TABLES IN SCHEMA public TO pound;
