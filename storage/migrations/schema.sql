CREATE TABLE IF NOT EXISTS book (
    id            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    title         TEXT NOT NULL,
    subtitle      TEXT,

    numOfPages    INTEGER DEFAULT 0,
    progress      INTEGER DEFAULT 0,
    rating        INTEGER DEFAULT 0,

    publisher     TEXT,
    datePublished TIMESTAMP,

    description   TEXT,
    notes         TEXT,
    cover         TEXT,

    dateStarted   TIMESTAMP,
    dateCompleted TIMESTAMP,
    dateAdded     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS author (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

-- M to M
CREATE TABLE IF NOT EXISTS book_author_link (
    book INTEGER NOT NULL REFERENCES book(id) ON DELETE CASCADE,
    author INTEGER NOT NULL REFERENCES author(id) ON DELETE RESTRICT,
    PRIMARY KEY(book, author)
);

CREATE TABLE IF NOT EXISTS tag (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

-- M to M
CREATE TABLE IF NOT EXISTS book_tag_link (
    book INTEGER NOT NULL REFERENCES book(id) ON DELETE CASCADE,
    tag INTEGER NOT NULL REFERENCES tag(id) ON DELETE CASCADE,
    PRIMARY KEY(book, tag)
);

-- M to 1
CREATE TABLE IF NOT EXISTS series (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    bookId INTEGER NOT NULL REFERENCES book(id) ON DELETE CASCADE,
    name TEXT NOT NULL
);

-- 1 to M
CREATE TABLE IF NOT EXISTS isbn10 (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    bookId INTEGER NOT NULL REFERENCES book(id) ON DELETE CASCADE,
    isbn TEXT NOT NULL UNIQUE
);

-- 1 to M
CREATE TABLE IF NOT EXISTS isbn13 (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    bookId INTEGER NOT NULL REFERENCES book(id) ON DELETE CASCADE,
    isbn TEXT NOT NULL UNIQUE
);

-- 1 to M
CREATE TABLE IF NOT EXISTS format (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    bookId INTEGER NOT NULL REFERENCES book(id) ON DELETE CASCADE,
    filepath TEXT NOT NULL UNIQUE
);

-- views
CREATE VIEW IF NOT EXISTS book_view AS
    SELECT b.*,
    GROUP_CONCAT(DISTINCT a.name) AS author_string,
    GROUP_CONCAT(DISTINCT t.name) AS tag_string,
    GROUP_CONCAT(DISTINCT it.isbn) AS isbn10_string,
    GROUP_CONCAT(DISTINCT ith.isbn) AS isbn13_string,
    GROUP_CONCAT(DISTINCT f.filepath) AS format_string,
    s.Name AS series_string
    FROM book b
        INNER JOIN book_author_link ba ON ba.book=b.id
        INNER JOIN author a ON ba.author=a.id
        LEFT JOIN  book_tag_link bt ON b.id=bt.book
        LEFT JOIN  tag t ON bt.tag=t.id
        LEFT JOIN  isbn10 it ON it.bookId=b.id
        LEFT JOIN  isbn13 ith ON ith.bookId=b.id
        LEFT JOIN  format f ON f.bookId=b.id
        LEFT JOIN  series s ON s.bookId=b.id
    GROUP BY b.id
    ORDER BY b.id;

-- FTS
CREATE VIRTUAL TABLE IF NOT EXISTS book_fts
	USING fts5(title, subtitle, tokenize = trigram, content = 'book', content_rowid = 'id');

CREATE TRIGGER IF NOT EXISTS book_fts_after_insert AFTER INSERT ON book BEGIN
	INSERT INTO book_fts (rowid, title, subtitle) VALUES (new.id, new.title, new.subtitle);
END;

CREATE TRIGGER IF NOT EXISTS book_fts_after_update AFTER UPDATE ON book BEGIN
  INSERT INTO book_fts (book_fts, rowid, title, subtitle) VALUES ('delete', old.id, old.title, old.subtitle);
  INSERT INTO book_fts (rowid, title, subtitle) VALUES (new.id, new.title, new.subtitle);
END;

CREATE TRIGGER IF NOT EXISTS book_fts_after_delete AFTER DELETE ON book BEGIN
  INSERT INTO book_fts (book_fts, rowid, title, subtitle) VALUES ('delete', old.id, old.title, old.subtitle);
END;


CREATE VIRTUAL TABLE IF NOT EXISTS author_fts
	USING fts5(name, tokenize = trigram, content = 'author', content_rowid = 'id');

CREATE TRIGGER IF NOT EXISTS author_fts_after_insert AFTER INSERT ON author BEGIN
	INSERT INTO author_fts (rowid, name) VALUES (new.id, new.name);
END;

CREATE TRIGGER IF NOT EXISTS author_fts_after_update AFTER UPDATE ON author BEGIN
  INSERT INTO author_fts (author_fts, rowid, name) VALUES ('delete', old.id, old.name);
  INSERT INTO author_fts (rowid, name) VALUES (new.id, new.name);
END;

CREATE TRIGGER IF NOT EXISTS author_fts_after_delete AFTER DELETE ON author BEGIN
  INSERT INTO author_fts (author_fts, rowid, name) VALUES ('delete', old.id, old.name);
END;


CREATE VIRTUAL TABLE IF NOT EXISTS tag_fts
	USING fts5(name, tokenize = trigram, content = 'tag', content_rowid = 'id');

CREATE TRIGGER IF NOT EXISTS tag_fts_after_insert AFTER INSERT ON tag BEGIN
	INSERT INTO tag_fts (rowid, name) VALUES (new.id, new.name);
END;

CREATE TRIGGER IF NOT EXISTS tag_fts_after_update AFTER UPDATE ON tag BEGIN
  INSERT INTO tag_fts (tag_fts, rowid, name) VALUES ('delete', old.id, old.name);
  INSERT INTO tag_fts (rowid, name) VALUES (new.id, new.name);
END;

CREATE TRIGGER IF NOT EXISTS tag_fts_after_delete AFTER DELETE ON tag BEGIN
  INSERT INTO tag_fts (tag_fts, rowid, name) VALUES ('delete', old.id, old.name);
END;
