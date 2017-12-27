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

-- These are copied in from POUND as sample config values. We'll replace
-- them with things appropriate for DWS as we move forward.
INSERT INTO config(name, value, avalues, description) VALUES ('xmlfeed', 10, 'i[1-100]', 'How many entries to show in the default XML feeds');
INSERT INTO config(name, value, avalues, description) VALUES ('entries_per_archpage', 10, 'i[1-30]', 'How many entries are part of a blog archive page');
INSERT INTO config(name, value, avalues, description) VALUES ('wiki_public', 0, 'b', 'Is the Wiki editable by the public?');

INSERT INTO config(name, value, avalues, description) VALUES ('blogstatic', 'http://localhost', 't[URL]', 'Base URL (includes http part) for the server');
INSERT INTO config(name, value, avalues, description) VALUES ('main_blogname', 'dachte', 't', 'Shortname of the "main" blog (if any)');
INSERT INTO config(name, value, avalues, description) VALUES ('doing_frontpage', 0, 'b', 'Are we doing a frontpage pointing at all hosted blogs?');

INSERT INTO config(name, value, avalues, description) VALUES ('postguard', 1, 'b', 'Enable postguard? This blocks some spam but will block AOL and Tor users from posting');

INSERT INTO blogentry(zeit, format, title, body, music) VALUES(1414343745, 'forcedtext', 'This article is about cats', 'I really like cats. I have often had them as pets', 'TMBG - Snail Shell');
INSERT INTO blogentry(zeit, format, title, body, music) VALUES(1514343745, 'forcedtext', 'This article is about blogs', 'I sometimes write blog software', 'Death Cab for Cutie - Good Help');

-- POUND had two features we might eventually add back in:
-- webpaths configured in the database, and uploading of files.
-- I left those out for now because even if we do them,
-- we might do them differently.

-- The latter should be revised if we ever do privilege separation;
-- we might reasonably make a second database user with writing privs
-- that we only use to do new posts, assuming ordinary users never
-- do things that write to the database.

GRANT ALL ON ALL TABLES IN SCHEMA public TO pound;
