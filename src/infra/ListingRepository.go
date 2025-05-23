package infra

import (
	"carwise"
	"database/sql"
	"fmt"

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
		rear_bumper, created_by, created_at, updated_at, status
	)
	VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, 
		$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, 
		$33, $34, $35, $36, $37
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
		listing.CreatedBy, listing.CreatedAt, listing.UpdatedAt, listing.Status,
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
		   rear_bumper, created_by, created_at, updated_at, status	
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
		&listing.CreatedBy, &listing.CreatedAt, &listing.UpdatedAt, &listing.Status,
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
		   rear_bumper, created_by, created_at, updated_at, status
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
		&listing.CreatedBy, &listing.CreatedAt, &listing.UpdatedAt, &listing.Status,
	)
	if err != nil {
		return nil, err
	}

	return &listing, nil
}

func (r *ListingRepository) UpdateListing(listing *carwise.Listing) error {
	query := `
	UPDATE listings
	SET
		brand_id = $1,
		series_id = $2,
		model_id = $3,
		title = $4,
		description = $5,
		currency = $6,
		price = $7,
		city = $8,
		district = $9,
		neighborhood = $10,
		fuel_type = $11,
		transmission_type = $12,
		body_type = $13,
		drive_type = $14,
		engine_power = $15,
		engine_volume = $16,
		kilometers = $17,
		year = $18,
		color = $19,
		heavy_damage = $20,
		front_bumper = $21,
		front_hood = $22,
		roof = $23,
		front_right_door = $24,
		rear_right_door = $25,
		front_left_mudguard = $26,
		front_left_door = $27,
		rear_left_door = $28,
		rear_left_mudguard = $29,
		rear_bumper = $30,
		updated_at = $31,
		status = $32
	WHERE id = $33
	`

	_, err := r.db.Exec(query,
		listing.BrandId,
		listing.SeriesId,
		listing.ModelId,
		listing.Title,
		listing.Description,
		listing.Currency,
		listing.Price,
		listing.City,
		listing.District,
		listing.Neighborhood,
		listing.FuelType,
		listing.TransmissionType,
		listing.BodyType,
		listing.DriveType,
		listing.EnginePower,
		listing.EngineVolume,
		listing.Kilometers,
		listing.Year,
		listing.Color,
		listing.HeavyDamage,
		listing.FrontBumper,
		listing.FrontHood,
		listing.Roof,
		listing.FrontRightDoor,
		listing.RearRightDoor,
		listing.FrontLeftMudguard,
		listing.FrontLeftDoor,
		listing.RearLeftDoor,
		listing.RearLeftMudguard,
		listing.RearBumper,
		listing.UpdatedAt,
		listing.Status,
		listing.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ListingRepository) DeleteListing(id string) error {
	query := `
	DELETE FROM listings
	WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *ListingRepository) ListListing(filter *carwise.ListingFilter) ([]carwise.Listing, error) {
	query := `
	SELECT id, slug, brand_id, series_id, model_id, title, description, currency, price, 
			city, district, neighborhood, images, fuel_type, transmission_type, body_type, 
			drive_type, engine_power, engine_volume, kilometers, year, color, heavy_damage, 	
			created_by, created_at, updated_at, status
	FROM listings
	WHERE 1=1
	`
	params := []interface{}{}
	paramCount := 1

	if filter.Query != "" {
		query += fmt.Sprintf(` AND (title ILIKE $%d OR description ILIKE $%d)`, paramCount, paramCount)
		params = append(params, "%"+filter.Query+"%")
		paramCount++
	}

	if filter.BrandId != "" {
		query += fmt.Sprintf(` AND brand_id = $%d`, paramCount)
		params = append(params, filter.BrandId)
		paramCount++
	}

	if filter.SeriesId != "" {
		query += fmt.Sprintf(` AND series_id = $%d`, paramCount)
		params = append(params, filter.SeriesId)
		paramCount++
	}

	if filter.ModelId != "" {
		query += fmt.Sprintf(` AND model_id = $%d`, paramCount)
		params = append(params, filter.ModelId)
		paramCount++
	}

	if filter.BodyType != "" {
		query += fmt.Sprintf(` AND body_type = $%d`, paramCount)
		params = append(params, filter.BodyType)
		paramCount++
	}

	if filter.FuelType != "" {
		query += fmt.Sprintf(` AND fuel_type = $%d`, paramCount)
		params = append(params, filter.FuelType)
		paramCount++
	}

	if filter.TransmissionType != "" {
		query += fmt.Sprintf(` AND transmission_type = $%d`, paramCount)
		params = append(params, filter.TransmissionType)
		paramCount++
	}

	if filter.DriveType != "" {
		query += fmt.Sprintf(` AND drive_type = $%d`, paramCount)
		params = append(params, filter.DriveType)
		paramCount++
	}

	if filter.CreatedBy != "" {
		query += fmt.Sprintf(` AND created_by = $%d`, paramCount)
		params = append(params, filter.CreatedBy)
		paramCount++
	}

	if filter.Status != 0 {
		query += fmt.Sprintf(` AND status = $%d`, paramCount)
		params = append(params, filter.Status)
		paramCount++
	}

	if filter.MinPrice != 0 {
		query += fmt.Sprintf(` AND price >= $%d`, paramCount)
		params = append(params, filter.MinPrice)
		paramCount++
	}

	if filter.MaxPrice != 0 {
		query += fmt.Sprintf(` AND price <= $%d`, paramCount)
		params = append(params, filter.MaxPrice)
		paramCount++
	}

	if filter.MinYear != 0 {
		query += fmt.Sprintf(` AND year >= $%d`, paramCount)
		params = append(params, filter.MinYear)
		paramCount++
	}

	if filter.MaxYear != 0 {
		query += fmt.Sprintf(` AND year <= $%d`, paramCount)
		params = append(params, filter.MaxYear)
		paramCount++
	}

	if filter.MinKilometers != 0 {
		query += fmt.Sprintf(` AND kilometers >= $%d`, paramCount)
		params = append(params, filter.MinKilometers)
		paramCount++
	}

	if filter.MaxKilometers != 0 {
		query += fmt.Sprintf(` AND kilometers <= $%d`, paramCount)
		params = append(params, filter.MaxKilometers)
		paramCount++
	}

	if filter.City != "" {
		query += fmt.Sprintf(` AND city = $%d`, paramCount)
		params = append(params, filter.City)
		paramCount++
	}

	if filter.District != "" {
		query += fmt.Sprintf(` AND district = $%d`, paramCount)
		params = append(params, filter.District)
		paramCount++
	}

	if filter.Neighborhood != "" {
		query += fmt.Sprintf(` AND neighborhood = $%d`, paramCount)
		params = append(params, filter.Neighborhood)
		paramCount++
	}

	if filter.MinEnginePower != 0 {
		query += fmt.Sprintf(` AND engine_power >= $%d`, paramCount)
		params = append(params, filter.MinEnginePower)
		paramCount++
	}

	if filter.MaxEnginePower != 0 {
		query += fmt.Sprintf(` AND engine_power <= $%d`, paramCount)
		params = append(params, filter.MaxEnginePower)
		paramCount++
	}

	if filter.Color != "" {
		query += fmt.Sprintf(` AND color = $%d`, paramCount)
		params = append(params, filter.Color)
		paramCount++
	}

	if filter.HeavyDamage != nil {
		query += fmt.Sprintf(` AND heavy_damage = $%d`, paramCount)
		params = append(params, *filter.HeavyDamage)
		paramCount++
	}

	if filter.SortBy != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", filter.SortBy, filter.SortOrder)
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	params = append(params, filter.Limit, (filter.Page-1)*filter.Limit)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve findings: %w", err)
	}
	defer rows.Close()

	listings := []carwise.Listing{}
	for rows.Next() {
		var listing carwise.Listing
		err = rows.Scan(&listing.Id, &listing.Slug, &listing.BrandId, &listing.SeriesId, &listing.ModelId,
			&listing.Title, &listing.Description, &listing.Currency, &listing.Price, &listing.City,
			&listing.District, &listing.Neighborhood, pq.Array(&listing.Images), &listing.FuelType,
			&listing.TransmissionType, &listing.BodyType, &listing.DriveType, &listing.EnginePower,
			&listing.EngineVolume, &listing.Kilometers, &listing.Year, &listing.Color, &listing.HeavyDamage,
			&listing.CreatedBy, &listing.CreatedAt, &listing.UpdatedAt, &listing.Status,
		)
		if err != nil {
			return nil, err
		}
		listings = append(listings, listing)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return listings, nil
}

func (r *ListingRepository) CountListing(filter *carwise.ListingFilter) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM listings
	WHERE 1=1
	`
	params := []interface{}{}
	paramCount := 1

	if filter.Query != "" {
		query += fmt.Sprintf(` AND (title ILIKE $%d OR description ILIKE $%d)`, paramCount, paramCount)
		params = append(params, "%"+filter.Query+"%")
		paramCount++
	}

	if filter.BrandId != "" {
		query += fmt.Sprintf(` AND brand_id = $%d`, paramCount)
		params = append(params, filter.BrandId)
		paramCount++
	}

	if filter.SeriesId != "" {
		query += fmt.Sprintf(` AND series_id = $%d`, paramCount)
		params = append(params, filter.SeriesId)
		paramCount++
	}

	if filter.ModelId != "" {
		query += fmt.Sprintf(` AND model_id = $%d`, paramCount)
		params = append(params, filter.ModelId)
		paramCount++
	}

	if filter.BodyType != "" {
		query += fmt.Sprintf(` AND body_type = $%d`, paramCount)
		params = append(params, filter.BodyType)
		paramCount++
	}

	if filter.FuelType != "" {
		query += fmt.Sprintf(` AND fuel_type = $%d`, paramCount)
		params = append(params, filter.FuelType)
		paramCount++
	}

	if filter.TransmissionType != "" {
		query += fmt.Sprintf(` AND transmission_type = $%d`, paramCount)
		params = append(params, filter.TransmissionType)
		paramCount++
	}

	if filter.DriveType != "" {
		query += fmt.Sprintf(` AND drive_type = $%d`, paramCount)
		params = append(params, filter.DriveType)
		paramCount++
	}

	if filter.CreatedBy != "" {
		query += fmt.Sprintf(` AND created_by = $%d`, paramCount)
		params = append(params, filter.CreatedBy)
		paramCount++
	}

	if filter.Status != 0 {
		query += fmt.Sprintf(` AND status = $%d`, paramCount)
		params = append(params, filter.Status)
		paramCount++
	}

	if filter.MinPrice != 0 {
		query += fmt.Sprintf(` AND price >= $%d`, paramCount)
		params = append(params, filter.MinPrice)
		paramCount++
	}

	if filter.MaxPrice != 0 {
		query += fmt.Sprintf(` AND price <= $%d`, paramCount)
		params = append(params, filter.MaxPrice)
		paramCount++
	}

	if filter.MinYear != 0 {
		query += fmt.Sprintf(` AND year >= $%d`, paramCount)
		params = append(params, filter.MinYear)
		paramCount++
	}

	if filter.MaxYear != 0 {
		query += fmt.Sprintf(` AND year <= $%d`, paramCount)
		params = append(params, filter.MaxYear)
		paramCount++
	}

	if filter.MinKilometers != 0 {
		query += fmt.Sprintf(` AND kilometers >= $%d`, paramCount)
		params = append(params, filter.MinKilometers)
		paramCount++
	}

	if filter.MaxKilometers != 0 {
		query += fmt.Sprintf(` AND kilometers <= $%d`, paramCount)
		params = append(params, filter.MaxKilometers)
		paramCount++
	}

	if filter.City != "" {
		query += fmt.Sprintf(` AND city = $%d`, paramCount)
		params = append(params, filter.City)
		paramCount++
	}

	if filter.District != "" {
		query += fmt.Sprintf(` AND district = $%d`, paramCount)
		params = append(params, filter.District)
		paramCount++
	}

	if filter.Neighborhood != "" {
		query += fmt.Sprintf(` AND neighborhood = $%d`, paramCount)
		params = append(params, filter.Neighborhood)
		paramCount++
	}

	if filter.MinEnginePower != 0 {
		query += fmt.Sprintf(` AND engine_power >= $%d`, paramCount)
		params = append(params, filter.MinEnginePower)
		paramCount++
	}

	if filter.MaxEnginePower != 0 {
		query += fmt.Sprintf(` AND engine_power <= $%d`, paramCount)
		params = append(params, filter.MaxEnginePower)
		paramCount++
	}

	if filter.Color != "" {
		query += fmt.Sprintf(` AND color = $%d`, paramCount)
		params = append(params, filter.Color)
		paramCount++
	}

	if filter.HeavyDamage != nil {
		query += fmt.Sprintf(` AND heavy_damage = $%d`, paramCount)
		params = append(params, *filter.HeavyDamage)
		paramCount++
	}

	var count int
	err := r.db.QueryRow(query, params...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
