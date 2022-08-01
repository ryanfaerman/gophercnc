-- +goose Up
-- +goose StatementBegin
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE resources;
DROP TABLE configs;
-- +goose StatementEnd
