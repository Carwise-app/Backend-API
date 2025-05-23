package carwise

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (i *Interactor) UploadImage(request *UploadImageRequest) (*Image, error) {
	if request.File == nil {
		return nil, errors.New("no file provided")
	}

	if request.File.Size > 5*1024*1024 {
		return nil, errors.New("file size exceeds 5MB limit")
	}

	contentType := request.File.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		return nil, errors.New("only jpeg, png and gif images are allowed")
	}

	image, err := i.services.ImageRepo.SaveImage(request.File, request.UserId)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (i *Interactor) DeleteImage(request *DeleteImageRequest) error {
	image, err := i.services.ImageRepo.GetImageById(request.ImageId)
	if err != nil {
		return err
	}

	if request.Role != 2 && image.CreatedBy != request.UserId {
		return errors.New("unauthorized to delete this image")
	}

	return i.services.ImageRepo.DeleteImage(request.ImageId)
}

func (i *Interactor) PredictImage(request *PredictImageRequest) (*PredictImageResponse, error) {
	image, err := i.services.ImageRepo.GetImageById(request.ImageId)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %v", err)
	}

	existingPrediction, _ := i.services.PredictionRepo.GetImagePredictionsByImageId(image.Id)
	log.Println(existingPrediction)
	if existingPrediction.CreatedAt != 0 {
		response := &PredictImageResponse{
			Image:      *image,
			Prediction: existingPrediction.Prediction,
			Confidence: existingPrediction.Confidence,
			CreatedAt:  existingPrediction.CreatedAt,
		}
		return response, nil
	}

	// Resim dosyasını aç
	file, err := os.Open(image.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %v", err)
	}
	defer file.Close()

	// Multipart form oluştur
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Dosyayı form'a ekle
	part, err := writer.CreateFormFile("file", filepath.Base(image.Path))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}

	// Dosya içeriğini kopyala
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %v", err)
	}

	// Form'u kapat
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %v", err)
	}

	// HTTP isteği oluştur
	req, err := http.NewRequest("POST", os.Getenv("CAR_RECOGNITION_API_URL"), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Content-Type header'ını ayarla
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// İsteği gönder
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Yanıt gövdesini oku
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Hata durumunu kontrol et
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prediction service returned status code %d: %s", resp.StatusCode, string(respBody))
	}

	// Yanıtı parse et
	var result struct {
		Prediction bool    `json:"prediction"`
		Confidence float64 `json:"confidence"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode prediction response: %v", err)
	}

	// Tahmini kaydet
	prediction := &ImagePrediction{
		ImageId:    image.Id,
		Prediction: result.Prediction,
		Confidence: result.Confidence,
		CreatedAt:  time.Now().Unix(),
	}

	if err := i.services.PredictionRepo.SaveImagePrediction(prediction); err != nil {
		return nil, fmt.Errorf("failed to save prediction: %v", err)
	}

	// Yanıtı hazırla
	response := &PredictImageResponse{
		Image:      *image,
		Prediction: result.Prediction,
		Confidence: result.Confidence,
		CreatedAt:  prediction.CreatedAt,
	}

	log.Printf("Successfully processed prediction for image %s: car=%v, confidence=%.2f",
		image.Id, result.Prediction, result.Confidence)

	return response, nil
}
