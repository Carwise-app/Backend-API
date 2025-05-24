package infra

import (
	"carwise"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

type GoogleAuth struct {
	oauth2 oauth2.Config
	state  string
}

func NewGoogleAuth() *GoogleAuth {
	oauth2Config := oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       strings.Split(os.Getenv("GOOGLE_SCOPES"), ","),
		Endpoint:     google.Endpoint,
	}

	return &GoogleAuth{
		oauth2: oauth2Config,
		state:  os.Getenv("OAUTH2_STATE"),
	}
}

func (g *GoogleAuth) VerifyIDToken(idTokenStr string) (*carwise.GoogleResponse, error) {
	payload, err := idtoken.Validate(context.Background(), idTokenStr, g.oauth2.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %v", err)
	}
	return &carwise.GoogleResponse{
		Id:        payload.Claims["sub"].(string),
		Email:     payload.Claims["email"].(string),
		FirstName: payload.Claims["given_name"].(string),
		LastName:  payload.Claims["family_name"].(string),
	}, nil
}

func (g *GoogleAuth) GetUrl() string {
	return g.oauth2.AuthCodeURL(g.state, oauth2.AccessTypeOffline)
}

func (g *GoogleAuth) Callback(state, code string) (*carwise.GoogleResponse, error) {
	if state != g.state {
		return nil, fmt.Errorf("invalid OAuth state parameter")
	}

	token, err := g.oauth2.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	client := g.oauth2.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &carwise.GoogleResponse{
		Id:        userInfo["id"].(string),
		Email:     userInfo["email"].(string),
		FirstName: userInfo["given_name"].(string),
		LastName:  userInfo["family_name"].(string),
	}, nil
}
