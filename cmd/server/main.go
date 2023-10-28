package main

import (
	"fmt"
	"net/http"

	"github.com/brenoproti/go-api/configs"
	"github.com/brenoproti/go-api/internal/entity"
	"github.com/brenoproti/go-api/internal/infra/database"
	"github.com/brenoproti/go-api/internal/infra/webserver/handlers"
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

	productDb := database.NewProductDB(db)
	productHandler := handlers.NewProductHandler(productDb)

	http.HandleFunc("/products", productHandler.Create)
	http.ListenAndServe(":8000", nil)
}
