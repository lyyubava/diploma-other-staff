package sendnotification

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type MessagePublishedData struct {
	Message PubSubMessage
}

type EventInfo struct {
	EventName      string
	EventTime      time.Time
	EventUser      string
	EventUserEmail string
}
type Event struct {
	EventDetails string
	EventInfo    EventInfo
}

func SendEmail(msgBody string, mailTo string, subject string) {
	from := os.Getenv("MAIL_FROM")
	password := os.Getenv("MAIL_PASSWORD")

	msg := fmt.Sprintf("From: %s\n To: %s\nSubject: %s\n\n %s", from, mailTo, subject, msgBody)

	err := smtp.SendMail(os.Getenv("SMTP_URI"),
		smtp.PlainAuth("", from, password, os.Getenv("SMTP_HOST")),
		from, []string{mailTo}, []byte(msg))

	if err != nil {
		log.Printf("error while sending email: %s", err)
		return
	}

	log.Printf("email was successfully sent to %s", mailTo)
}

func SendNotification(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %v", err)
	}

	event := Event{}
	data := msg.Message.Data
	json.Unmarshal(data, &event)
	fmt.Println(data)
	subject := event.EventInfo.EventName
	mailTo := event.EventInfo.EventUserEmail
	msgBody := fmt.Sprintf("%s at %s", event.EventDetails, event.EventInfo.EventTime.Format(time.Kitchen))
	SendEmail(msgBody, mailTo, subject)
	return nil
}

func init() {
	functions.CloudEvent("SendNotification", SendNotification)
}
