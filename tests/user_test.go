package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInterests(t *testing.T) {
	router := ServerInit()
	go func() {
		_ = router.Run(":3113")
	}()

	user := Login("user_001", t)

	resp, err := user.authRequest(http.MethodGet, "http://localhost:3113/api/user/interests", nil)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var interests []Interest
	err = json.NewDecoder(resp.Body).Decode(&interests)
	assert.NoError(t, err)

	assert.NotEmpty(t, interests, "expected at least one interest")
	for _, i := range interests {
		t.Logf("Interest: %s (%s)", i.Name, i.Id)
	}
}

func TestGetUserInfo(t *testing.T) {
	router := ServerInit()
	go func() {
		_ = router.Run(":3113")
	}()

	user := Login("user_001", t)

	resp, err := user.authRequest(http.MethodGet, "http://localhost:3113/api/user/info", nil)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var info UserInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	assert.NoError(t, err)

	assert.Equal(t, "aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa", info.Id)
}
