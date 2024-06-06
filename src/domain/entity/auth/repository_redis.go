package auth

import (
	"DX/src/pkg/error_service"
	"DX/src/pkg/response"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	redisContext = context.Background()
	redisClient  *redis.Client
	redisName    = "REDIS_HOST"
)

type redisRepository struct {
	errorService error_service.Service
}

func NewDatabaseRepository(errorService error_service.Service) Repository {
	return &redisRepository{
		errorService: errorService,
	}
}

func (repository *redisRepository) CreateToken(crudPayToken *Token, userId string) *response.BaseResponse {
	accessTokenExpiration := time.Unix(int64(crudPayToken.AccessTokenExpires), 0).Sub(time.Now())
	refreshTokenExpiration := time.Unix(crudPayToken.RefreshTokenExpires, 0).Sub(time.Now())

	if redisSetError := redisClient.Set(redisContext, crudPayToken.AccessUuid, userId, accessTokenExpiration).Err(); redisSetError != nil {
		return repository.errorService.HandleRedisDbError(redisSetError)
	}

	if redisSetError := redisClient.Set(redisContext, crudPayToken.RefreshUuid, userId, refreshTokenExpiration).Err(); redisSetError != nil {
		return repository.errorService.HandleRedisDbError(redisSetError)
	}

	return nil
}

func (repository *redisRepository) Get(tokenUuid string) (*string, *response.BaseResponse) {
	userId, resultError := redisClient.Get(redisContext, tokenUuid).Result()
	if resultError != nil {
		return nil, repository.errorService.HandleRedisDbError(resultError)
	}
	return &userId, nil
}

func (repository *redisRepository) DeleteToken(uuids string) *response.BaseResponse {
	redisClient.Del(redisContext, uuids)
	return nil
}

func (repository *redisRepository) RefreshToken(crudPayToken *Token, userId string, oldTokenUuid string) *response.BaseResponse {
	accessTokenExpiration := time.Unix(int64(crudPayToken.AccessTokenExpires), 0).Sub(time.Now())
	refreshTokenExpiration := time.Unix(crudPayToken.RefreshTokenExpires, 0).Sub(time.Now())

	if redisSetError := redisClient.Set(redisContext, crudPayToken.AccessUuid, userId, accessTokenExpiration).Err(); redisSetError != nil {
		return repository.errorService.HandleRedisDbError(redisSetError)
	}
	if redisDeleteError := redisClient.Del(redisContext, oldTokenUuid).Err(); redisDeleteError != nil {
		return repository.errorService.HandleRedisDbError(redisDeleteError)
	}
	if redisSetError := redisClient.Set(redisContext, crudPayToken.RefreshUuid, userId, refreshTokenExpiration).Err(); redisSetError != nil {
		return repository.errorService.HandleRedisDbError(redisSetError)
	}
	return nil
}
