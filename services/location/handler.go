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
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err := h.locationModel.UpdateLocation(id, location)

	if err != nil {
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
	var location Location

	id := c.GetString("x-user-id")

	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	println(location.Lat)
	println(location.Lon)

	data, err := h.locationModel.getProximity(id, location)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
