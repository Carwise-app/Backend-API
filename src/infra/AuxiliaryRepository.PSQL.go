package infra

import (
	"carwise"
	"database/sql"
)

type AuxiliaryRepository struct {
	db *sql.DB
}

func NewAuxiliaryRepository() *AuxiliaryRepository {
	database := ConnectDb()
	return &AuxiliaryRepository{db: database}
}

func (repo *AuxiliaryRepository) GetBrands() ([]carwise.Brand, error) {
	rows, err := repo.db.Query("SELECT id, logo, name FROM brands")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brands []carwise.Brand
	for rows.Next() {
		var brand carwise.Brand
		if err := rows.Scan(&brand.ID, &brand.Name); err != nil {
			return nil, err
		}
		brands = append(brands, brand)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return brands, nil
}

func (repo *AuxiliaryRepository) GetSeriesByBrand(brandID int) ([]carwise.Series, error) {
	rows, err := repo.db.Query("SELECT id, brand_id, name FROM series WHERE brand_id = $1", brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var series []carwise.Series
	for rows.Next() {
		var s carwise.Series
		if err := rows.Scan(&s.ID, &s.BrandID, &s.Name); err != nil {
			return nil, err
		}
		series = append(series, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return series, nil
}

func (repo *AuxiliaryRepository) GetModelsBySeries(seriesID int) ([]carwise.Model, error) {
	rows, err := repo.db.Query("SELECT id, series_id, name FROM models WHERE series_id = $1", seriesID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []carwise.Model
	for rows.Next() {
		var model carwise.Model
		if err := rows.Scan(&model.ID, &model.SeriesID, &model.Name); err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return models, nil
}
