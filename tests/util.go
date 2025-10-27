package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"linkedout/services/auth"
	"linkedout/services/location"
	"linkedout/services/user"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Token   string
	Refresh string
	Client  *http.Client
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type UserInfo struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Profession string     `json:"profession"`
	Bio        string     `json:"bio"`
	Interests  []Interest `json:"interests"`
}

type Interest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type LoginResponse struct {
	Token    string   `json:"token"`
	New      bool     `json:"new"`
	UserInfo UserInfo `json:"userInfo"`
}

type Proximity struct {
	Location Location `json:"location" binding:"required"`
	Distance int32    `json:"distance" binding:"required"`
}

type UserProx struct {
	Id         string     `json:"id"         binding:"required"`
	Name       string     `json:"name"       binding:"required"`
	Bio        string     `json:"bio"`
	Distance   string     `json:"distance"   binding:"required"`
	Profession string     `json:"profession"`
	Tags       []Interest `json:"tags"`
}

type CreatePayload struct {
	SenderName   string `json:"sender" binding:"required"`
	To           string `json:"to" binding:"required"`
	ReceiverName string `json:"receiver" binding:"required"`
	Message      string `json:"message" binding:"required"`
}

var Home = Location{56.050248, 14.148894}

var Uni = Location{56.048533, 14.145616}

var PhoneDafault = Location{37.4219983, -122.084}

func PgTestInit() *sql.DB {
	pwd := os.Getenv("POSTGRES_PASSWORD")
	name := os.Getenv("POSTGRES_DB")
	addr := os.Getenv("POSTGRES_ADDR")

	url := fmt.Sprintf("postgres://postgres:%s@%s/%s?sslmode=disable", pwd, addr, name)
	db, err := sql.Open("pgx", url)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func ServerInit() *gin.Engine {

	if err := godotenv.Load("../.test.env"); err != nil {
		log.Fatal("failed to load .test.env:", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	pg := PgTestInit()

	authGroup := r.Group("/auth")
	auth.Routes(authGroup, pg)

	api := r.Group("/api")
	api.Use(auth.TokenMiddleware())
	location.Routes(api, pg)
	user.Routes(api, pg)

	return r
}

func Login(code string, t *testing.T) User {
	client := &http.Client{}

	payload := map[string]string{
		"code":       code,
		"deviceCode": "test-device-001",
	}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:3113/auth/dev/login", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "expected 200 OK on login")

	cookies := resp.Cookies()
	assert.NotEmpty(t, cookies, "expected refresh cookie after login")

	refresh := cookies[0].Value

	var loginRes LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginRes)
	assert.NoError(t, err)
	assert.NotEmpty(t, loginRes.Token, "expected non-empty access token")

	return User{Token: loginRes.Token, Refresh: refresh, Client: client}
}

func (u *User) UpdateLocation(loc Location) {
	u.authRequest(http.MethodPatch, "http://localhost:3113/api/location", loc)
}

func (u *User) UpdateState(state string) {
	u.authRequest(http.MethodPatch, "http://localhost:3113/api/location/status/"+state, nil)
}

func (u *User) authRequest(method string, endpoint string, data any) (*http.Response, error) {
	jsonData, _ := json.Marshal(data)

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+u.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.Client.Do(req)

	return resp, err
}
