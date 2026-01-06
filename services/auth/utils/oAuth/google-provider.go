package oauth

import (
	"encoding/json"
	"io"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (u GoogleUserInfo) ToUserInfo(body io.Reader) UserInfo {
	json.NewDecoder(body).Decode(&u)
	return UserInfo{
		Id:   "google_" + u.ID,
		Name: u.Name,
	}
}

func getGoogleConf() OAuthConfig {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_GOOGLE_ID"),
		ClientSecret: os.Getenv("OAUTH_GOOGLE_SECRET"),
		RedirectURL:  os.Getenv("OAUTH_GOOGLE_REDIRECT"),
		Scopes:       []string{"openid"},
		Endpoint:     google.Endpoint,
	}
	endpoint := "https://www.googleapis.com/oauth2/v2/userinfo"
	var googleUser GoogleUserInfo

	return OAuthConfig{Config: conf, Endpoint: endpoint, Res: googleUser}
}
