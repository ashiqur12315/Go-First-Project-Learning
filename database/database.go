package database

import (
	"errors"
	"sync"

	"go-first-project/models"
)

// ErrNotFound is returned when a requested product does not exist in the database.
var ErrNotFound = errors.New("product not found")

// In-memory data store protected by a read-write mutex.
var (
	mu       sync.RWMutex
	products = []models.Product{
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
	nextID = 4
)

// Reset restores the in-memory product database to its default seed data.
// Useful for automated tests to ensure test isolation.
func Reset() {
	mu.Lock()
	defer mu.Unlock()

	products = []models.Product{
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
	nextID = 4
}

// GetAllProducts returns a copy of all products currently stored.
func GetAllProducts() []models.Product {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]models.Product, len(products))
	copy(result, products)
	return result
}

// GetProductByID retrieves a single product by its ID.
func GetProductByID(id int) (*models.Product, error) {
	mu.RLock()
	defer mu.RUnlock()

	for _, p := range products {
		if p.ID == id {
			found := p
			return &found, nil
		}
	}
	return nil, ErrNotFound
}

// CreateProduct adds a new product to the database and assigns it a unique ID.
func CreateProduct(p models.Product) models.Product {
	mu.Lock()
	defer mu.Unlock()

	p.ID = nextID
	nextID++
	products = append(products, p)
	return p
}

// UpdateProduct completely replaces an existing product by its ID.
func UpdateProduct(id int, p models.Product) (*models.Product, error) {
	mu.Lock()
	defer mu.Unlock()

	for i, existing := range products {
		if existing.ID == id {
			p.ID = id
			products[i] = p
			updated := products[i]
			return &updated, nil
		}
	}
	return nil, ErrNotFound
}

// PatchProduct applies partial updates to an existing product by its ID.
func PatchProduct(id int, input models.UpdateProductInput) (*models.Product, error) {
	mu.Lock()
	defer mu.Unlock()

	for i, p := range products {
		if p.ID == id {
			if input.Name != nil {
				products[i].Name = *input.Name
			}
			if input.Description != nil {
				products[i].Description = *input.Description
			}
			if input.Price != nil {
				products[i].Price = *input.Price
			}
			if input.Stock != nil {
				products[i].Stock = *input.Stock
			}
			patched := products[i]
			return &patched, nil
		}
	}
	return nil, ErrNotFound
}

// DeleteProduct removes a product by its ID from the database.
func DeleteProduct(id int) error {
	mu.Lock()
	defer mu.Unlock()

	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
