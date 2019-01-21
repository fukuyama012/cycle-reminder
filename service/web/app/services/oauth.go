package services

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"os"
	v2 "google.golang.org/api/oauth2/v2"
)

const (
	authorizeEndpoint = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenEndpoint     = "https://www.googleapis.com/oauth2/v4/token"
)

// GetAuthCodeURLWithSessionKey 確認キーを渡しつつAuth接続先URL取得
func GetAuthCodeURLWithSessionKey(key string) string {
	oauth := getConnectConfig()
	return oauth.AuthCodeURL(key)
}

// GetOauthInfo OAUTHにより情報を取得
func GetOauthInfo(code string) (*v2.Tokeninfo, error) {
	oauth := getConnectConfig()
	myContext := context.Background()

	t, err := oauth.Exchange(myContext, code)
	if err != nil {
		return nil, err
	}

	if t.Valid() == false {
		return nil, errors.New("invaild token")
	}

	s, err := v2.New(oauth.Client(myContext, t))
	if err != nil {
		return nil, err
	}

	info, err := s.Tokeninfo().AccessToken(t.AccessToken).Context(myContext).Do()
	if err != nil {
		return nil, err
	}
	return info, nil
}

// getConnectConfig 接続を取得する
func getConnectConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  authorizeEndpoint,
			TokenURL: tokenEndpoint,
		},
		Scopes:      []string{"openid", "email"},
		RedirectURL: os.Getenv("CALL_BACK_URL") + "/auth/callback",
	}
}


