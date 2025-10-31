package auth

import (
	"database/sql"
	"fmt"
	"linkedout/services/auth/utils/JWT"
	"linkedout/services/auth/utils/oAuth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authModel AuthModel
}

func newAuthHandler(db *sql.DB) AuthHandler {
	return AuthHandler{authModel: AuthModel{DB: db}}
}

func (h *AuthHandler) devLogin(c *gin.Context) {

	var code oAuthPayload
	if err := c.ShouldBindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user_id, err := h.authModel.userExists(code.Code)

	var new = false
	if err != nil {
		user_id, err = h.authModel.creatUser(code.Code, "dev_user-"+code.Code)

		new = true
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
	}

	if err := h.authModel.setDeviceCode(user_id, code.DeviceCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	tokens, err := jwt.CreatTokenPair(user_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate tokens"})
		return
	}

	userInfo := UserInfo{Id: user_id, Name: "dev_user-" + code.Code}
	res := LoginRes{tokens.Access, new, userInfo}

	setCookie(c, tokens.Refresh)
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) loginLinkedin(c *gin.Context) {
	println("stargin loginLinkedin")
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "missing code"})
		return
	}
	println("got code")
	userInfo, err := oauth.ExchangeCode(code, oauth.LINKEDIN)

	println("fot info")
	println(userInfo.Name)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user_id, err := h.authModel.userExists(userInfo.Id)
	var new = false
	if err != nil {
		user_id, err = h.authModel.creatUser(userInfo.Id, userInfo.Name)
		new = true

		if err != nil {
			println("failed insertion")
			println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
	}

	println("generating token")
	token, err := jwt.NewAuthToken(user_id)

	println("got token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	redirect := fmt.Sprintf(
		"auth.example.linkedout://auth/callback?token=%s&firstTime=%t",
		token,
		new,
	)

	println("redirectin")
	c.Redirect(http.StatusTemporaryRedirect, redirect)
}

func (h *AuthHandler) linkinCallback(c *gin.Context) {
	println("in call back")
	var payload oAuthPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	println("verifying token")
	tok, err := jwt.Verify(payload.Code, jwt.AuthCode)
	if err != nil {
		println(err.Error())
		c.JSON(http.StatusForbidden, gin.H{"message": "invalid Auth"})
		return
	}

	println("token verifyed")
	tokens, err := jwt.CreatTokenPair(tok.Subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate tokens"})
		return
	}

	if err := h.authModel.setDeviceCode(tok.Subject, payload.DeviceCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	name, err := h.authModel.getUserName(tok.Subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get name"})
		return
	}

	println("created new token pair")

	userInfoRes := UserInfo{Name: name}

	res := LoginRes{tokens.Access, false, userInfoRes}

	println(res.UserInfo.Name)
	setCookie(c, tokens.Refresh)
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) loginGoogle(c *gin.Context) {

	var code oAuthPayload
	if err := c.ShouldBindJSON(&code); err != nil {
		println("failed parse")
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	userInfo, err := oauth.ExchangeCode(code.Code, oauth.GOOGLE)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	h.handleUserLogin(UserInfo{Id: userInfo.Id, Name: userInfo.Name}, code.DeviceCode, c)
}

func (h *AuthHandler) handleUserLogin(userInfo UserInfo, deviceCode string, c *gin.Context) {
	user_id, err := h.authModel.userExists(userInfo.Id)

	var new = false
	if err != nil {
		user_id, err = h.authModel.creatUser(userInfo.Id, userInfo.Name)
		new = true

		if err != nil {
			println("failed insertion")
			println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
	}

	tokens, err := jwt.CreatTokenPair(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate tokens"})
		return
	}

	if err := h.authModel.setDeviceCode(user_id, deviceCode); err != nil {
		println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	userInfoRes := UserInfo{Name: userInfo.Name}

	res := LoginRes{tokens.Access, new, userInfoRes}

	setCookie(c, tokens.Refresh)
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) logOutUser(c *gin.Context) {
	removeCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "logedOut"})
}

func (h *AuthHandler) accessToken(c *gin.Context) {

	reqToken, err := c.Cookie("linkedOut-refresh")
	if err != nil {
		println(err.Error())
		c.JSON(http.StatusForbidden, gin.H{"message": "invalid Refresh"})
		return
	}

	token, err := jwt.Verify(reqToken, jwt.Refresh)
	if err != nil {
		removeCookie(c)
		c.JSON(http.StatusForbidden, gin.H{"message": "invalid Refresh"})
		return
	}

	id := token.Subject

	tokenRotation(c, id)
}

func tokenRotation(c *gin.Context, id string) {
	tokens, err := jwt.CreatTokenPair(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate tokens"})
		return
	}

	setCookie(c, tokens.Refresh)
	c.JSON(http.StatusOK, gin.H{"token": tokens.Access})
}

func removeCookie(c *gin.Context) {

	host := c.Request.Host

	if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		host = parts[0]
	}
	c.SetCookie("linkedOut-refresh", "", -1, "/auth/", host, false, true)
}

func setCookie(c *gin.Context, token string) {

	host := c.Request.Host

	if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		host = parts[0]
	}

	c.SetCookie("linkedOut-refresh", token, 3600, "/auth/", host, false, true)

}
