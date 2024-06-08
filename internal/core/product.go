package core

type Product struct {
	Id        string `db:"id"`
	Name      string `db:"name"`
	Url       string `db:"url"`
	Price     int    `db:"price"`
	Currency  string `db:"currency"`
	BelongsTo int    `db:"belongs_to_id"`
}
