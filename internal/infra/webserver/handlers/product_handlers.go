package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/brenoproti/go-api/internal/dto"
	"github.com/brenoproti/go-api/internal/entity"
	"github.com/brenoproti/go-api/internal/infra/database"
	pkg "github.com/brenoproti/go-api/pkg/entity"
	"github.com/go-chi/chi"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// Create product godoc
// @Summary Create a new product
// @Description Create a new product
// @Tags products
// @Accept  json
// @Produce  json
// @Param request body dto.ProductDTO true "Product info"
// @Success 201 {string} string	"Product created"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string	"Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /products [post]
// @Security ApiKeyAuth
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
	w.Header().Add("id", entity.ID.String())
	w.WriteHeader(http.StatusCreated)
}

// FindById product godoc
// @Summary Find a product by ID
// @Description Find a product by ID
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductDTO
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /products/{id} [get]
// @Security ApiKeyAuth
func (p *ProductHandler) FindById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := p.ProductDB.FindById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// Update product godoc
// @Summary Update a product
// @Description Update a product
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Param request body dto.ProductDTO true "Product info"
// @Success 200 {string} string	"Product updated"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /products/{id} [put]
// @Security ApiKeyAuth
func (p *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	id := chi.URLParam(r, "id")
	product.ID, err = pkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = p.ProductDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Delete product godoc
// @Summary Delete a product
// @Description Delete a product
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {string} string	"Product deleted"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /products/{id} [delete]
// @Security ApiKeyAuth
func (p *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := p.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetProducts godoc
// @Summary Get products paginated
// @Description Get products paginated
// @Tags products
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param sort query string false "Sort by field"
// @Success 200 {array} string
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /products [get]
// @Security ApiKeyAuth
func (p *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	queryPage := r.URL.Query().Get("page")
	queryLimit := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")
	page, err := strconv.Atoi(queryPage)
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(queryLimit)
	if err != nil {
		limit = 0
	}
	products, err := p.ProductDB.FindAll(page, limit, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
