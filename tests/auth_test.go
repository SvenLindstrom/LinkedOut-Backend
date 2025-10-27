package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (u *User) getNewAcces(client *http.Client) {
	type Token struct {
		Token string
	}
	req, err := http.NewRequest(http.MethodGet,
		"http://localhost:3113/auth/access_token",
		bytes.NewBufferString(""),
	)

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "linkedOut-refresh", Value: u.Refresh})

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()

	refresh := cookies[0].Value

	acces_token := &Token{}
	json.NewDecoder(resp.Body).Decode(acces_token)

	if acces_token.Token == u.Token {
		fmt.Println("Same")
	}

	u.Refresh = refresh
	u.Token = acces_token.Token

}
func loginOauth(code string, client *http.Client) {

	data := map[string]string{"code": code}
	jsonData, _ := json.Marshal(data)

	req, err := http.NewRequest(http.MethodPost,
		"http://localhost:3113/auth/login",
		bytes.NewBuffer(jsonData),
	)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("POST /login ->", string(body))
}
