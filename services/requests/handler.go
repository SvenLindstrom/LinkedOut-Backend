package requests

import (
	"database/sql"
	fcm "linkedout/services/FCM"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RequestsHandler struct {
	rm  *RequestsModel
	fcm *fcm.FcmClient
}

type CreatePayload struct {
	SenderName   string `json:"sender"   binding:"required"`
	To           string `json:"to"       binding:"required"`
	ReceiverName string `json:"receiver" binding:"required"`
	Message      string `json:"message"  binding:"required"`
}

type UpdatePayload struct {
	Status    string `json:"status"    binding:"required"`
	RequestID string `json:"requestID" binding:"required"`
}

func NewRequestsHandler(rdb *redis.Client, pg *sql.DB) *RequestsHandler {
	return &RequestsHandler{rm: &RequestsModel{rdb: rdb}, fcm: fcm.GetClient(pg)}
}

func (rh *RequestsHandler) PostRequest(c *gin.Context) {
	ctx := c.Request.Context()
	from := c.GetString("x-user-id")

	var payload CreatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req, err := rh.rm.CreateRequest(
		ctx,
		from,
		payload.SenderName,
		payload.To,
		payload.ReceiverName,
		payload.Message,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = rh.fcm.Send(payload.SenderName, payload.Message, payload.To, true, ToMap(req))

	if err != nil {
		println(err.Error())
	}

	c.JSON(http.StatusCreated, req)
}

func (rh *RequestsHandler) PatchStatus(c *gin.Context) {
	ctx := c.Request.Context()
	from := c.GetString("x-user-id")

	var payload UpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !rh.rm.CheckRequestMember(ctx, from, payload.RequestID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request status updated not allowed"})
		return
	}

	err := rh.rm.UpdateRequestStatus(ctx, payload.RequestID, payload.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedRequest, err := rh.rm.FindRequest(ctx, payload.RequestID)

	if err != nil {
		println("failed notification could not find request")
	}

	notify := payload.Status == "ACCEPTED"
	err = rh.fcm.Send(
		updatedRequest.ReceiverName,
		"Requst Accepted",
		updatedRequest.From,
		notify,
		ToMap(*updatedRequest),
	)

	if err != nil {
		println(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request status updated successfully"})
}

func (rh *RequestsHandler) GetRequestsByUser(c *gin.Context) {
	ctx := c.Request.Context()
	from := c.GetString("x-user-id")

	requests, err := rh.rm.FindRequestsByUser(ctx, from)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}
