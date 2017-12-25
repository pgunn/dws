--

CREATE TABLE blogentry
	(
	id SERIAL PRIMARY KEY NOT NULL,
	zeit BIGINT NOT NULL, -- External epochtime timestamp
	format TEXT, -- this is a string describing a rendering policy
	title TEXT,
	body TEXT,
	music TEXT, -- I like to say what music I'm listening to
	hidden BOOLEAN DEFAULT(FALSE)
	);

CREATE TABLE tag
	(
	id SERIAL PRIMARY KEY NOT NULL,
	descrip TEXT, -- Not sure if we'll use this
	tagname TEXT UNIQUE NOT NULL
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
	topic INTEGER REFERENCES reivew_topic(id)
	);

CREATE TABLE review
	(
	id SERIAL PRIMARY KEY NOT NULL,
	zeit BIGINT NOT NULL, -- Epochtime
	title TEXT,
	body TEXT,
	rating TEXT,
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

INSERT INTO config(name,value, avalues, description) VALUES ('xmlfeed', 10, 'i[1-100]', 'How many entries to show in the default XML feeds');
INSERT INTO config(name,value, avalues, description) VALUES ('entries_per_archpage', 10, 'i[1-30]', 'How many entries are part of a blog archive page');
INSERT INTO config(name,value, avalues, description) VALUES ('wiki_public', 0, 'b', 'Is the Wiki editable by the public?');

INSERT INTO config(name,value, avalues, description) VALUES ('blogstatic', 'http://localhost', 't[URL]', 'Base URL (includes http part) for the server');
INSERT INTO config(name,value, avalues, description) VALUES ('main_blogname', 'dachte', 't', 'Shortname of the "main" blog (if any)');
INSERT INTO config(name,value, avalues, description) VALUES ('doing_frontpage', 0, 'b', 'Are we doing a frontpage pointing at all hosted blogs?');

INSERT INTO config(name,value, avalues, description) VALUES ('postguard', 1, 'b', 'Enable postguard? This blocks some spam but will block AOL and Tor users from posting');

-- POUND had two features we might eventually add back in:
-- webpaths configured in the database, and uploading of files.
-- I left those out for now because even if we do them,
-- we might do them differently.
