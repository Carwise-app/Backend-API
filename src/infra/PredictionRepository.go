package infra

import (
	"database/sql"

	"carwise"
)

type PredictionRepository struct {
	db *sql.DB
}

func NewPredictionRepository() *PredictionRepository {
	return &PredictionRepository{db: ConnectDb()}
}

func (r *PredictionRepository) SaveImagePrediction(prediction *carwise.ImagePrediction) error {
	query := `
		INSERT INTO image_predictions (image_id, prediction, confidence)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, prediction.ImageId, prediction.Prediction, prediction.Confidence)
	return err
}

func (r *PredictionRepository) GetImagePredictionsByImageId(imageId string) (*carwise.ImagePrediction, error) {
	query := `
		SELECT image_id, prediction, confidence, created_at FROM image_predictions WHERE image_id = $1
	`
	rows, err := r.db.Query(query, imageId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prediction carwise.ImagePrediction
	for rows.Next() {
		err := rows.Scan(&prediction.ImageId, &prediction.Prediction, &prediction.Confidence, &prediction.CreatedAt)
		if err != nil {
			return nil, err
		}
	}
	return &prediction, nil
}
