package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/brenoproti/go-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateNewProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	product, err := entity.NewProduct("Product 1", 10.5)
	assert.NoError(t, err)
	productDb := NewProductDB(db)
	err = productDb.Create(product)
	assert.NoError(t, err)
	product, err = productDb.FindById(product.ID.String())
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)
}

func TestFindAllProducts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	productDb := NewProductDB(db)
	for i := 1; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100)
		assert.NoError(t, err)
		err = productDb.Create(product)
		assert.NoError(t, err)
	}
	products, err := productDb.FindAll(1, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = productDb.FindAll(2, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)

	products, err = productDb.FindAll(3, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Equal(t, "Product 21", products[0].Name)
	assert.Equal(t, "Product 23", products[2].Name)

	products, err = productDb.FindAll(1, 24, "desc")
	assert.NoError(t, err)
	assert.Len(t, products, 23)
	assert.Equal(t, "Product 23", products[0].Name)
	assert.Equal(t, "Product 1", products[22].Name)
}

func TestFindById(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	productDb := NewProductDB(db)
	product, err := entity.NewProduct("Product 1", 10.5)
	assert.NoError(t, err)
	err = productDb.Create(product)
	assert.NoError(t, err)
	product, err = productDb.FindById(product.ID.String())
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)
	product, err = productDb.FindById("123")
	assert.Error(t, err)
	assert.Empty(t, product)
}

func TestUpdateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	productDb := NewProductDB(db)
	product, err := entity.NewProduct("Product 1", 10.5)
	assert.NoError(t, err)
	err = productDb.Create(product)
	assert.NoError(t, err)
	product, err = productDb.FindById(product.ID.String())
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)
	product.Name = "Product 2"
	err = productDb.Update(product)
	assert.NoError(t, err)
	product, err = productDb.FindById(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Product 2", product.Name)
}

// Copilot Create tests for Delete method
func TestDeleteProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	productDb := NewProductDB(db)
	product, err := entity.NewProduct("Product 1", 10.5)
	assert.NoError(t, err)
	err = productDb.Create(product)
	assert.NoError(t, err)
	product, err = productDb.FindById(product.ID.String())
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)
	err = productDb.Delete(product.ID.String())
	assert.NoError(t, err)
	product, err = productDb.FindById(product.ID.String())
	assert.Error(t, err)
	assert.Empty(t, product)
}
