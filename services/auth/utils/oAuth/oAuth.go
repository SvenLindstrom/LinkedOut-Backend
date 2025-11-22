package oauth

import (
	"golang.org/x/net/context"
)

func ExchangeCode(code string, provider string) (UserInfo, error) {
	authConfig := getConfig[provider]
	tok, err := authConfig.Config.Exchange(context.Background(), code)
	if err != nil {
		println(err.Error())
		return UserInfo{}, err
	}

	client := authConfig.Config.Client(context.Background(), tok)

	res, err := client.Get(authConfig.Endpoint)
	if err != nil {
		return UserInfo{}, err
	}
	defer res.Body.Close()

	info := authConfig.Res.ToUserInfo(res.Body)

	return info, nil
}
