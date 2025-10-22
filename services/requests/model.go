package requests

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RequestsModel struct {
	rdb *redis.Client
}

func NewRequestsModel(rdb *redis.Client) *RequestsModel {
	return &RequestsModel{rdb: rdb}
}

func (rm *RequestsModel) CreateRequest(
	ctx context.Context,
	from, sender, to, receiver, message string,
) (Request, error) {
	newRequest := NewRequest(from, sender, to, receiver, message)
	fields := map[string]any{
		"id":        newRequest.ID,
		"from":      newRequest.From,
		"sender":    newRequest.SenderName,
		"to":        newRequest.To,
		"receiver":  newRequest.ReceiverName,
		"status":    newRequest.Status,
		"message":   newRequest.Message,
		"timestamp": newRequest.Timestamp.Format(time.RFC3339),
	}

	tx := rm.rdb.TxPipeline()
	tx.HSet(ctx, "requests:"+newRequest.ID, fields)
	tx.SAdd(ctx, "users:"+from, newRequest.ID)
	tx.SAdd(ctx, "users:"+to, newRequest.ID)
	ok, err := tx.Expire(ctx, "requests:"+newRequest.ID, 30*time.Minute).Result()
	if err != nil {
		println("error when adding expiration")
	}
	if ok {
		println("added 30 min expiration successfully")
	}
	_, err = tx.Exec(ctx)

	if err != nil {
		return *newRequest, err
	}

	return *newRequest, nil
}

func (rm *RequestsModel) CheckRequestMember(ctx context.Context, userID, requestID string) bool {
	isMember, err := rm.rdb.SIsMember(ctx, "users:"+userID, requestID).Result()
	if err != nil {
		return false
	}

	if isMember {
		return true
	}

	return false
}

func (rm *RequestsModel) UpdateRequestStatus(ctx context.Context, requestID, status string) error {
	tx := rm.rdb.TxPipeline()
	tx.HSet(ctx, "requests:"+requestID, "status", status)
	tx.Expire(ctx, "requests:"+requestID, 30*time.Minute)
	_, err := tx.Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (rm *RequestsModel) FindRequestsIds(ctx context.Context, userID string) ([]string, error) {
	IDs, err := rm.rdb.SMembers(ctx, "users:"+userID).Result()
	if err != nil {
		return nil, err
	}

	return IDs, nil
}

func (rm *RequestsModel) FindRequestsByUser(
	ctx context.Context,
	userID string,
) ([]*Request, error) {
	ids, err := rm.FindRequestsIds(ctx, userID)
	if err != nil {
		return nil, err
	}

	requests := make([]*Request, 0)
	for _, id := range ids {
		req, err := rm.FindRequest(ctx, id)
		if err != nil {
			if err.Error() == "request ID expired" {
				rm.rdb.SRem(ctx, "users:"+userID, id)
				continue
			}
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (rm *RequestsModel) FindRequest(ctx context.Context, requestID string) (*Request, error) {
	val, err := rm.rdb.HGetAll(ctx, "requests:"+requestID).Result()
	if err != nil {
		return nil, err
	}

	if len(val) == 0 {
		return nil, fmt.Errorf("request ID expired")
	}

	time, err := time.Parse(time.RFC3339, val["timestamp"])
	if err != nil {
		return nil, fmt.Errorf("could not parse timestamp")
	}

	req := &Request{
		ID:           val["id"],
		From:         val["from"],
		SenderName:   val["sender"],
		To:           val["to"],
		ReceiverName: val["receiver"],
		Status:       val["status"],
		Message:      val["message"],
		Timestamp:    time,
	}

	return req, nil
}
