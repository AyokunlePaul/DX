package auth

import (
	"DX/src/pkg/error_service"
	"DX/src/pkg/response"
	"DX/src/utils/logger"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type mongoRepository struct {
	Collection *mongo.Collection
	error_service.Service
}

func NewMongoCacheRepository(collection *mongo.Collection) Repository {
	return &mongoRepository{Collection: collection}
}

func (mr *mongoRepository) CreateToken(token *Token, userId string) *response.BaseResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token.UserId = userId
	_, err := mr.Collection.InsertOne(ctx, token)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("value: %v", token))
		logger.GetLogger().Error(fmt.Sprintf("couldn't create token. error: %v", err))
		return response.NewInternalServerError("unable to process request")
	}
	return nil
}

func (mr *mongoRepository) Get(tokenUuid string) (*string, *response.BaseResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{
		"access_uuid": tokenUuid,
	}
	token := &Token{}
	if err := mr.Collection.FindOne(ctx, filter).Decode(token); err != nil {
		logger.Error(fmt.Sprintf("error fetching token: %v", tokenUuid), err)
		return nil, mr.HandleMongoDbError("token", err)
	}

	return &token.UserId, nil
}

func (mr *mongoRepository) DeleteToken(tokenUuid string) *response.BaseResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"access_uuid": tokenUuid,
	}
	if _, err := mr.Collection.DeleteOne(ctx, filter); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("couldn't delete token. error: %v", err))
		return response.NewInternalServerError("error deleting user")
	}
	return nil
}

func (mr *mongoRepository) RefreshToken(crudPayToken *Token, userId string, oldTokenUuid string) *response.BaseResponse {
	return nil
}
