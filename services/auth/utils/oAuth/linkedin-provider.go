package oauth

import (
	"encoding/json"
	"io"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

type LinkedInUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (u LinkedInUserInfo) ToUserInfo(body io.Reader) UserInfo {
	json.NewDecoder(body).Decode(&u)
	return UserInfo{
		Id:   "linkedin_" + u.Sub,
		Name: u.Name,
	}
}

func getLinkedinConf() OAuthConfig {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_LINKEDIN_ID"),
		ClientSecret: os.Getenv("OAUTH_LINKEDIN_SECRET"),
		RedirectURL:  "https://svenlindstroem.dev/auth/linkedin",
		Scopes:       []string{"openid"},
		Endpoint:     linkedin.Endpoint,
	}
	endpoint := "https://api.linkedin.com/v2/userinfo"
	var linkedinUser LinkedInUserInfo

	return OAuthConfig{Config: conf, Endpoint: endpoint, Res: linkedinUser}
}
