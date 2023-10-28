package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brenoproti/go-api/configs"
	"github.com/brenoproti/go-api/internal/dto"
	"github.com/brenoproti/go-api/internal/entity"
	"github.com/brenoproti/go-api/internal/infra/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	config := configs.LoadConfig("cmd/server")
	fmt.Printf("%v", config)
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{}, &entity.Product{})

	http.ListenAndServe(":8000", nil)
}

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product dto.ProductDTO
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	entity, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.ProductDB.Create(entity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
