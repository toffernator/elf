-- +goose Up
-- +goose StatementBegin
CREATE TABLE user (
  id INTEGER PRIMARY KEY,
  sub TEXT,
  name TEXT,
  email TEXT
);

CREATE TABLE wishlist (
  id INTEGER PRIMARY KEY,
  owner_id INTEGER,
  FOREIGN KEY(owner_id) REFERENCES  user(id)
);

CREATE TABLE product (
  id INTEGER PRIMARY KEY,
  name TEXT,
  url TEXT,
  price INTEGER
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
    'down SQL query';
-- +goose StatementEnd
