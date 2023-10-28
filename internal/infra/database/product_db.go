package database

import (
	"fmt"

	"github.com/brenoproti/go-api/internal/entity"
	"gorm.io/gorm"
)

type ProductDB struct {
	DB *gorm.DB
}

func NewProductDB(db *gorm.DB) *ProductDB {
	return &ProductDB{
		DB: db,
	}
}

func (p *ProductDB) Create(product *entity.Product) error {
	return p.DB.Create(product).Error
}

func (p *ProductDB) FindById(id string) (*entity.Product, error) {
	var product entity.Product
	err := p.DB.First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *ProductDB) FindAll(page, limit int, sort string) ([]entity.Product, error) {
	var products []entity.Product
	var err error
	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	if page != 0 && limit != 0 {
		offset := (page - 1) * limit
		err = p.DB.Offset(offset).Limit(limit).Order(fmt.Sprintf("created_at %s", sort)).Find(&products).Error
	} else {
		err = p.DB.Order(fmt.Sprintf("created_at %s", sort)).Find(&products).Error
	}
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *ProductDB) Update(product *entity.Product) error {
	_, err := p.FindById(product.ID.String())
	if err != nil {
		return err
	}
	return p.DB.Save(product).Error
}

func (p *ProductDB) Delete(id string) error {
	product, err := p.FindById(id)
	if err != nil {
		return err
	}
	return p.DB.Delete(product).Error
}
