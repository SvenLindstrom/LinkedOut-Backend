package fcm

import (
	"context"
	"database/sql"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

type FcmClient struct {
	client *messaging.Client
	DB     *sql.DB
}

func GetClient(pg *sql.DB) (FcmClient, error) {

	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return FcmClient{}, err
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return FcmClient{}, err
	}

	return FcmClient{client: client, DB: pg}, nil
}

func (f *FcmClient) Send(from string, name string, message string) error {

	token, err := f.getDeviceCode(from)

	if err != nil {
		return err
	}

	conf := toNotification(name, message)
	msg := messaging.Message{Android: &conf, Token: token}

	ctx := context.Background()
	f.client.Send(ctx, &msg)

	return nil
}

func toNotification(name string, message string) messaging.AndroidConfig {
	notification := messaging.AndroidNotification{
		Title: name,
		Body:  message,
	}
	return messaging.AndroidConfig{Notification: &notification}
}

func (f *FcmClient) getDeviceCode(id string) (string, error) {

	var deviceCode string
	err := f.DB.QueryRow("SELECT deviceCode FROM users WHERE id=$1", id).
		Scan(&deviceCode)

	return deviceCode, err
}
