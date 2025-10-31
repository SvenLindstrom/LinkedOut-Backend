package oauth

import (
	"encoding/json"
	"io"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/linkedin"
)

const (
	GOOGLE   string = "google"
	LINKEDIN string = "linkedin"
)

type UserInfo struct {
	Id   string
	Name string
}

type OAuthRes interface {
	ToUserInfo(body io.Reader) UserInfo
}

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

func (u GoogleUserInfo) ToUserInfo(body io.Reader) UserInfo {
	json.NewDecoder(body).Decode(&u)
	return UserInfo{
		Id:   "google_" + u.ID,
		Name: u.Name,
	}
}

type OAuthConfig struct {
	Config   *oauth2.Config
	Endpoint string
	Res      OAuthRes
}

func getGoogleConf() OAuthConfig {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_ID"),
		ClientSecret: os.Getenv("OAUTH_SECRET"),
		RedirectURL:  os.Getenv("OAUTH_REDIRECT"),
		Scopes:       []string{"openid"},
		Endpoint:     google.Endpoint,
	}
	endpoint := "https://www.googleapis.com/oauth2/v2/userinfo"
	var googleUser GoogleUserInfo
	return OAuthConfig{Config: conf, Endpoint: endpoint, Res: googleUser}
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

func getConf(provider string) OAuthConfig {
	if provider == GOOGLE {
		return getGoogleConf()
	} else {
		return getLinkedinConf()
	}
}

func ExchangeCode(code string, provider string) (UserInfo, error) {
	println("starting exchange")
	authConfig := getConf(provider)
	tok, err := authConfig.Config.Exchange(context.Background(), code)
	if err != nil {
		println(err.Error())
		return UserInfo{}, err
	}

	println("got token")
	client := authConfig.Config.Client(context.Background(), tok)

	res, err := client.Get(authConfig.Endpoint)
	if err != nil {
		return UserInfo{}, err
	}
	println("got data")
	defer res.Body.Close()

	println(res.Body)

	// body, _ := io.ReadAll(res.Body)
	// fmt.Println("POST /login ->", string(body))

	info := authConfig.Res.ToUserInfo(res.Body)

	println(info.Name)
	return info, nil
}
