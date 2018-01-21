-- This loads sample data you might consider loading to demo the software. They are not necessary for the
-- software to function.

-- Sample blogentry content
INSERT INTO blogentry(zeit, format, title, body, music) VALUES(1414343745, 'forcedtext', 'This article is about cats',  'I really like cats. I have often had them as pets', 'TMBG - Snail Shell');
INSERT INTO blogentry(zeit, format, title, body, music) VALUES(1514343745, 'forcedtext', 'This article is about blogs', 'I sometimes write blog software', 'Death Cab for Cutie - Good Help');
INSERT INTO blogentry(zeit, format, title, body, music) VALUES(1514343770, 'forcedtext', 'This article is about cars',  'I do not have a car', 'Firewater - Psychoparmacology');
INSERT INTO blogentry(zeit, format, title, body, music) VALUES(1516515509, 'forcedtext', 'This article has markup',  '''''Thoughts'''' [http://foo.com|foo]', 'Flogging Molly - Lightning Storm');

INSERT INTO tag(name, safename, descrip) VALUES('Pets', 'pets', 'When I talk about pets');
INSERT INTO tag(name, safename, descrip) VALUES('Allergies', 'sniffle', 'I sometimes am ill');

INSERT INTO blogentry_tags(beid, tagid) VALUES( (SELECT id FROM blogentry WHERE title='This article is about cats'), (SELECT id FROM tag WHERE name='Pets'));

-- Sample review content

INSERT INTO review_topic(name, safename) VALUES ('Places to Eat', 'restaurants');
INSERT INTO review_topic(name, safename) VALUES ('Board Games',   'boardgames');

INSERT INTO review_target(name, safename, topic) VALUES('Meat Monarch',   'meatmonarch',  (SELECT id FROM review_topic WHERE safename='restaurants') );
INSERT INTO review_target(name, safename, topic) VALUES('Milk Monarch',   'milkmonarch',  (SELECT id FROM review_topic WHERE safename='restaurants') );
INSERT INTO review_target(name, safename, topic) VALUES('World War Game', 'worldwargame', (SELECT id FROM review_topic WHERE safename='boardgames' ) );

INSERT INTO review(zeit, title, body, rating, target) VALUES(1515981295, 'My burger experience', 'I do not eat meat. This was a poor restaurant choice', '1/5', (SELECT id FROM review_target WHERE safename='meatmonarch'));
INSERT INTO review(zeit, title, body, rating, target) VALUES(1515981495, 'My fries experience', 'Turns out their fries are decent. My second visit went better', '4/5', (SELECT id FROM review_target WHERE safename='meatmonarch'));
INSERT INTO review(zeit, title, body, rating, target) VALUES(1215981495, 'A fun gaming night', 'Had a very pleasant game of six. Took a few hours. Had a grand time. Highly recommended', '9/10', (SELECT id FROM review_target WHERE safename='worldwargame') );

-- We should store CSS in the database. Here are two themes from my previous blogging software.
INSERT INTO theme(name, description) VALUES('EasyReading', 'Theme designed for easier reading');
INSERT INTO theme(name, description) VALUES('LitGlass',    'Funky lit-glass theme');

-- CSS
-- Load first extra theme
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='EasyReading'), 'ID',    'caption', 'color',        'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='EasyReading'), 'CLASS', 'jbody',   'background',   'lightgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='EasyReading'), 'CLASS', 'jbody',   'color',        'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='EasyReading'), 'CLASS', 'body',    'background',   'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='EasyReading'), 'CLASS', 'caption', 'background',   'lightgrey');
-- Load second extra theme
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'),    'CLASS', 'jbody',   'background',   'url(/ripple.jpg)');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'),    'CLASS', 'jbody',   'color',        'black');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'),    'TAG',   'body',    'background',   'url(/wexner.jpg)');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'),    'ID',    'caption', 'filter',       'alpha(opacity=70)');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'),    'ID',    'caption', '-moz-opacity', '0.7');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'),    'CLASS', 'jentry',  'color',        'darkgrey');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'),    'CLASS', 'jentry',  'background',   'darkblue');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'),    'CLASS', 'jentry',  '-moz-opacity', '0.9');
INSERT INTO themedata(themeid, csstype, csselem, cssprop, cssval) VALUES((SELECT id FROM theme WHERE name='LitGlass'),    'CLASS', 'jentry',  'filter',       'alpha(opacity=90)');

