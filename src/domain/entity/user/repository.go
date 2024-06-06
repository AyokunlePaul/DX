package user

import (
	"DX/src/domain/entity"
	"DX/src/pkg/error_service"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type reader interface {
	Phone(string) error
	Get(*User) error
	GetWithId(string) (*User, error)
	GetWithPhone(string) (*User, error)
	GetAllUsers() ([]User, error)
	GetAllDeletedUsers() ([]User, error)
	GetAllSuspendedUsers() ([]User, error)
}

type writer interface {
	Create(*User) error
	Update(*User) error
	CompleteErrand(string) error
	RateUser(string, int64) error
	Suspend(string, string) error
	SuspendMany(string, []string) error
	Restore(string, string) error
	Delete(string, string) error
	DeleteMany(string, []string) error
}

type Repository interface {
	reader
	writer
}

type repository struct {
	Collection *mongo.Collection
}

func NewDatabaseRepository(collection *mongo.Collection) Repository {
	return &repository{Collection: collection}
}

func (r *repository) Phone(phone string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{
		"phone_number": phone,
	}
	res := r.Collection.FindOne(ctx, filter)
	// If the number exist, return an error
	if res.Err() == nil {
		return error_service.ErrDuplicatePhoneNumber
	}
	// If the phone doesn't exist, return nil
	if res.Err() == mongo.ErrNoDocuments {
		return nil
	}
	return res.Err()
}

func (r *repository) Get(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"_id": user.Id},
			{"phone_number": user.PhoneNumber},
		},
	}
	if err := r.Collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return err
	}
	if user.IsSuspended || user.IsDeleted {
		return errors.New("user is suspended")
	}
	return nil
}

func (r *repository) GetAllUsers() (users []User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{
		{"$and", []bson.M{
			{"is_deleted": bson.M{
				"$ne": true},
			},
			{"is_suspended": bson.M{
				"$ne": true},
			}},
		},
	}

	crs, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var user User
	for crs.Next(ctx) {
		err = crs.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *repository) GetAllDeletedUsers() (users []User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{
		{"is_deleted", true},
	}

	crs, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var user User
	for crs.Next(ctx) {
		err = crs.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *repository) GetAllSuspendedUsers() (users []User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{
		{"is_suspended", true},
	}

	crs, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var user User
	for crs.Next(ctx) {
		err = crs.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *repository) GetWithId(userId string) (user *User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, _ := entity.StringToErrandId(userId)

	filter := bson.M{
		"_id": id,
	}
	if err := r.Collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}
	if user.IsSuspended || user.IsDeleted {
		return nil, errors.New("user is suspended")
	}
	return user, nil
}

func (r *repository) GetWithPhone(phone string) (user *User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"phone_number": phone,
	}
	if err := r.Collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) Create(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.Collection.InsertOne(ctx, user)
	return err
}

func (r *repository) CompleteErrand(uId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId, _ := entity.StringToErrandId(uId)

	filter := bson.D{
		{"_id", userId},
	}
	param := bson.D{
		{"$inc", bson.D{
			{"errands_completed", 1},
		}},
	}

	if _, err := r.Collection.UpdateOne(ctx, filter, param); err != nil {
		return err
	}

	return nil
}

func (r *repository) RateUser(uId string, rating int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId, _ := entity.StringToErrandId(uId)

	filter := bson.D{
		{"_id", userId},
	}
	param := bson.D{
		{"$push", bson.D{
			{"ratings", rating},
		}},
		{"$addFields", bson.D{
			{"$rating", bson.D{
				{"$avg", "$ratings"},
			}},
		}},
	}

	if _, err := r.Collection.UpdateOne(ctx, filter, param); err != nil {
		return err
	}

	return nil
}

func (r *repository) Update(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	upParam := bson.D{
		{"$set", bson.D{
			{"first_name", user.FirstName},
			{"last_name", user.LastName},
			{"phone_number", user.PhoneNumber},
			{"is_admin", user.IsAdmin},
			{"updated_at", "$$NOW"},
			{"is_suspended", user.IsSuspended},
			{"is_deleted", user.IsDeleted},
			{"password", user.Password},
			{"email", user.Email},
			{"profile_picture", user.ProfilePicture},
			{"modified_by", user.ModifiedBy},
			{"has_verified_email", user.HasVerifiedEmail},
			{"has_verified_phone", user.HasVerifiedPhone},
			{"has_verified_banking_details", user.HasVerifiedBankingDetails},
			{"has_transaction_pin", user.HasTransactionPin},
			{"has_verified_address", user.HasVerifiedAddress},
			{"admin_id", user.AdminId},
			{"created_by", user.CreatedBy},
			{"ratings", user.Ratings},
			{"rating", user.Rating},
			{"errands_completed", user.ErrandsCompleted},
			{"errands_cancelled", user.ErrandsCancelled},
			{"type", user.UserType},
		}},
	}
	filter := bson.M{"_id": user.Id}
	if _, err := r.Collection.UpdateOne(ctx, filter, upParam); err != nil {
		return err
	}
	return nil
}

func (r *repository) Suspend(userId string, adminId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, _ := entity.StringToErrandId(userId)
	cTime := time.Now()

	filter := bson.M{
		"_id": id,
	}

	param := bson.D{
		{"$set", bson.D{
			{"is_suspended_by", adminId},
			{"is_suspended", true},
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

func (r *repository) SuspendMany(adminId string, ids []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cTime := time.Now()

	var bulkWrites []mongo.WriteModel

	for _, id := range ids {
		userId, _ := entity.StringToErrandId(id)
		filter := bson.M{
			"_id": userId,
		}
		param := bson.M{
			"$set": bson.M{
				"is_suspended_by": adminId,
				"is_suspended":    true,
				"updated_at":      cTime,
			},
			"$push": bson.M{
				"modified_by": entity.ModifiedBy{
					Id:   adminId,
					Date: cTime,
				}},
		}
		bulkWrites = append(bulkWrites, mongo.NewUpdateManyModel().SetFilter(filter).SetUpdate(param))
	}

	if _, err := r.Collection.BulkWrite(ctx, bulkWrites); err != nil {
		return err
	}

	return nil
}

func (r *repository) Restore(userId string, adminId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cTime := time.Now()

	id, _ := entity.StringToErrandId(userId)

	filter := bson.M{
		"_id": id,
	}

	param := bson.D{
		{"$set", bson.D{
			{"is_suspended_by", nil},
			{"is_suspended", false},
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

func (r *repository) Delete(adminId string, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId, _ := entity.StringToErrandId(id)
	cTime := time.Now()

	filter := bson.M{
		"_id": userId,
	}

	param := bson.D{
		{"$set", bson.D{
			{"is_deleted_by", adminId},
			{"is_deleted", true},
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

func (r *repository) DeleteMany(adminId string, ids []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var bulkWrites []mongo.WriteModel
	cTime := time.Now()

	for _, id := range ids {
		userId, _ := entity.StringToErrandId(id)
		filter := bson.M{
			"_id": userId,
		}
		param := bson.M{
			"$set": bson.M{
				"is_deleted_by": adminId,
				"is_deleted":    true,
				"updated_at":    cTime,
			},
			"$push": bson.M{
				"modified_by": entity.ModifiedBy{
					Id:   adminId,
					Date: cTime,
				}},
		}
		bulkWrites = append(bulkWrites, mongo.NewUpdateManyModel().SetFilter(filter).SetUpdate(param))
	}

	if _, err := r.Collection.BulkWrite(ctx, bulkWrites); err != nil {
		return err
	}

	return nil
}
