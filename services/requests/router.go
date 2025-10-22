package requests

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Routes(rg *gin.RouterGroup, rdb *redis.Client, pg *sql.DB) {
	rh := NewRequestsHandler(rdb, pg)
	req := rg.Group("/requests")
	{
		req.POST("", rh.PostRequest)
		req.PATCH("", rh.PatchStatus)
		req.GET("", rh.GetRequestsByUser)
	}
}
