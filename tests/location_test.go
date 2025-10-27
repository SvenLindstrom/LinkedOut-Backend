package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProxUsers(t *testing.T) {
	router := ServerInit()
	go func() {
		_ = router.Run(":3113")
	}()

	user := Login("user_001", t)

	var Home = Location{56.050248, 14.148894}
	user2 := Login("user_002", t) // interests in common: UI/UX, Project Management
	user2.UpdateLocation(Home)
	user2.UpdateState("true")
	user3 := Login("user_003", t) // nothing in common
	user3.UpdateLocation(Home)
	user3.UpdateState("true")

	prox := Proximity{Home, 100}

	resp, err := user.authRequest(http.MethodPost, "http://localhost:3113/api/location/", prox)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var nearbyUsers []UserProx
	err = json.NewDecoder(resp.Body).Decode(&nearbyUsers)
	assert.NoError(t, err)

	assert.NotEmpty(t, nearbyUsers, "expected user_002")
	assert.Equal(t, "Bob Test", nearbyUsers[0].Name)

	matchingInterests := []Interest{
		{"55555555-5555-5555-5555-555555555555", "UI/UX Design"},
		{"77777777-7777-7777-7777-777777777777", "Project Management"}}
	assert.Equal(t, matchingInterests, nearbyUsers[0].Tags)
}
