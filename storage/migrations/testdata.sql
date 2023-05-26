-- ORDER MATTERS, APPEND NEW DATA TO END
-- remember to add structs to testdata.go as well

INSERT INTO author (
	name
) VALUES ('John Adams'), ('Alice Brown'), ('Billy Foo'), ('Carl Baz'), ('Daniel Bar');

INSERT INTO tag (
	name
) VALUES ('testTag'), ('Favourites'), ('Starred');

INSERT INTO book (
	title, isbn, numOfPages, rating, state
) VALUES ('Book 1', '1', 250, 5, 'read');

INSERT INTO book_author_link (
	book, author
	) VALUES (
	(SELECT id FROM book WHERE title = 'Book 1'),
	(SELECT id FROM author WHERE name = 'John Adams')
);

INSERT INTO book_tag_link (
	book, tag
	) VALUES (
	(SELECT id FROM book WHERE title = 'Book 1'),
	(SELECT id FROM tag WHERE name = 'testTag')
);

INSERT INTO book (
	title, isbn, numOfPages, rating, state
) VALUES ('Book 2', '2', 900, 4, 'unread');

INSERT INTO book_author_link (
	book, author
	) VALUES (
	(SELECT id FROM book WHERE title = 'Book 2'),
	(SELECT id FROM author WHERE name = 'Alice Brown')
);

INSERT INTO book (
	title, isbn
) VALUES ('Many Authors', '3');

INSERT INTO book_author_link (
	book, author
	) VALUES
	((SELECT id FROM book WHERE title = 'Many Authors'), (SELECT id FROM author WHERE name = 'Billy Foo')),
	((SELECT id FROM book WHERE title = 'Many Authors'), (SELECT id FROM author WHERE name = 'Carl Baz')),
	((SELECT id FROM book WHERE title = 'Many Authors'), (SELECT id FROM author WHERE name = 'Daniel Bar'));

INSERT INTO book_tag_link (
	book, tag
	) VALUES
	((SELECT id FROM book WHERE title = 'Many Authors'), (SELECT id FROM tag WHERE name = 'Favourites')),
	((SELECT id FROM book WHERE title = 'Many Authors'), (SELECT id FROM tag WHERE name = 'Starred'));

INSERT INTO book (
	title, isbn
) VALUES ('Book 4', '4');

INSERT INTO book_author_link (
	book, author
	) VALUES
	((SELECT id FROM book WHERE title = 'Book 4'), (SELECT id FROM author WHERE name = 'Daniel Bar'));
