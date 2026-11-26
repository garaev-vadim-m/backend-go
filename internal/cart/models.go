package cart

type Item struct {
	ProductID int    `json:"product_id"`
	Name      string `json:"name"`
	Price     string `json:"price"`
	Quantity  int    `json:"quantity"`
	Total     string `json:"total"`
}

type Cart struct {
	Items []Item `json:"items"`
	Total string `json:"total"`
}
