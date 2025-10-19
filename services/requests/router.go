package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Routes(rg *gin.RouterGroup, rdb *redis.Client) {
	rh := NewRequestsHandler(rdb)
	req := rg.Group("/requests")
	{
		req.POST("/", rh.PostRequest)
		req.PATCH("/", rh.PatchStatus)
		req.GET("/", rh.GetRequestsByUser)
	}
}
