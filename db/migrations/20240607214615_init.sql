-- +goose Up
-- +goose StatementBegin
CREATE TABLE user (
  id INTEGER PRIMARY KEY,
  name TEXT
);

CREATE TABLE wishlist (
  id INTEGER PRIMARY KEY,
  owner_id INTEGER,
  name string,
  image string,

  FOREIGN KEY(owner_id) REFERENCES user(id)
);

CREATE TABLE product (
  id INTEGER PRIMARY KEY,
  name TEXT,
  url TEXT,
  price INTEGER,
  currency INTEGER,
  belongs_to_id INTEGER,
  
  FOREIGN KEY(belongs_to_id) REFERENCES wishlist(id)

);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE user;
DROP TABLE wishlist;
DROP TABLE product;
-- +goose StatementEnd
