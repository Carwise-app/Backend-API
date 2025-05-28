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

func (r *BrandRepository) GetAll() ([]carwise.Brand, error) {
	query := `
	SELECT id, image_path, name FROM brands
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brands []carwise.Brand
	for rows.Next() {
		var brand carwise.Brand
		err := rows.Scan(&brand.Id, &brand.ImagePath, &brand.Name)
		if err != nil {
			return nil, err
		}
		brands = append(brands, brand)
	}

	return brands, nil
}

func (r *BrandRepository) Create(brand *carwise.Brand) error {
	query := `
	INSERT INTO brands (id, image_path, name)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	err := r.db.QueryRow(query, brand.Id, brand.ImagePath, brand.Name).Scan(&brand.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *BrandRepository) Update(brand *carwise.Brand) error {
	query := `
	UPDATE brands
	SET image_path = $1, name = $2
	WHERE id = $3
	RETURNING id
	`

	err := r.db.QueryRow(query, brand.ImagePath, brand.Name, brand.Id).Scan(&brand.Id)
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
	SELECT id, image_path, name FROM brands WHERE id = $1
	`

	var brand carwise.Brand
	err := r.db.QueryRow(query, id).Scan(&brand.Id, &brand.ImagePath, &brand.Name)
	if err != nil {
		return nil, err
	}
	return &brand, nil
}

func (r *BrandRepository) GetAllSeriesByBrandId(brandId string) ([]carwise.Series, error) {
	query := `
	SELECT id, brand_id, name FROM series WHERE brand_id = $1
	`

	rows, err := r.db.Query(query, brandId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var series []carwise.Series
	for rows.Next() {
		var s carwise.Series
		err := rows.Scan(&s.Id, &s.BrandId, &s.Name)
		if err != nil {
			return nil, err
		}
		series = append(series, s)
	}

	return series, nil
}

func (r *BrandRepository) CreateSeries(series *carwise.Series) error {
	query := `
	INSERT INTO series (id, brand_id, name)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	err := r.db.QueryRow(query, series.Id, series.BrandId, series.Name).Scan(&series.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *BrandRepository) UpdateSeries(series *carwise.Series) error {
	query := `
	UPDATE series
	SET name = $1
	WHERE id = $2
	RETURNING id
	`

	err := r.db.QueryRow(query, series.Name, series.Id).Scan(&series.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *BrandRepository) DeleteSeries(id string) error {
	query := `
	DELETE FROM series WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *BrandRepository) GetSeriesById(id string) (*carwise.Series, error) {
	query := `
	SELECT id, brand_id, name FROM series WHERE id = $1
	`

	var series carwise.Series
	err := r.db.QueryRow(query, id).Scan(&series.Id, &series.BrandId, &series.Name)
	if err != nil {
		return nil, err
	}
	return &series, nil
}

func (r *BrandRepository) GetAllModelsBySeriesId(seriesId string) ([]carwise.Model, error) {
	query := `
	SELECT id, brand_id, series_id, name FROM models WHERE series_id = $1
	`

	rows, err := r.db.Query(query, seriesId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []carwise.Model
	for rows.Next() {
		var m carwise.Model
		err := rows.Scan(&m.Id, &m.BrandId, &m.SeriesId, &m.Name)
		if err != nil {
			return nil, err
		}
		models = append(models, m)
	}

	return models, nil
}

func (r *BrandRepository) CreateModel(model *carwise.Model) error {
	query := `
	INSERT INTO models (id, brand_id, series_id, name)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	err := r.db.QueryRow(query, model.Id, model.BrandId, model.SeriesId, model.Name).Scan(&model.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *BrandRepository) UpdateModel(model *carwise.Model) error {
	query := `
	UPDATE models
	SET name = $1
	WHERE id = $2
	RETURNING id
	`

	err := r.db.QueryRow(query, model.Name, model.Id).Scan(&model.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *BrandRepository) DeleteModel(id string) error {
	query := `
	DELETE FROM models WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *BrandRepository) GetModelById(id string) (*carwise.Model, error) {
	query := `
	SELECT id, brand_id, series_id, name FROM models WHERE id = $1
	`

	var model carwise.Model
	err := r.db.QueryRow(query, id).Scan(&model.Id, &model.BrandId, &model.SeriesId, &model.Name)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *BrandRepository) GetAllWithDetails() ([]carwise.BrandWithDetails, error) {
	// First get all brands
	brands, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	var brandsWithDetails []carwise.BrandWithDetails
	for _, brand := range brands {
		// Get series for this brand
		series, err := r.GetAllSeriesByBrandId(brand.Id)
		if err != nil {
			return nil, err
		}

		// For each series, get its models
		var seriesDetails []carwise.SeriesDetail
		for _, s := range series {
			models, err := r.GetAllModelsBySeriesId(s.Id)
			if err != nil {
				return nil, err
			}

			seriesDetail := carwise.SeriesDetail{
				Id:     s.Id,
				Name:   s.Name,
				Models: models,
			}
			seriesDetails = append(seriesDetails, seriesDetail)
		}

		brandWithDetails := carwise.BrandWithDetails{
			Id:        brand.Id,
			ImagePath: brand.ImagePath,
			Name:      brand.Name,
			Series:    seriesDetails,
		}
		brandsWithDetails = append(brandsWithDetails, brandWithDetails)
	}

	return brandsWithDetails, nil
}
