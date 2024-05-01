-- ORDER MATTERS, APPEND NEW DATA TO END
-- remember to add structs to testdata.go as well

INSERT INTO author (
    name
) VALUES ('John Adams'), ('Alice Brown'), ('Billy Foo'), ('Carl Baz'), ('Daniel Bar');

INSERT INTO tag (
    name
) VALUES ('testTag'), ('Favourites'), ('Starred'), ('Super duper long tag for testing');

INSERT INTO book (
    title, isbn, numOfPages, rating
) VALUES ('Book 1', '1000000000', 250, 5);

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
    title, isbn, numOfPages, rating
) VALUES ('Book 2', '2000000000', 900, 4);

INSERT INTO book_author_link (
    book, author
    ) VALUES (
    (SELECT id FROM book WHERE title = 'Book 2'),
    (SELECT id FROM author WHERE name = 'Alice Brown')
);

INSERT INTO book (
    title, isbn, description, rating
) VALUES ('Many Authors', '3000000000', 'Lorem ipsum dolor sit amet, officia excepteur ex fugiat reprehenderit enim labore culpa sint ad nisi Lorem pariatur mollit ex esse exercitation amet. Nisi anim cupidatat excepteur officia. Reprehenderit nostrud nostrud ipsum Lorem est aliquip amet voluptate voluptate dolor minim nulla est proident. Nostrud officia pariatur ut officia. Sit irure elit esse ea nulla sunt ex occaecat reprehenderit commodo officia dolor Lorem duis laboris cupidatat officia voluptate. Culpa proident adipisicing id nulla nisi laboris ex in Lorem sunt duis officia eiusmod. Aliqua reprehenderit commodo ex non excepteur duis sunt velit enim. Voluptate laboris sint cupidatat ullamco ut ea consectetur et est culpa et culpa duis.', 7);

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
    ((SELECT id FROM book WHERE title = 'Many Authors'), (SELECT id FROM tag WHERE name = 'Starred')),
    ((SELECT id FROM book WHERE title = 'Many Authors'), (SELECT id FROM tag WHERE name = 'Super duper long tag for testing'));

INSERT INTO book (
    title, isbn
) VALUES ('Book 4', '4000000000');

INSERT INTO book_author_link (
    book, author
    ) VALUES
    ((SELECT id FROM book WHERE title = 'Book 4'), (SELECT id FROM author WHERE name = 'Daniel Bar'));
