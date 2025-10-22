package location

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	locationModel LocationModel
}

func newLocationHandler(db *sql.DB) LocationHandler {
	return LocationHandler{locationModel: LocationModel{DB: db}}
}

func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	var location Location

	id := c.GetString("x-user-id")
	if err := c.ShouldBindJSON(&location); err != nil {
		println("cant bind")
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err := h.locationModel.UpdateLocation(id, location)

	if err != nil {
		println("cant update")
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, location)
}

func (h *LocationHandler) UpdateStatus(c *gin.Context) {

	id := c.GetString("x-user-id")
	status, err := strconv.ParseBool(c.Param("status"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = h.locationModel.UpdateStatus(id, status)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

func (h *LocationHandler) GetProxUsers(c *gin.Context) {
	var prox Proximity

	id := c.GetString("x-user-id")

	if err := c.ShouldBindJSON(&prox); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	data, err := h.locationModel.getProximity(id, prox)

	if err != nil {
		println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
