package authentication

import (
	"DX/src/domain/entity/notification"
	"DX/src/domain/entity/user"
	"DX/src/pkg/response"
)

type UseCase interface {
	ValidatePhoneNumber(string) *response.BaseResponse
	CreateUser(*user.User) *response.BaseResponse
	LoginUser(*user.User) *response.BaseResponse
	LogOut(string) *response.BaseResponse
	ResetUserPassword(*user.User) *response.BaseResponse
	SuspendUser(string) *response.BaseResponse
	DeleteUserAccount(*user.User) *response.BaseResponse
	GetUserProfile(string) (*user.User, error)
	GetUser(string) (*user.User, error)
	UpdateUserPassword(string, string) error
	GetNotifications(string) ([]notification.Notification, error)
}
