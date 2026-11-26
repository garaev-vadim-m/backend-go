package product

import "hello-go-backend/internal/status"

type Product struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description *string        `json:"description"`
	Status      *status.Status `json:"status"`
	Price       string         `json:"price"`
	Stock       int            `json:"stock"`
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	StatusId    *int    `json:"status_id"`
	Price       string  `json:"price"`
	Stock       int     `json:"stock"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	StatusId    *int    `json:"status_id"`
	Price       string  `json:"price"`
	Stock       int     `json:"stock"`
}
