package services

import (
	"golang.org/x/oauth2"
	"os"
)

const (
	authorizeEndpoint = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenEndpoint     = "https://www.googleapis.com/oauth2/v4/token"
)

func GetAuthCodeUrlWithSessionKey(key string) string {
	oauth := getConnectConfig()
	return oauth.AuthCodeURL(key)
}

// GetConnect 接続を取得する
func getConnectConfig() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  authorizeEndpoint,
			TokenURL: tokenEndpoint,
		},
		Scopes:      []string{"openid", "email", "profile"},
		RedirectURL: os.Getenv("CALL_BACK_URL") + "/login/callback",
	}
	return config
}


