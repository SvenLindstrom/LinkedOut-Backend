package auth

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Routes(rg *gin.RouterGroup, db *sql.DB) {
	h := newAuthHandler(db)

	rg.GET("access_token", h.accessToken)
	rg.POST("login", h.login)
	rg.POST("dev/login", h.devLogin)
}
