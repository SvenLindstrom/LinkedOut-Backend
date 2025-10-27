package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserProx struct {
	Id       string `json:"id"       binding:"required"`
	Name     string `json:"name"     binding:"required"`
	Bio      string `json:"bio"`
	Distance string `json:"distance" binding:"required"`
}

type Proximity struct {
	Location Location `json:"location" binding:"required"`
	Distance int32    `json:"distance" binding:"required"`
}

type CreatePayload struct {
	SenderName   string `json:"sender" binding:"required"`
	To           string `json:"to" binding:"required"`
	ReceiverName string `json:"receiver" binding:"required"`
	Message      string `json:"message" binding:"required"`
}

// type UpdatePayload struct {
// 	Status    string `json:"status" binding:"required"`
// 	RequestID string `json:"requestID" binding:"required"`
// }

// type User struct {
// 	Token   string
// 	Refresh string
// 	Client  *http.Client
// }

// type Location struct {
// 	Lat float64 `json:"lat"`
// 	Lon float64 `json:"lon"`
// }

// type UserInfo struct {
// 	Id         string     `json:"id"`
// 	Name       string     `json:"name"`
// 	Profession string     `json:"profession"`
// 	Bio        string     `json:"bio"`
// 	Interests  []Interest `json:"interests"`
// }

// type Interest struct {
// 	Id   string `json:"id"`
// 	Name string `json:"name"`
// }

// type LoginResponse struct {
// 	Token    string   `json:"token"`
// 	New      bool     `json:"new"`
// 	UserInfo UserInfo `json:"userInfo"`
// }

var Home = Location{56.050248, 14.148894}

var Uni = Location{56.048533, 14.145616}

var PhoneDafault = Location{37.4219983, -122.084}

func main() {

	client := &http.Client{}

	// loginOauth(code, client)

	// Example POST request with JSON
	//
	user1 := loginUser("123", client)
	user1.updateLoc(PhoneDafault, client)
	user1.updateState("true", client)

	user2 := loginUser("456", client)
	user2.updateLoc(PhoneDafault, client)
	user2.updateState("true", client)

	//proxs := user2.getProx(Proximity{PhoneDafault, 50}, client)

	// user2.authReq()

	// user2.updateLoc(PhoneDafault, client)
	// user2.updateState("true", client)

	//
	// user2.updateLoc(Uni, client)
	//
	// user2.updateState(client, "false")
	//
	// user1.getProx(Uni, client)
	//
	// user1.getNewAcces(client)

	p := &CreatePayload{
		SenderName:   "Bianca",
		To:           "12347",
		ReceiverName: "Sven",
		Message:      "want to hang out?",
	}

	// u := &UpdatePayload{
	// 	Status:    "accepted",
	// 	RequestID: "b7b719d8-22d6-4fba-8b16-275458724a92",
	// }

	// create request
	user1.authReq(http.MethodPost, p, "http://localhost:3113/api/requests/", client)
	// update request status
	//user1.authReq(http.MethodPatch, u, "http://localhost:3113/api/requests/", client)
	// get requests by user
	//user1.authReq(http.MethodGet, p, "http://localhost:3113/api/requests/", client)

	//user1.authReq(http.MethodGet, p, "http://localhost:3113/api/user/info", client)

}

func (u *User) authReq(method string, data any, endpoint string, client *http.Client) {

	req := reqWithAuth(
		method,
		data,
		endpoint,
		u,
	)

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println("${method}  ->", string(body))
}

func (u *User) getProx(location Proximity, client *http.Client) *[]UserProx {
	req := reqWithAuth(http.MethodPost,
		location,
		"http://localhost:3113/api/location",
		u)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	proxs := &[]UserProx{}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(proxs)
	return proxs
}

func (u *User) updateState(state string, client *http.Client) {
	req := reqWithAuth(
		http.MethodPatch,
		"",
		"http://localhost:3113/api/location/status/"+state,
		u)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("POST /login ->", string(body))

}

func (u *User) updateLoc(location Location, client *http.Client) {

	req := reqWithAuth(
		http.MethodPatch,
		location,
		"http://localhost:3113/api/location",
		u)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("POST /login ->", string(body))
}

func reqWithAuth(method string, data any, endpoint string, u *User) *http.Request {

	jsonData, _ := json.Marshal(data)

	req, err := http.NewRequest(
		method,
		endpoint,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		panic(err)
	}

	bear := "Bearer " + u.Token

	req.Header.Set("Authorization", bear)
	req.Header.Set("Content-Type", "application/json")
	return req
}

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

func loginUser(code string, client *http.Client) User {

	type UserInfo struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	type LoginRes struct {
		Token    string   `json:"token"`
		New      bool     `json:"new"`
		UserInfo UserInfo `json:"userInfo"`
	}

	data := map[string]string{"code": code}
	jsonData, _ := json.Marshal(data)

	req, err := http.NewRequest(http.MethodPost,
		"http://localhost:3113/auth/dev/login",
		bytes.NewBuffer(jsonData),
	)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()

	refresh := cookies[0].Value

	// body, _ := io.ReadAll(resp.Body)
	// fmt.Println("POST /login ->", string(body))
	loginRes := &LoginRes{}
	json.NewDecoder(resp.Body).Decode(loginRes)

	user := User{loginRes.Token, refresh, client}

	fmt.Println("Refresh: ", user.Refresh)
	fmt.Println("Acces: ", user.Token)

	println(loginRes.New)

	return user
}
