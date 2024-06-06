package error_service

import "DX/src/pkg/response"

type Service interface {
	HandleMongoDbError(string, error) *response.BaseResponse
	HandleRedisDbError(error) *response.BaseResponse
	HandleGoogleStorageError(error) *response.BaseResponse
	HandleUtilityError(error) *response.BaseResponse
}
