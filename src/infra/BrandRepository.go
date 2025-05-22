package infra

import (
	"carwise"
	"database/sql"
)

type BrandRepository struct {
	db *sql.DB
}

func NewBrandRepository() *BrandRepository {
	db := ConnectDb()
	return &BrandRepository{db: db}
}

func (r *BrandRepository) Create(brand *carwise.Brand) error {
	query := `
	INSERT INTO brands (id, image_id, name)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	err := r.db.QueryRow(query, brand.Id, brand.ImageId, brand.Name).Scan(&brand.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *BrandRepository) Update(brand *carwise.Brand) error {
	query := `
	UPDATE brands
	SET image_id = $1, name = $2
	WHERE id = $3
	RETURNING id
	`

	err := r.db.QueryRow(query, brand.ImageId, brand.Name, brand.Id).Scan(&brand.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *BrandRepository) Delete(id string) error {
	query := `
	DELETE FROM brands WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *BrandRepository) GetById(id string) (*carwise.Brand, error) {
	query := `
	SELECT id, image_id, name FROM brands WHERE id = $1
	`

	var brand carwise.Brand
	err := r.db.QueryRow(query, id).Scan(&brand.Id, &brand.ImageId, &brand.Name)
	if err != nil {
		return nil, err
	}
	return &brand, nil
}
