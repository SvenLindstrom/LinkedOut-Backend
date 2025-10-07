package auth

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"linkedout/services/auth/utils/JWT"
	"linkedout/services/auth/utils/oAuth"
	"net/http"
)

type AuthHandler struct {
	authModel AuthModel
}

func newAuthHandler(db *sql.DB) AuthHandler {
	return AuthHandler{authModel: AuthModel{DB: db}}
}

type oAuthPayload struct {
	Code string `json:"code" binding:"required"`
}

type LoginRes struct {
	Token     string `json:"token"`
	FirstTime bool   `json:"firstTime"`
}

func (h *AuthHandler) devLogin(c *gin.Context) {

	var code oAuthPayload
	if err := c.ShouldBindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user_id, err := h.authModel.userExists(code.Code)

	var loginRes = LoginRes{FirstTime: false}
	if err != nil {
		user_id, err = h.authModel.creatUser(code.Code, "dev_user-"+code.Code)

		loginRes = LoginRes{FirstTime: true}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
	}

	tokens, err := jwt.CreatTokenPair(user_id)

	loginRes.Token = tokens.Access
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate tokens"})
		return
	}

	c.SetCookie("linkedOut-refresh", tokens.Refresh, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, loginRes)
}

func (h *AuthHandler) login(c *gin.Context) {

	var code oAuthPayload
	if err := c.ShouldBindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	googleInfo, err := oauth.ExchangeCode(code.Code)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user_id, err := h.authModel.userExists(googleInfo.ID)

	var loginRes = LoginRes{FirstTime: false}
	if err != nil {
		user_id, err = h.authModel.creatUser(googleInfo.ID, googleInfo.Name)
		loginRes = LoginRes{FirstTime: true}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
	}

	tokens, err := jwt.CreatTokenPair(user_id)

	loginRes.Token = tokens.Access
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate tokens"})
		return
	}

	c.SetCookie("linkedOut-refresh", tokens.Refresh, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, loginRes)
}

func (h *AuthHandler) accessToken(c *gin.Context) {

	reqToken, err := c.Cookie("linkedOut-refresh")

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "invalid Refresh"})
		return
	}

	token, err := jwt.Verify(reqToken, jwt.Refresh)

	if err != nil {
		c.SetCookie("linkedOut-refresh", "", -1, "/", "localhost", false, true)
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

	c.SetCookie("linkedOut-refresh", tokens.Refresh, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"token": tokens.Access})
}
