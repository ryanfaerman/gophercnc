-- +goose Up
-- +goose StatementBegin
CREATE TABLE resources (
  id integer PRIMARY KEY,
  name text NOT NULL,
  path text,
  format text,
  kind text,
  UNIQUE(name, kind) ON CONFLICT ABORT
);


CREATE TABLE configs (
  uri text PRIMARY KEY,
  data text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE resources;
DROP TABLE configs;
-- +goose StatementEnd
