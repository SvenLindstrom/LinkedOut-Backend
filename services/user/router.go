package user

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Routes(rg *gin.RouterGroup, db *sql.DB) {
	h := newUserHandler(db)
	loc := rg.Group("/user")

	loc.GET("/info", h.GetInfo)
	loc.PUT("/info", h.PutUserInfo)
	loc.GET("/interests", h.GetInterests)
}
