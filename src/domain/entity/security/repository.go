package security

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type repository struct {
	Collection *mongo.Collection
}

func NewDatabaseRepository(collection *mongo.Collection) Repository {
	return &repository{Collection: collection}
}

func (r *repository) Get(value string) (*Security, error) {
	sec := &Security{}
	mongoContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"user_id": value},
			{"phone_number": fmt.Sprintf("+%s", value)}, // Adding the "+" because the query coming from the API doesn't contain the +.
		},
	}
	if err := r.Collection.FindOne(mongoContext, filter).Decode(&sec); err != nil {
		return nil, err
	}
	return sec, nil
}

func (r *repository) Create(sec *Security) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.Collection.InsertOne(ctx, sec)
	return err
}
