PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS book (
    id            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    title         TEXT NOT NULL,
    isbn          TEXT NOT NULL UNIQUE,
    numOfPages    INTEGER DEFAULT 0,
    rating        INTEGER DEFAULT 0,
    state         TEXT NOT NULL DEFAULT "unread",
    description   TEXT,
    notes         TEXT,
    series        TEXT,
    dateAdded     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    dateUpdated   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    dateCompleted TIMESTAMP
);

-- create state enum

CREATE TABLE IF NOT EXISTS author (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS book_author_link (
	book INTEGER NOT NULL REFERENCES book(id) ON DELETE CASCADE,
	author INTEGER NOT NULL REFERENCES author(id) ON DELETE RESTRICT,
	PRIMARY KEY(book, author)
);

CREATE TABLE IF NOT EXISTS tag (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS book_tag_link (
	book INTEGER NOT NULL REFERENCES book(id) ON DELETE CASCADE,
	tag INTEGER NOT NULL REFERENCES tag(id) ON DELETE CASCADE,
	PRIMARY KEY(book, tag)
);
