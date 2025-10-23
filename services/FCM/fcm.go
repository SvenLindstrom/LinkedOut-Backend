package fcm

import (
	"context"
	"database/sql"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

type FcmClient struct {
	client *messaging.Client
	DB     *sql.DB
}

func GetClient(pg *sql.DB) *FcmClient {

	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return &FcmClient{client: client, DB: pg}
}

func (f *FcmClient) Send(
	title string,
	message string,
	reciver string,
	notify bool,
	request map[string]string,
) error {

	token, err := f.getDeviceCode(reciver)

	if err != nil {
		return err
	}

	var msg messaging.Message

	if notify {
		conf := toNotification(title, message)
		msg = messaging.Message{Android: &conf, Token: token, Data: request}
	} else {
		msg = messaging.Message{Token: token, Data: request}
	}

	ctx := context.Background()
	mss, err := f.client.Send(ctx, &msg)
	if err != nil {
		println(err.Error())
	} else {
		println(mss)
	}

	return nil
}

func toNotification(name string, message string) messaging.AndroidConfig {
	notification := messaging.AndroidNotification{
		Title:       name,
		Body:        message,
		ClickAction: "OPEN_REQUESTS",
	}
	return messaging.AndroidConfig{Notification: &notification}
}

func (f *FcmClient) getDeviceCode(id string) (string, error) {

	var deviceCode string
	err := f.DB.QueryRow("SELECT deviceCode FROM users WHERE id=$1", id).
		Scan(&deviceCode)

	return deviceCode, err
}
