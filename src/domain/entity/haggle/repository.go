package haggle

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type repository struct {
	Collection *mongo.Collection
}

func NewRepository(Collection *mongo.Collection) Repository {
	return &repository{Collection: Collection}
}

func (r *repository) GetHagglesForBid(bidId string) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) CreateHaggle(haggle *Haggle) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.Collection.InsertOne(ctx, haggle)
	if err != nil {
		return err
	}

	return nil
}
