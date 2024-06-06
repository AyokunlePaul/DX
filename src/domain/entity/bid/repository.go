package bid

import (
	"DX/src/domain/entity"
	"DX/src/domain/entity/haggle"
	"DX/src/utils/logger"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type repository struct {
	*mongo.Collection
}

func NewRepository(coll *mongo.Collection) Repository {
	return &repository{Collection: coll}
}

func (r *repository) GetBid(id string) (bid *Bid, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bidId, _ := entity.StringToErrandId(id)

	filter := bson.M{
		"_id": bidId,
	}

	if err = r.FindOne(ctx, filter).Decode(&bid); err != nil {
		return nil, err
	}
	return bid, nil
}

func (r *repository) GetBidWithRunnerAndErrandId(errandId string, runnerId string) (bid *Bid, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$and": []bson.M{
			{"errand_id": errandId},
			{"runner": runnerId},
			{"bid_state": Open},
		},
	}
	if err = r.FindOne(ctx, filter).Decode(&bid); err != nil {
		return nil, err
	}
	return bid, nil
}

func (r *repository) GetBidsBetween(errandId, firstId, secondId string) (bids []Bid, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$and": []bson.M{
			{"errand_id": errandId},
			{
				"$or": []bson.M{
					{"bidder": firstId},
					{"bidder": secondId},
				},
			},
		},
	}
	sort := options.Find().SetSort(bson.D{
		{"created_at", -1},
	})

	crs, err := r.Find(ctx, filter, sort)
	if err != nil {
		return nil, err
	}
	for crs.Next(ctx) {
		var bid Bid
		crsErr := crs.Decode(&bid)
		if crsErr != nil {
			logger.Error("fetching bids", crsErr)
		}
		bids = append(bids, bid)
	}

	return bids, nil
}

func (r *repository) UpdateBidHaggle(id string, haggle *haggle.Haggle) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bidId, _ := entity.StringToErrandId(id)
	filter := bson.M{
		"_id": bidId,
	}
	param := bson.M{
		"$push": bson.M{
			"haggles": haggle,
		},
	}
	_, err := r.Collection.UpdateOne(ctx, filter, param)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetAllBids(errandId string) (bids []Bid, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$and": []bson.M{
			{"errand_id": errandId},
			{"bid_state": Open},
		},
	}
	sort := options.Find().SetSort(bson.D{
		{"created_at", -1},
	})

	crs, err := r.Find(ctx, filter, sort)
	if err != nil {
		return nil, err
	}
	for crs.Next(ctx) {
		var bid Bid
		crsErr := crs.Decode(&bid)
		if crsErr != nil {
			logger.Error("fetching bids", crsErr)
		}
		bids = append(bids, bid)
	}

	return bids, nil
}

func (r *repository) CreateBid(bid *Bid) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.InsertOne(ctx, bid)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateBid(bid *Bid) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": bid.Id,
	}

	param := bson.D{
		{"$set", bson.D{
			{"bid_state", bid.BidState},
			{"state", bid.State},
			{"updated_at", bid.UpdatedAt},
		}},
	}

	_, err := r.Collection.UpdateOne(ctx, filter, param)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteBidsFor(errandId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$and": []bson.M{
			{"errand_id": errandId},
			{"$not": []bson.M{
				{"is_primary_bid": true},
			}},
		},
	}

	_, err := r.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
