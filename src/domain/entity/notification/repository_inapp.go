package notification

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type inAppNotification struct {
	*mongo.Collection
}

func NewInAppNotificationRepository(Collection *mongo.Collection) Repository {
	return &inAppNotification{
		Collection: Collection,
	}
}

func (i *inAppNotification) GetAllNotifications(userId string) (notifications []Notification, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": userId,
	}

	crs, err := i.Collection.Find(ctx, filter)
	if err = crs.All(ctx, &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (i *inAppNotification) SendNotification(notification Notification) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := i.Collection.InsertOne(ctx, notification); err != nil {
		return err
	}

	return nil
}
