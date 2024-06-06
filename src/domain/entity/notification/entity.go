package notification

import (
	"DX/src/domain/entity"
	"fmt"
	"time"
)

type Notification struct {
	Id               entity.DatabaseId `json:"id" bson:"_id"`
	UserId           string            `json:"user_id" bson:"user_id"`
	ItemId           string            `json:"item_id,omitempty" bson:"item_id"`
	Type             string            `json:"type" bson:"type"`
	NotificationType Type              `json:"-" bson:"notification_type"`
	Title            string            `json:"title" bson:"title"`
	Message          string            `json:"message" bson:"message"`
	Link             string            `json:"link" bson:"link"`
	CreatedAt        time.Time         `json:"created_at" bson:"created_at"`
}

type Type int

const (
	SMS Type = iota
	Email
	InApp
)

func NewBidNotification(userId, errandId string) Notification {
	cTime := time.Now()
	return Notification{
		Id:               entity.NewDatabaseId(),
		UserId:           userId,
		Type:             InApp.Id(),
		NotificationType: InApp,
		Title:            "New errand bid",
		Message:          fmt.Sprintf("You have received a new bid for your errand."),
		CreatedAt:        cTime,
		ItemId:           errandId,
	}
}

func NewErrandStartedNotification(userId, errandId string) Notification {
	cTime := time.Now()
	return Notification{
		Id:               entity.NewDatabaseId(),
		UserId:           userId,
		Type:             InApp.Id(),
		NotificationType: InApp,
		Title:            "Errand started",
		Message:          fmt.Sprintf("The runner has accepted your bid contract and the errand's started."),
		CreatedAt:        cTime,
		ItemId:           errandId,
	}
}

func NewBidProposalRejectedNotification(userId, errandId string) Notification {
	cTime := time.Now()
	return Notification{
		Id:               entity.NewDatabaseId(),
		UserId:           userId,
		Type:             InApp.Id(),
		NotificationType: InApp,
		Title:            "Bid rejected",
		Message:          fmt.Sprintf("The runner has rejected your bid contract."),
		CreatedAt:        cTime,
		ItemId:           errandId,
	}
}

func NewBidAcceptedNotification(userId, bidId string) Notification {
	cTime := time.Now()
	return Notification{
		Id:               entity.NewDatabaseId(),
		UserId:           userId,
		Type:             InApp.Id(),
		NotificationType: InApp,
		Title:            "Bid Accepted",
		Message:          fmt.Sprintf("Congratulations!!! Your bid has been accepted."),
		CreatedAt:        cTime,
		ItemId:           bidId,
	}
}

func NewHaggleNotification(userId, bidId string) Notification {
	cTime := time.Now()
	return Notification{
		Id:               entity.NewDatabaseId(),
		UserId:           userId,
		Type:             InApp.Id(),
		NotificationType: InApp,
		Title:            "New bid haggle",
		Message:          fmt.Sprintf("You have received a new haggle for your bid."),
		CreatedAt:        cTime,
		ItemId:           bidId,
	}
}

func NewErrandUpdateRequestNotification(userId, errandId, title, message string) Notification {
	cTime := time.Now()
	return Notification{
		Id:               entity.NewDatabaseId(),
		UserId:           userId,
		Type:             InApp.Id(),
		NotificationType: InApp,
		Title:            title,
		Message:          message,
		CreatedAt:        cTime,
		ItemId:           errandId,
	}
}

func NewSenderErrandCompletedNotification(userId, errandId string) Notification {
	cTime := time.Now()
	return Notification{
		Id:               entity.NewDatabaseId(),
		UserId:           userId,
		Type:             InApp.Id(),
		NotificationType: InApp,
		Title:            "Errand completed",
		Message:          "Congratulations! Your errand has been marked completed by the sender and your account has been credited",
		CreatedAt:        cTime,
		ItemId:           errandId,
	}
}

func NewRunnerErrandCompletedNotification(userId, errandId string) Notification {
	cTime := time.Now()
	return Notification{
		Id:               entity.NewDatabaseId(),
		UserId:           userId,
		Type:             InApp.Id(),
		NotificationType: InApp,
		Title:            "Errand completed",
		Message:          "The runner for your errand has marked it as completed. Kindly review to accept. ",
		CreatedAt:        cTime,
		ItemId:           errandId,
	}
}

func (t Type) Id() string {
	if t == SMS {
		return "sms"
	}
	if t == Email {
		return "email"
	}
	if t == InApp {
		return "in-app"
	}
	return ""
}

func (t Type) String() string {
	if t == SMS {
		return "SMS"
	}
	if t == Email {
		return "Email"
	}
	if t == InApp {
		return "In-App"
	}
	return ""
}
