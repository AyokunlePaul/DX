package error_service

import (
	"DX/src/pkg/response"
	"DX/src/utils/logger"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type errandError struct{}

func New() Service {
	return &errandError{}
}

func (errandErr *errandError) HandleMongoDbError(from string, err error) *response.BaseResponse {
	logger.Error("mongo error:", err)
	if writeException, ok := err.(mongo.WriteException); ok {
		for _, exception := range writeException.WriteErrors {
			switch exception.Code {
			case 11000:
				return response.NewBadRequestError(fmt.Sprintf("%s already exist", from))
			}
		}
		return response.NewBadRequestError(fmt.Sprintf("%s already exist", from))
	}
	switch err {
	case mongo.ErrNoDocuments:
		return response.NewNotFoundError(fmt.Sprintf("%s doesn't exist", from))
	case mongo.ErrClientDisconnected:
		return response.NewInternalServerError("internal server error")
	case mongo.ErrNilDocument:
		return response.NewBadRequestError("internal server error")
	case ErrDuplicatePhoneNumber:
		return response.NewBadRequestError("phone number already exist")
	case ErrNoUser:
		return response.NewBadRequestError("user doesn't exist")
	default:
		return response.NewInternalServerError(err.Error())
	}
}

func (errandErr *errandError) HandleRedisDbError(err error) *response.BaseResponse {
	logger.Error("redis error", err)
	switch err {
	case redis.Nil:
		return response.NewUnAuthorizedError()
	default:
		return response.NewInternalServerError(err.Error())
	}
}

func (errandErr *errandError) HandleGoogleStorageError(err error) *response.BaseResponse {
	logger.Error("google storage error", err)
	return response.NewInternalServerError(fmt.Sprintf("an error occurred: %s", err.Error()))
}

func (errandErr *errandError) HandleUtilityError(err error) *response.BaseResponse {
	logger.Error("utility error", err)
	return response.NewInternalServerError(fmt.Sprintf("an error occurred: %s", err.Error()))
}
