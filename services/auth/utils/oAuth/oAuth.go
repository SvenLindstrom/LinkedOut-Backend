package oauth

import (
	"encoding/json"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func ExchangeCode(code string) (GoogleUser, error) {

	conf := &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_ID"),
		ClientSecret: os.Getenv("OAUTH_SECRET"),
		RedirectURL:  os.Getenv("OAUTH_REDIRECT"),
		Scopes:       []string{"openid"},
		Endpoint:     google.Endpoint,
	}

	tok, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return GoogleUser{}, err
	}

	client := conf.Client(context.Background(), tok)

	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return GoogleUser{}, err
	}
	defer res.Body.Close()

	var userinfo GoogleUser
	json.NewDecoder(res.Body).Decode(&userinfo)

	return userinfo, nil
}
