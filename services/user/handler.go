package user

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userModel UserModel
}

type InfoPayload struct {
	Profession string     `json:"profession" binding:"required"`
	Bio        string     `json:"bio"        binding:"required"`
	Interests  []Interest `json:"interests"  binding:"required"`
}

func newUserHandler(db *sql.DB) UserHandler {
	return UserHandler{userModel: UserModel{DB: db}}
}

func (h *UserHandler) GetInfo(c *gin.Context) {
	user_id := c.GetString("x-user-id")

	info, err := h.userModel.getInfo(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (uh *UserHandler) PutUserInfo(c *gin.Context) {
	id := c.GetString("x-user-id")

	var payload InfoPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uh.userModel.SaveInfo(id, payload.Profession, payload.Bio, payload.Interests)
	if err != nil {
		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "info updated successfully"})
}

func (uh *UserHandler) GetInterests(c *gin.Context) {
	interests, err := uh.userModel.FindAllInterests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, interests)
}
