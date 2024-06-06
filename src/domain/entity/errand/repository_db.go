package errand

import (
	"DX/src/domain/entity"
	"DX/src/domain/entity/bid"
	"DX/src/domain/entity/haggle"
	"DX/src/domain/entity/timeline"
	"DX/src/pkg/error_service"
	"DX/src/utils/logger"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type repository struct {
	*mongo.Collection
}

func NewRepository(collection *mongo.Collection) Repository {
	return &repository{Collection: collection}
}

func (r *repository) GetDraft(id string) (errand *Errand, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info(fmt.Sprintf("user id for draft: %s", id))
	filter := bson.M{
		"$and": []bson.M{
			{"user_id": id},
			{"state": Draft},
		},
	}
	err = r.Collection.FindOne(ctx, filter).Decode(&errand)
	// If errand does not have a draft PR, return nothing
	if err == mongo.ErrNoDocuments {
		logger.Error("no draft errand error", err)
		return nil, nil
	}
	if err != nil {
		logger.Error("error decoding draft errand", err)
		return nil, err
	}

	return errand, nil
}

func (r *repository) Get(id string) (errand *Errand, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errandId, _ := entity.StringToErrandId(id)

	filter := bson.M{
		"_id": errandId,
	}

	err = r.Collection.FindOne(ctx, filter).Decode(&errand)
	if err != nil {
		return nil, err
	}

	return errand, nil
}

func (r *repository) GetAllMarketErrands() (errands []Errand, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"state": Pending},
			{"state": Open},
		},
	}

	crs, err := r.Collection.Find(ctx, filter)
	for crs.Next(ctx) {
		var nErrand Errand
		crsErr := crs.Decode(&nErrand)
		if crsErr != nil {
			logger.Error("fetching errands", crsErr)
			return nil, err
		}
		errands = append(errands, nErrand)
	}
	return errands, nil
}

func (r *repository) GetAll() (errands []Errand, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	filter := bson.M{}

	crs, err := r.Collection.Find(ctx, filter)
	if err = crs.All(ctx, &errands); err != nil {
		return nil, err
	}

	return errands, nil
}

func (r *repository) GetAllDraftErrands() (errands []Errand, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	filter := bson.M{
		"state": Draft,
	}

	crs, err := r.Collection.Find(ctx, filter)
	for crs.Next(ctx) {
		var nErrand Errand
		crsErr := crs.Decode(&nErrand)
		if crsErr != nil {
			logger.Error("fetching errands", crsErr)
			return nil, err
		}
		errands = append(errands, nErrand)
	}
	return errands, nil
}

func (r *repository) GetAllCompletedErrands() (errands []Errand, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	filter := bson.M{
		"state": Completed,
	}

	crs, err := r.Collection.Find(ctx, filter)
	for crs.Next(ctx) {
		var nErrand Errand
		crsErr := crs.Decode(&nErrand)
		if crsErr != nil {
			logger.Error("fetching errands", crsErr)
			return nil, err
		}
		errands = append(errands, nErrand)
	}
	return errands, nil
}

func (r *repository) GetAllCancelledErrands() (errands []Errand, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	filter := bson.M{
		"state": Cancelled,
	}

	crs, err := r.Collection.Find(ctx, filter)
	for crs.Next(ctx) {
		var nErrand Errand
		crsErr := crs.Decode(&nErrand)
		if crsErr != nil {
			logger.Error("fetching errands", crsErr)
			return nil, err
		}
		errands = append(errands, nErrand)
	}
	return errands, nil
}

func (r *repository) GetAllActiveErrands() (errands []Errand, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	filter := bson.M{
		"state": Active,
	}

	crs, err := r.Collection.Find(ctx, filter)
	for crs.Next(ctx) {
		var nErrand Errand
		crsErr := crs.Decode(&nErrand)
		if crsErr != nil {
			logger.Error("fetching errands", crsErr)
			return nil, err
		}
		errands = append(errands, nErrand)
	}
	return errands, nil
}

func (r *repository) GetAllAbandonedErrands() (errands []Errand, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	filter := bson.M{
		"state": Abandoned,
	}

	crs, err := r.Collection.Find(ctx, filter)
	for crs.Next(ctx) {
		var nErrand Errand
		crsErr := crs.Decode(&nErrand)
		if crsErr != nil {
			logger.Error("fetching errands", crsErr)
			return nil, err
		}
		errands = append(errands, nErrand)
	}
	return errands, nil
}

func (r *repository) GetBidForUser(id string, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	errandId, _ := entity.StringToErrandId(id)

	filter := bson.M{
		"$and": []bson.M{
			{"_id": errandId},
			{"bids": bson.M{
				"$elemMatch": bson.M{
					"$and": []bson.M{
						{"runner": userId},
						{"bid_state": bid.Open},
					},
				},
			}},
		},
	}

	res := r.Collection.FindOne(ctx, filter)

	return res.Err()
}

func (r *repository) GetFor(userId string) (errands []Errand, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	filter := bson.M{
		"$and": []bson.M{
			{"$or": []bson.M{
				{"user_id": userId},
				{"runner_id": userId},
				{"bids.runner": userId},
			}},
			{"state": bson.M{
				"$ne": Draft,
			}},
		},
	}

	crs, err := r.Collection.Find(ctx, filter)
	if err = crs.All(ctx, &errands); err != nil {
		return nil, err
	}

	return errands, nil
}

func (r *repository) AssignErrandToSender(adminId string, eId string, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errandId, _ := entity.StringToErrandId(eId)
	cTime := time.Now()

	filter := bson.M{
		"_id": errandId,
	}

	param := bson.D{
		{"$set", bson.D{
			{"user_id", userId},
			{"state", Open},
			{"status", Open.Id()},
			{"updated_at", cTime},
		}},
		{"$push", bson.D{
			{"modified_by", entity.ModifiedBy{
				Id:   adminId,
				Date: cTime,
			}},
		}},
	}

	if _, err := r.Collection.UpdateOne(ctx, filter, param); err != nil {
		return err
	}

	return nil
}

func (r *repository) AssignErrandToRunner(adminId, eId, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errandId, _ := entity.StringToErrandId(eId)
	cTime := time.Now()

	filter := bson.M{
		"_id": errandId,
	}

	param := bson.D{
		{"$set", bson.D{
			{"runner_id", userId},
			{"state", Pending},
			{"status", Pending.Id()},
			{"updated_at", cTime},
		}},
		{"$push", bson.D{
			{"modified_by", entity.ModifiedBy{
				Id:   adminId,
				Date: cTime,
			}},
		}},
	}

	if _, err := r.Collection.UpdateOne(ctx, filter, param); err != nil {
		return err
	}

	return nil
}

func (r *repository) AssignErrandToOfflineRunner(adminId, eId, userId string, bid *bid.Bid) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errandId, _ := entity.StringToErrandId(eId)
	cTime := time.Now()

	filter := bson.M{
		"_id": errandId,
	}

	param := bson.D{
		{"$set", bson.D{
			{"runner_id", userId},
			{"state", Active},
			{"status", Active.Id()},
			{"updated_at", cTime},
		}},
		{"$push", bson.D{
			{"bids", bid},
			{"modified_by", entity.ModifiedBy{
				Id:   adminId,
				Date: cTime,
			}},
		}},
	}

	if _, err := r.Collection.UpdateOne(ctx, filter, param); err != nil {
		return err
	}

	return nil
}

func (r *repository) AddBidToErrand(id, userId string, bid *bid.Bid) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errandId, _ := entity.StringToErrandId(id)
	cTime := time.Now()

	filter := bson.M{
		"_id": errandId,
	}
	param := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
		"$inc": bson.M{
			"total_bids": 1,
		},
		"$push": bson.M{
			"bids": bid,
			"modified_by": entity.ModifiedBy{
				Id:   userId,
				Date: cTime,
			},
		},
	}

	if _, err := r.Collection.UpdateOne(ctx, filter, param); err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateBidHaggle(eId, bId string, haggle *haggle.Haggle) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errandId, _ := entity.StringToErrandId(eId)
	bidId, _ := entity.StringToErrandId(bId)

	filter := bson.M{
		"_id":      errandId,
		"bids._id": bidId,
	}

	param := bson.M{
		"$push": bson.M{
			"bids.$.haggles": haggle,
		},
	}

	if _, err := r.Collection.UpdateOne(ctx, filter, param); err != nil {
		return err
	}

	return nil
}

func (r *repository) AcceptBid(eId, bId, senderId string, amount int64, update timeline.Update) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errandId, _ := entity.StringToErrandId(eId)
	bidId, _ := entity.StringToErrandId(bId)
	cTime := time.Now()

	filter := bson.M{
		"$and": []bson.M{
			{"_id": errandId},
			{"bids._id": bidId},
		},
	}
	param := bson.D{
		{"$set", bson.D{
			{"state", Pending},
			{"status", Pending.Id()},
			{"amount", amount},
			{"bids.$.bid_state", bid.Accepted},
			{"bids.$.state", bid.Accepted.Id()},
			{"updated_at", cTime},
			{"timeline.updated_at", cTime},
		}},
		{"$push", bson.D{
			{"timeline.updates", update},
			{"modified_by", entity.ModifiedBy{
				Id:   senderId,
				Date: cTime,
			}},
		}},
	}

	if res, err := r.Collection.UpdateOne(ctx, filter, param); err != nil {
		return err
	} else {
		if res.ModifiedCount == 0 {
			return error_service.ErrBidAcceptance
		}
	}

	return nil
}

func (r *repository) RejectBid(eId string, bId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	errandId, _ := entity.StringToErrandId(eId)
	bidId, _ := entity.StringToErrandId(bId)
	cTime := time.Now()

	filter := bson.M{
		"$and": []bson.M{
			{"_id": errandId},
			{"bids._id": bidId},
		},
	}

	param := bson.D{
		{"$set", bson.D{
			{"bids.$.bid_state", bid.Rejected},
			{"bids.$.state", bid.Rejected.Id()},
			{"bids.$.updated_at", cTime},
			{"updated_at", cTime},
		}},
	}

	if _, err := r.Collection.UpdateOne(ctx, filter, param); err != nil {
		return err
	}

	return nil
}

func (r *repository) StartErrand(eId, runnerId string, update timeline.Update) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	errandId, _ := entity.StringToErrandId(eId)
	cTime := time.Now()

	filter := bson.M{
		"$and": []bson.M{
			{"_id": errandId},
			{"bids.state": bid.Accepted.Id()},
		},
	}
	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{"bidElem.state", bid.Open.Id()}},
		},
	})

	param := bson.D{
		{"$set", bson.D{
			{"runner_id", runnerId},
			{"state", Active},
			{"status", Active.Id()},
			{"bids.$[bidElem].bid_state", bid.Rejected},
			{"bids.$[bidElem].state", bid.Rejected.Id()},
			{"bids.$[bidElem].updated_at", cTime},
			{"timeline.updated_at", cTime},
			{"updated_at", cTime},
		}},
		{"$push", bson.D{
			{"timeline.updates", update},
		}},
	}

	if _, err := r.Collection.UpdateOne(ctx, filter, param, opts); err != nil {
		return err
	}

	return nil
}

func (r *repository) ResetErrandBids(eId, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	errandId, _ := entity.StringToErrandId(eId)
	cTime := time.Now()

	filter := bson.M{
		"_id": errandId,
	}
	param := bson.D{
		{"$set", bson.D{
			{"state", Open},
			{"status", Open.Id()},
			{"bids.$[].bid_state", bid.Open},
			{"bids.$[].state", bid.Open.Id()},
			{"bids.$[].updated_at", cTime},
			{"updated_at", cTime},
			{"runner_id", ""},
			{"amount", 0},
		}},
		{"$push", bson.D{
			{"modified_by", entity.ModifiedBy{
				Id:   userId,
				Date: cTime,
			}},
		}},
	}

	if _, err := r.Collection.UpdateOne(ctx, filter, param); err != nil {
		return err
	}

	return nil
}

func (r *repository) RunnerComplete(eId string, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	errandId, _ := entity.StringToErrandId(eId)
	cTime := time.Now()

	filter := bson.M{
		"_id": errandId,
	}

	param := bson.D{
		{"$set", bson.D{
			{"state", Completed},
			{"status", Completed.Id()},
			{"updated_at", cTime},
		}},
		{"$push", bson.D{
			{"modified_by", entity.ModifiedBy{
				Id:   userId,
				Date: cTime,
			}},
		}},
	}

	if _, err := r.Collection.UpdateOne(ctx, filter, param); err != nil {
		return err
	}

	return nil
}

func (r *repository) SenderComplete(eId string, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	errandId, _ := entity.StringToErrandId(eId)
	cTime := time.Now()

	filter := bson.M{
		"_id": errandId,
	}

	param := bson.D{
		{"$set", bson.D{
			{"state", RunnerCompleted},
			{"status", RunnerCompleted.Id()},
			{"updated_at", cTime},
		}},
		{"$push", bson.D{
			{"modified_by", entity.ModifiedBy{
				Id:   userId,
				Date: cTime,
			}},
		}},
	}

	if _, err := r.Collection.UpdateOne(ctx, filter, param); err != nil {
		return err
	}

	return nil
}

func (r *repository) Search(keyword string) ([]string, error) {
	return nil, nil
}

func (r *repository) Create(errand *Errand) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := r.Collection.InsertOne(ctx, errand)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Update(errand *Errand) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	param := bson.D{
		{"$set", bson.D{
			{"description", errand.Description},
			{"category", errand.Category},
			{"status", errand.Status},
			{"state", errand.State},
			{"duration", errand.Duration},
			{"images", errand.Images},
			{"audio", errand.Audio},
			{"step", errand.Step},
			{"restriction", errand.Restriction},
			{"restrict_by", errand.RestrictBy},
			{"pickup_address", errand.PickupAddress},
			{"dropoff_address", errand.DropOffAddress},
			{"budget", errand.Budget},
			{"expiry_date", errand.ExpiryDate},
			{"created_by", errand.CreatedBy},
			{"created_at", errand.CreatedAt},
			{"updated_at", errand.UpdatedAt},
			{"cancellation_reason", errand.CancellationReason},
			{"modified_by", errand.ModifiedBy},
			{"runner_id", errand.RunnerId},
			{"timeline", errand.Timeline},
		}},
	}
	filter := bson.M{"_id": errand.Id}
	if _, err := r.Collection.UpdateOne(ctx, filter, param); err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateTimeline(eId, userId string, update timeline.Update) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	errandId, _ := entity.StringToErrandId(eId)

	filter := bson.M{
		"_id": errandId,
	}
	param := bson.D{
		{"$set", bson.D{
			{"updated_at", time.Now()},
		}},
		{"$push", bson.D{
			{"timeline.updates", update},
			{"modified_by", entity.ModifiedBy{
				Id:   userId,
				Date: time.Now(),
			}},
		}},
	}

	if _, err := r.Collection.UpdateOne(ctx, filter, param); err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}
