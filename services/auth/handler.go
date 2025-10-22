package auth

import (
	"database/sql"
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

	userInfo := UserInfo{Id: user_id, Name: "testname"}
	res := LoginRes{tokens.Access, new, userInfo}

	setCookie(c, tokens.Refresh)
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) login(c *gin.Context) {

	var code oAuthPayload
	if err := c.ShouldBindJSON(&code); err != nil {
		println("failed parse")
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	googleInfo, err := oauth.ExchangeCode(code.Code)
	if err != nil {
		println("failed exchange")
		println(code.Code)
		println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user_id, err := h.authModel.userExists(googleInfo.ID)

	var new = false

	if err != nil {
		user_id, err = h.authModel.creatUser(googleInfo.ID, googleInfo.Name)
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

	userInfo := UserInfo{Name: googleInfo.Name}

	res := LoginRes{tokens.Access, new, userInfo}

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
