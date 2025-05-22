package infra

import (
	"carwise"
	"database/sql"

	"github.com/lib/pq"
)

type ListingRepository struct {
	db *sql.DB
}

func NewListingRepository() *ListingRepository {
	db := ConnectDb()
	return &ListingRepository{
		db: db,
	}
}

func (r *ListingRepository) CreateListing(listing *carwise.Listing) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
	INSERT INTO listings (
		id, slug, brand_id, series_id, model_id, title, description, currency, price, 
		city, district, neighborhood, images, fuel_type, transmission_type, body_type, 
		drive_type, engine_power, engine_volume, kilometers, year, color, heavy_damage, 
		front_bumper, front_hood, roof, front_right_door, rear_right_door, 
		front_left_mudguard, front_left_door, rear_left_door, rear_left_mudguard, 
		rear_bumper, created_by, created_at, updated_at
	)
	VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, 
		$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, 
		$33, $34, $35, $36
	)
	RETURNING id`

	err = tx.QueryRow(query,
		listing.Id, listing.Slug, listing.BrandId, listing.SeriesId, listing.ModelId,
		listing.Title, listing.Description, listing.Currency, listing.Price,
		listing.City, listing.District, listing.Neighborhood, pq.Array(listing.Images),
		listing.FuelType, listing.TransmissionType, listing.BodyType, listing.DriveType,
		listing.EnginePower, listing.EngineVolume, listing.Kilometers, listing.Year,
		listing.Color, listing.HeavyDamage, listing.FrontBumper, listing.FrontHood,
		listing.Roof, listing.FrontRightDoor, listing.RearRightDoor, listing.FrontLeftMudguard,
		listing.FrontLeftDoor, listing.RearLeftDoor, listing.RearLeftMudguard, listing.RearBumper,
		listing.CreatedBy, listing.CreatedAt, listing.UpdatedAt,
	).Scan(&listing.Id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ListingRepository) GetListingById(id string) (*carwise.Listing, error) {
	query := `
	SELECT id, slug, brand_id, series_id, model_id, title, description, currency, price, 
		   city, district, neighborhood, images, fuel_type, transmission_type, body_type, 
		   drive_type, engine_power, engine_volume, kilometers, year, color, heavy_damage, 
		   front_bumper, front_hood, roof, front_right_door, rear_right_door, 
		   front_left_mudguard, front_left_door, rear_left_door, rear_left_mudguard, 
		   rear_bumper, created_by, created_at, updated_at
	FROM listings
	WHERE id = $1
	`

	var listing carwise.Listing
	err := r.db.QueryRow(query, id).Scan(
		&listing.Id, &listing.Slug, &listing.BrandId, &listing.SeriesId, &listing.ModelId,
		&listing.Title, &listing.Description, &listing.Currency, &listing.Price,
		&listing.City, &listing.District, &listing.Neighborhood, pq.Array(&listing.Images),
		&listing.FuelType, &listing.TransmissionType, &listing.BodyType, &listing.DriveType,
		&listing.EnginePower, &listing.EngineVolume, &listing.Kilometers, &listing.Year,
		&listing.Color, &listing.HeavyDamage, &listing.FrontBumper, &listing.FrontHood,
		&listing.Roof, &listing.FrontRightDoor, &listing.RearRightDoor, &listing.FrontLeftMudguard,
		&listing.FrontLeftDoor, &listing.RearLeftDoor, &listing.RearLeftMudguard, &listing.RearBumper,
		&listing.CreatedBy, &listing.CreatedAt, &listing.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &listing, nil
}

func (r *ListingRepository) GetListingBySlug(slug string) (*carwise.Listing, error) {
	query := `
	SELECT id, slug, brand_id, series_id, model_id, title, description, currency, price, 
		   city, district, neighborhood, images, fuel_type, transmission_type, body_type, 
		   drive_type, engine_power, engine_volume, kilometers, year, color, heavy_damage, 
		   front_bumper, front_hood, roof, front_right_door, rear_right_door, 
		   front_left_mudguard, front_left_door, rear_left_door, rear_left_mudguard, 
		   rear_bumper, created_by, created_at, updated_at
	FROM listings
	WHERE slug = $1
	`

	var listing carwise.Listing
	err := r.db.QueryRow(query, slug).Scan(
		&listing.Id, &listing.Slug, &listing.BrandId, &listing.SeriesId, &listing.ModelId,
		&listing.Title, &listing.Description, &listing.Currency, &listing.Price,
		&listing.City, &listing.District, &listing.Neighborhood, pq.Array(&listing.Images),
		&listing.FuelType, &listing.TransmissionType, &listing.BodyType, &listing.DriveType,
		&listing.EnginePower, &listing.EngineVolume, &listing.Kilometers, &listing.Year,
		&listing.Color, &listing.HeavyDamage, &listing.FrontBumper, &listing.FrontHood,
		&listing.Roof, &listing.FrontRightDoor, &listing.RearRightDoor, &listing.FrontLeftMudguard,
		&listing.FrontLeftDoor, &listing.RearLeftDoor, &listing.RearLeftMudguard, &listing.RearBumper,
		&listing.CreatedBy, &listing.CreatedAt, &listing.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &listing, nil
}
