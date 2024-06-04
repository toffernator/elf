package core

type Product struct {
	Name     string  `json:"name"`
	Url      string  `json:"url"`
	Price    float32 `json:"price"`
	Currency string  `json:"currency"`
}
