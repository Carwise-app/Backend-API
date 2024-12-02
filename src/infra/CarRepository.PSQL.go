package infra

import (
	"carwise"
	"database/sql"
	"fmt"
)

type CarRepository struct {
	db *sql.DB
}

func NewCarRepository() *CarRepository {
	database := ConnectDb()
	return &CarRepository{db: database}
}

func (r *CarRepository) Create(car *carwise.Car) error {
	query := `
		INSERT INTO cars (
			id, 
			owner_id, 
			title, 
			description, 
			currency, 
			price, 
			city, 
			district, 
			neighborhood, 
			listing_number, 
			listing_date, 
			brand_id, 
			series_id, 
			model_id, 
			year, 
			fuel_type, 
			transmission, 
			mileage, 
			body_type, 
			engine_power, 
			engine_volume, 
			drive_type, 
			color, 
			warranty, 
			heavy_damage, 
			seller_type, 
			trade_option, 
			front_bumper, 
			front_hood, 
			roof, 
			front_right_door, 
			rear_right_door, 
			front_left_mudguard, 
			front_left_door, 
			rear_left_door, 
			rear_left_mudguard, 
			rear_bumper
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37
		)`
	_, err := r.db.Exec(query,
		car.ID,
		car.OwnerId,
		car.Title,
		car.Description,
		car.Currency,
		car.Price,
		car.City,
		car.District,
		car.Neighborhood,
		car.ListingNumber,
		car.ListingDate,
		car.BrandId,
		car.SeriesId,
		car.ModelId,
		car.Year,
		car.FuelType,
		car.Transmission,
		car.Mileage,
		car.BodyType,
		car.EnginePower,
		car.EngineVolume,
		car.DriveType,
		car.Color,
		car.Warranty,
		car.HeavyDamage,
		car.SellerType,
		car.TradeOption,
		car.FrontBumper,
		car.FrontHood,
		car.Roof,
		car.FrontRightDoor,
		car.RearRightDoor,
		car.FrontLeftMudguard,
		car.FrontLeftDoor,
		car.RearLeftDoor,
		car.RearLeftMudguard,
		car.RearBumper,
	)
	if err != nil {
		return fmt.Errorf("failed to create car: %w", err)
	}
	return nil
}
