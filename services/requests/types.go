package requests

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	ID           string    `json:"id"`
	From         string    `json:"from"`
	SenderName   string    `json:"sender"`
	To           string    `json:"to"`
	ReceiverName string    `json:"receiver"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	Timestamp    time.Time `json:"timestamp"`
}

func NewRequest(from, sender, to, receiver, message string) *Request {
	return &Request{
		ID:           uuid.NewString(),
		From:         from,
		SenderName:   sender,
		To:           to,
		ReceiverName: receiver,
		Status:       "PENDING",
		Message:      message,
		Timestamp:    time.Now(),
	}
}

func ToMap(request Request) map[string]string {
	fields := map[string]string{
		"id":        request.ID,
		"origin":    request.From,
		"sender":    request.SenderName,
		"to":        request.To,
		"receiver":  request.ReceiverName,
		"status":    request.Status,
		"message":   request.Message,
		"timestamp": request.Timestamp.Format(time.RFC3339),
	}
	return fields
}
