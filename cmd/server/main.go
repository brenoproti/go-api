package main

import (
	"fmt"
	"net/http"

	"github.com/brenoproti/go-api/configs"
	"github.com/brenoproti/go-api/internal/entity"
	"github.com/brenoproti/go-api/internal/infra/database"
	"github.com/brenoproti/go-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"

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

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	productDb := database.NewProductDB(db)
	productHandler := handlers.NewProductHandler(productDb)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/", productHandler.Create)
		r.Put("/{id}", productHandler.Update)
		r.Get("/{id}", productHandler.FindById)
		r.Get("/", productHandler.GetProducts)
		r.Delete("/{id}", productHandler.Delete)
	})

	userDb := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDb, config.TokenAuth, config.JWTExpiresIn)

	r.Post("/users", userHandler.Create)
	r.Post("/users/generate_token", userHandler.GetJWT)

	http.ListenAndServe(":8000", r)
}
