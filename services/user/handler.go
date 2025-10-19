package user

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userModel UserModel
}

func newUserHandler(db *sql.DB) UserHandler {
	return UserHandler{userModel: UserModel{DB: db}}
}

func (h *UserHandler) getInfo(c *gin.Context) {
	user_id := c.GetString("x-user-id")

	info, err := h.userModel.getInfo(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (h *UserHandler) updateBio(c *gin.Context) {

	type Payload struct {
		Bio string `json:"bio"`
	}

	user_id := c.GetString("x-user-id")

	var bio Payload
	if err := c.ShouldBindJSON(&bio); err != nil {
		println("cant bind")
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err := h.userModel.updateBio(user_id, bio.Bio)

	if err != nil {
		println("cant update")
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bio)

}
