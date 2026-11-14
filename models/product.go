package models

import "sync"

// Product represents our data model for a product.
// Capitalized (Exported) so it is accessible from handlers and other packages.
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// UpdateProductInput is used for PATCH requests where fields are optional.
type UpdateProductInput struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
}

// In-memory product storage with initial seed data.
// Exported variables (capitalized) allow handlers to read and write them.
var (
	Mu       sync.RWMutex
	Products = []Product{
		{
			ID:          1,
			Name:        "Mechanical Keyboard",
			Description: "RGB mechanical gaming keyboard with blue switches",
			Price:       79.99,
			Stock:       25,
		},
		{
			ID:          2,
			Name:        "Wireless Mouse",
			Description: "Ergonomic 2.4GHz wireless mouse with adjustable DPI",
			Price:       29.99,
			Stock:       50,
		},
		{
			ID:          3,
			Name:        "Gaming Monitor",
			Description: "27-inch 144Hz IPS display with 1ms response time",
			Price:       199.99,
			Stock:       15,
		},
	}
	NextID = 4
)
