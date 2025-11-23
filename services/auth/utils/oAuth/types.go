package oauth

import (
	"io"

	"golang.org/x/oauth2"
)

var getConfig = map[string]func() OAuthConfig{
	"google":   getGoogleConf,
	"linkedin": getLinkedinConf,
}

type OAuthConfig struct {
	Config   *oauth2.Config
	Endpoint string
	Res      OAuthRes
}

type OAuthRes interface {
	ToUserInfo(body io.Reader) UserInfo
}

type UserInfo struct {
	Id   string
	Name string
}
