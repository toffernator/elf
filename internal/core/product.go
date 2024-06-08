package core

type Product struct {
	Name      string `json:"name" db:"name"`
	Url       string `json:"url" db:"url"`
	Price     int    `json:"price" db:"price"`
	Currency  string `json:"currency" db:"currency"`
	BelongsTo int    `json:"belongsTo" db:"belongs_to_id"`
}
