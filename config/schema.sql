CREATE TABLE resources (
  name text PRIMARY KEY,
  path text,
  format text,
  kind text
);

CREATE TABLE configs (
  uri text PRIMARY KEY,
  data text
);
