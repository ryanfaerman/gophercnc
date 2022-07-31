CREATE TABLE authors (
  id   integer PRIMARY KEY AUTOINCREMENT,
  name text      NOT NULL,
  bio  text
);

CREATE TABLE libraries (
  name text  PRIMARY KEY,
  path text not null,
  active bool DEFAULT false
)
