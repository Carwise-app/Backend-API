package carwise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type PredictionResponse struct {
	TahminiFiyat float64 `json:"tahmini_fiyat"`
	R2Skoru      float64 `json:"r2_skoru"`
	MAE          float64 `json:"mae"`
}

func (i *Interactor) CreatePredict(request *PredictRequest) (*PredictionResponse, error) {
	apiURL := os.Getenv("CAR_PRICE_PREDICTION_API_URL")
	if apiURL == "" {
		return nil, fmt.Errorf("CAR_PRICE_PREDICTION_API_URL environment variable not set")
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to make prediction request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prediction API returned non-200 status code: %d", resp.StatusCode)
	}

	var predictionResp PredictionResponse
	if err := json.NewDecoder(resp.Body).Decode(&predictionResp); err != nil {
		return nil, fmt.Errorf("failed to decode prediction response: %v", err)
	}

	if request.UserId != "" {
		predict := &Predict{
			Brand:            request.Brand,
			Series:           request.Series,
			Model:            request.Model,
			Year:             request.Year,
			Mileage:          request.Mileage,
			EngineVolume:     request.EngineVolume,
			EnginePower:      request.EnginePower,
			AccidentHistory:  request.AccidentHistory,
			TransmissionType: request.TransmissionType,
			FuelType:         request.FuelType,
			BodyType:         request.BodyType,
			Color:            request.Color,
			OriginalParts:    request.OriginalParts,
			ReplacedParts:    request.ReplacedParts,
			PaintedParts:     request.PaintedParts,
			CreatedBy:        request.UserId,
			CreatedAt:        time.Now().Unix(),
			Price:            predictionResp.TahminiFiyat,
			R2:               predictionResp.R2Skoru,
			MAE:              predictionResp.MAE,
		}

		err = i.services.PredictRepo.SavePredict(predict)
		if err != nil {
			return nil, err
		}
	}

	return &predictionResp, nil
}

func (i *Interactor) GetPredicts(request *GetPredictsRequest) (*GetPredictsResponse, error) {
	predicts, err := i.services.PredictRepo.GetPredictByUserId(request.UserId, request.Page, request.Limit)
	if err != nil {
		return nil, err
	}
	total, err := i.services.PredictRepo.CountPredictByUserId(request.UserId)
	if err != nil {
		return nil, err
	}
	return &GetPredictsResponse{
		Predicts: predicts,
		Total:    total,
	}, nil
}
