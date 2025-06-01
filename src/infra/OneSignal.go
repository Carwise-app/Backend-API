package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type OneSignalRepository struct {
	Client     *http.Client
	AppID      string
	RestAPIKey string
	BaseURL    string
}
type OneSignalNotification struct {
	AppID                string            `json:"app_id"`
	IncludeAndroidRegIDs []string          `json:"include_android_reg_ids"`
	Contents             map[string]string `json:"contents"`
	Headings             map[string]string `json:"headings,omitempty"`
	Data                 map[string]string `json:"data,omitempty"`
}

type OneSignalResponse struct {
	ID         string   `json:"id"`
	Recipients int      `json:"recipients"`
	Errors     []string `json:"errors,omitempty"`
}

func NewOneSignalRepository() *OneSignalRepository {
	return &OneSignalRepository{
		Client:     &http.Client{},
		AppID:      os.Getenv("ONESIGNAL_APP_ID"),
		RestAPIKey: os.Getenv("ONESIGNAL_REST_API_KEY"),
		BaseURL:    "https://onesignal.com/api/v1/notifications",
	}
}

func (r *OneSignalRepository) PushNotification(deviceTokens []string, title, message string, customData map[string]string) error {
	url := r.BaseURL

	notification := OneSignalNotification{
		AppID:                r.AppID,
		IncludeAndroidRegIDs: deviceTokens,
		Contents: map[string]string{
			"en": message,
			"tr": message,
		},
	}

	if title != "" {
		notification.Headings = map[string]string{
			"en": title,
			"tr": title,
		}
	}

	if customData != nil {
		notification.Data = customData
	}

	jsonData, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("JSON marshal hatası: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("HTTP request oluşturma hatası: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+r.RestAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request hatası: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("response read error: %v", err)
	}

	var response OneSignalResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("JSON parse error: %v", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("OneSignal API hatası (%d): %s", resp.StatusCode, string(body))
	}

	if len(response.Errors) > 0 {
		return fmt.Errorf("OneSignal hataları: %v", response.Errors)
	}

	fmt.Printf("Bildirim başarıyla gönderildi! ID: %s, Alıcı sayısı: %d\n", response.ID, response.Recipients)
	return nil
}
