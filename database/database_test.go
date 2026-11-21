package database

import (
	"errors"
	"testing"

	"go-first-project/models"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	Reset()
}

func TestGetAllProducts(t *testing.T) {
	setupTestDB(t)

	products := GetAllProducts()
	if len(products) != 3 {
		t.Fatalf("expected 3 products initially, got %d", len(products))
	}
}

func TestGetProductByID(t *testing.T) {
	setupTestDB(t)

	// Valid ID
	prod, err := GetProductByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prod.ID != 1 || prod.Name != "Mechanical Keyboard" {
		t.Fatalf("unexpected product retrieved: %+v", prod)
	}

	// Non-existent ID
	_, err = GetProductByID(999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCreateProduct(t *testing.T) {
	setupTestDB(t)

	newProd := models.Product{
		Name:        "Webcam 1080p",
		Description: "HD streaming camera",
		Price:       49.99,
		Stock:       10,
	}

	created := CreateProduct(newProd)
	if created.ID != 4 {
		t.Fatalf("expected created product ID to be 4, got %d", created.ID)
	}

	// Verify it can be retrieved
	fetched, err := GetProductByID(4)
	if err != nil {
		t.Fatalf("failed to retrieve newly created product: %v", err)
	}
	if fetched.Name != "Webcam 1080p" {
		t.Fatalf("expected name Webcam 1080p, got %s", fetched.Name)
	}
}

func TestUpdateProduct(t *testing.T) {
	setupTestDB(t)

	updateData := models.Product{
		Name:        "Mechanical Keyboard Pro",
		Description: "Updated description",
		Price:       89.99,
		Stock:       30,
	}

	updated, err := UpdateProduct(1, updateData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.ID != 1 || updated.Name != "Mechanical Keyboard Pro" || updated.Price != 89.99 {
		t.Fatalf("unexpected updated product: %+v", updated)
	}

	// Non-existent ID
	_, err = UpdateProduct(999, updateData)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for non-existent product update, got %v", err)
	}
}

func TestPatchProduct(t *testing.T) {
	setupTestDB(t)

	newName := "Ultra Wireless Mouse"
	newPrice := 34.99
	patchInput := models.UpdateProductInput{
		Name:  &newName,
		Price: &newPrice,
	}

	patched, err := PatchProduct(2, patchInput)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if patched.Name != newName || patched.Price != newPrice {
		t.Fatalf("unexpected patched product: %+v", patched)
	}
	// Check untouched fields remain intact
	if patched.Stock != 50 {
		t.Fatalf("expected stock to remain 50, got %d", patched.Stock)
	}

	// Non-existent ID
	_, err = PatchProduct(999, patchInput)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for non-existent product patch, got %v", err)
	}
}

func TestDeleteProduct(t *testing.T) {
	setupTestDB(t)

	err := DeleteProduct(3)
	if err != nil {
		t.Fatalf("unexpected error deleting product 3: %v", err)
	}

	// Verify it cannot be retrieved
	_, err = GetProductByID(3)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound after deletion, got %v", err)
	}

	// Deleting again should return ErrNotFound
	err = DeleteProduct(3)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound when deleting already deleted product, got %v", err)
	}
}
