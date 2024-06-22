-- +goose Up
-- +goose StatementBegin
INSERT INTO user (name) 
VALUES
    ("test user 1"),
    ("test user 2"),
    ("test user 3");
INSERT INTO wishlist (owner_id, name)
VALUES
    (1, "test wishlist 1 belonging to user with id 1"),
    (1, "test wishlist 2 belonging to user with id 1"),
    (2, "test wishlist 3 belonging to user with id 2");
INSERT INTO product (name, url, price, currency, belongs_to_id)
VALUES
    ("iPad Pro", "https://www.apple.com/ipad-pro/", 89900, 0, 1),
    ("iPhone 15 Pro", "https://www.apple.com/iphone-15-pro/", 99900, 0, 2),
    ("Macbook Pro", "https://www.apple.com/macbook-pro/", 149900, 0, 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM user;
DELETE FROM wishlist;
DELETE FROM product;
-- +goose StatementEnd
