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
		ID:           "requests:" + uuid.NewString(),
		From:         from,
		SenderName:   sender,
		To:           to,
		ReceiverName: receiver,
		Status:       "pending",
		Message:      message,
		Timestamp:    time.Now(),
	}
}
