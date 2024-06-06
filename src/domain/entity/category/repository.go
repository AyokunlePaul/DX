package category

import (
	"DX/src/domain/entity"
	"DX/src/utils/logger"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type reader interface {
	Get(string) (*Category, error)
	GetAll() ([]Category, error)
	GetTopCategories(int) ([]Category, error)
}

type writer interface {
	Create(*Category) error
	Update(string, string, string, string) error
	Delete(string) error
}

type Repository interface {
	reader
	writer
}

type repository struct {
	*mongo.Collection
}

func NewRepository(collection *mongo.Collection) Repository {
	return &repository{Collection: collection}
}

func (r *repository) Get(id string) (category *Category, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	catId, _ := entity.StringToErrandId(id)

	filter := bson.M{
		"$or": []bson.M{
			{"_id": catId},
			{"identifier": id},
		},
	}

	if err = r.Collection.FindOne(ctx, filter).Decode(&category); err != nil {
		return nil, err
	}

	return category, nil
}

func (r *repository) GetAll() (categories []Category, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	crs, err := r.Collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var nCategory Category
	for crs.Next(ctx) {
		crsErr := crs.Decode(&nCategory)
		if crsErr != nil {
			logger.Error("fetching categories", crsErr)
			return nil, crsErr
		}
		categories = append(categories, nCategory)
	}

	return categories, nil
}

func (r *repository) GetTopCategories(limit int) (categories []Category, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	sort := bson.M{
		"created_at": 1,
	}
	opts := options.Find().SetSort(sort).SetLimit(int64(limit))

	crs, err := r.Collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	var nCategory Category
	for crs.Next(ctx) {
		crsErr := crs.Decode(&nCategory)
		if crsErr != nil {
			logger.Error("fetching categories", crsErr)
			return nil, crsErr
		}
		categories = append(categories, nCategory)
	}

	return categories, nil
}

func (r *repository) Create(category *Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if _, err := r.Collection.InsertOne(ctx, category); err != nil {
		return err
	}

	return nil
}

func (r *repository) Update(id, name, identifier, adminId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	catId, _ := entity.StringToErrandId(id)

	filter := bson.M{
		"_id": catId,
	}
	param := bson.M{
		"$set": bson.M{
			"name":       name,
			"identifier": identifier,
		},
		"$push": bson.M{
			"modified_by": entity.ModifiedBy{
				Id:   adminId,
				Date: time.Now(),
			},
		},
	}

	if _, err := r.Collection.UpdateOne(ctx, filter, param); err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	catId, _ := entity.StringToErrandId(id)

	filter := bson.M{
		"_id": catId,
	}

	if _, err := r.Collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
