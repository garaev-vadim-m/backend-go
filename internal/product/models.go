package product

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	StatusId    *int    `json:"status_id"`
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	StatusId    *int    `json:"status_id"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	StatusId    *int    `json:"status_id"`
}
