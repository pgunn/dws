--

CREATE TABLE blogentry
	(
	id SERIAL PRIMARY KEY NOT NULL,
	zeit BIGINT NOT NULL, -- External epochtime timestamp
	format TEXT, -- this is a string describing a rendering policy
	title TEXT,
	body TEXT,
	music TEXT, -- I like to say what music I'm listening to
	private BOOLEAN DEFAULT(FALSE)
	);

CREATE TABLE tag
	(
	id SERIAL PRIMARY KEY NOT NULL,
	descrip TEXT, -- Not sure if we'll use this
	tagname TEXT UNIQUE NOT NULL
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
	name TEXT UNIQUE NOT NULL
	);

CREATE TABLE review_target
	(
	id SERIAL PRIMARY KEY NOT NULL,
	name TEXT NOT NULL, -- Unique is tempting
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
	avalues VARCHAR(10),
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

-- Sample blogentry content
INSERT INTO blogentry(zeit, format, title, body, music) VALUES(1414343745, 'forcedtext', 'This article is about cats', 'I really like cats. I have often had them as pets', 'TMBG - Snail Shell');
INSERT INTO blogentry(zeit, format, title, body, music) VALUES(1514343745, 'forcedtext', 'This article is about blogs', 'I sometimes write blog software', 'Death Cab for Cutie - Good Help');

-- We should store CSS in the database. Here are two themes from my previous blogging software.
INSERT INTO theme(name, description) VALUES('BaseTheme',   'All other themes use this as a basis');
INSERT INTO theme(name, description) VALUES('EasyReading', 'Theme designed for easier reading');
INSERT INTO theme(name, description) VALUES('LitGlass',    'Funky lit-glass theme');
-- CSS defaults
-- TODO: Why did I do "E"? Is it important or legacy?
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'TAG',   'body',          'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'TAG',   'body',          'background',            '#aaaccc');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'TAG',   'body',          'font-family',           '"Verdana", sans-serif');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'TAG',   'body',          'margin',                '0 0 0 0');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'TAG',   'body',          'font-size',             '10pt');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'entrypart',     'E',                     '');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'entrypart',     'width',                 '70%');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'entrypart',     'margin-right',          '12%');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'entrypart',     'float',                 'right');

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

-- No "E" for this?
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jehead',        'border',                '3px solid green');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jehead',        '-moz-border-radius',    '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jehead',        '-webkit-border-radius', '10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jehead',        'border-radius',         '10px 10px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jehead',        'padding',               '2px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jehead',        'background',            'rgb(150,150,140)');

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

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'ID',    'centrearea',    'Position',              'absolute');
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

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetimetext',    'margin-left',           '1em');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetimetext',    'margin-top',            '0.5em');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetimedesc',    'float',                 'left');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetimedesc',    'min-height',            '2em');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetimedesc',    'height',                '100%');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetimedesc',    'min-width',             '20%');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jeheadtime',    'border',                '1px solid darkgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jeheadtime',    'height',                '38px');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jeheadtimev',   'float',                 'left');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jeheadtimev',   'min-height',            '2em');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jeheadtimev',   'height',                '100%');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jeheadtimev',   'min-width',             '20%');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jeheadtimet',   'border-left',           '1px solid darkgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jeheadtimet',   'float',                 'left');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jeheadtimet',   'height',                '100%');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetailfield',   'border-left',           '2px solid rgb(170,170,170)');

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

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jemisc',        'border',                '1px solid lightgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jemisc',        'margin-top',            '0px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jemisc',        'height',                '1.4em');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetopic',       'border-left',           '1px solid lightgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetopic',       'float',                 'left');

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetitle',       'border',                '0px solid grey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetitle',       'margin',                '3px');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetitle',       'font-family',           'Serif');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'jetitle',       'text-decoration',       'underline');

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

INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'goodlink',      'color',                 'lightgreen');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'noexist',       'color',                 'red');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='BaseTheme'), 'CLASS', 'namespace',     'color',                 'orange');

-- Load first theme
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='EasyReading'), 'ID',    'caption', 'color',        'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='EasyReading'), 'CLASS', 'jbody',   'background',   'lightgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='EasyReading'), 'CLASS', 'jbody',   'color',        'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='EasyReading'), 'CLASS', 'body',    'background',   'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='EasyReading'), 'CLASS', 'caption', 'background',   'lightgrey');
-- Load second theme
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'), 'CLASS',    'jbody',   'background',   'url(/ripple.jpg)');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'), 'CLASS',    'jbody',   'color',        'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'), 'TAG',      'body',    'background',   'url(/wexner.jpg)');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'), 'ID',       'caption', 'filter',       'alpha(opacity=70)');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'), 'ID',       'caption', '-moz-opacity', '0.7');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'), 'CLASS',    'jentry',  'color',        'darkgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'), 'CLASS',    'jentry',  'background',   'darkblue');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'), 'CLASS',    'jentry',  '-moz-opacity', '0.9');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'), 'CLASS',    'jentry',  'filter',       'alpha(opacity=90)');
-- POUND had two features we might eventually add back in:
-- webpaths configured in the database, and uploading of files.
-- I left those out for now because even if we do them,
-- we might do them differently.

-- The latter should be revised if we ever do privilege separation;
-- we might reasonably make a second database user with writing privs
-- that we only use to do new posts, assuming ordinary users never
-- do things that write to the database.

GRANT ALL ON ALL TABLES IN SCHEMA public TO pound;
