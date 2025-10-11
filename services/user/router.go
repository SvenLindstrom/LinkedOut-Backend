package user

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Routes(rg *gin.RouterGroup, db *sql.DB) {
	h := newUserHandler(db)
	loc := rg.Group("/user")

	loc.GET("info", h.getInfo)
	loc.PATCH("bio", h.updateBio)
}
