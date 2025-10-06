package location

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)




func Routes(rg *gin.RouterGroup, db *sql.DB){
	h := newLocationHandler(db)
	loc := rg.Group("/location")
	{
		loc.PATCH("/", h.UpdateLocation)
		loc.PATCH("/status/:status", h.UpdateStatus)
		loc.POST("/", h.GetProxUsers)
	}
}
