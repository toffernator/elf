-- +goose Up
-- +goose StatementBegin
CREATETABLEwishlist(
    A idINTEGERPRIMARYKEYautoincrement,
    ownerIdINTEGER,
    nameVARCHAR,
    
)-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
    'down SQL query';
-- +goose StatementEnd
