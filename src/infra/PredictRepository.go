package infra

import (
	"carwise"
	"database/sql"
)

type PredictRepository struct {
	db *sql.DB
}

func NewPredictRepository() *PredictRepository {
	db := ConnectDb()
	return &PredictRepository{db: db}
}

func (r *PredictRepository) SavePredict(predict *carwise.Predict) error {
	query := `
		INSERT INTO predicts (created_by, created_at, brand, series, model, year, 
		mileage, engine_volume, engine_power, accident_history, transmission_type, fuel_type, 
		body_type, color, original_parts, replaced_parts, painted_parts, price, r2, mae)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	`
	_, err := r.db.Exec(query, predict.CreatedBy, predict.CreatedAt, predict.Brand,
		predict.Series, predict.Model, predict.Year, predict.Mileage, predict.EngineVolume,
		predict.EnginePower, predict.AccidentHistory, predict.TransmissionType, predict.FuelType,
		predict.BodyType, predict.Color, predict.OriginalParts, predict.ReplacedParts, predict.PaintedParts,
		predict.Price, predict.R2, predict.MAE)
	if err != nil {
		return err
	}
	return nil
}

func (r *PredictRepository) GetPredictByUserId(userId string, page, limit int) ([]carwise.Predict, error) {
	query := `
		SELECT * FROM predicts WHERE created_by = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, userId, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var predicts []carwise.Predict
	for rows.Next() {
		var predict carwise.Predict
		err := rows.Scan(&predict.Id, &predict.CreatedBy, &predict.CreatedAt, &predict.Brand, &predict.Series,
			&predict.Model, &predict.Year, &predict.Mileage, &predict.EngineVolume, &predict.EnginePower,
			&predict.AccidentHistory, &predict.TransmissionType, &predict.FuelType, &predict.BodyType,
			&predict.Color, &predict.OriginalParts, &predict.ReplacedParts, &predict.PaintedParts,
			&predict.Price, &predict.R2, &predict.MAE)
		if err != nil {
			return nil, err
		}
		predicts = append(predicts, predict)
	}
	return predicts, nil
}

func (r *PredictRepository) CountPredictByUserId(userId string) (int, error) {
	query := `
		SELECT COUNT(*) FROM predicts WHERE created_by = $1
	`
	var count int
	err := r.db.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
