package authentication

import (
	"DX/src/domain/entity"
	"DX/src/domain/entity/auth"
	"DX/src/domain/entity/notification"
	"DX/src/domain/entity/user"
	"DX/src/pkg/error_service"
	"DX/src/pkg/password_service"
	"DX/src/pkg/response"
	"DX/src/utils"
	"DX/src/utils/logger"
	"errors"
	"fmt"
)

type impl struct {
	repository user.Repository
	error_service.Service
	password         password_service.Service
	manager          auth.Manager
	NotificationRepo notification.Repository
}

func NewUseCase(repo user.Repository, errorHandler error_service.Service, passwordService password_service.Service, authManager auth.Manager,
	notificationRepo notification.Repository) UseCase {
	return &impl{
		repository:       repo,
		Service:          errorHandler,
		password:         passwordService,
		manager:          authManager,
		NotificationRepo: notificationRepo,
	}
}

func (i *impl) ValidatePhoneNumber(phone string) *response.BaseResponse {
	if !utils.IsValidPhoneNumber(phone) {
		return response.NewBadRequestError("invalid phone number")
	}
	err := i.repository.Phone(phone)
	if err != nil {
		return i.Service.HandleMongoDbError("phone", err)
	}

	return nil
}

func (i *impl) CreateUser(user *user.User) *response.BaseResponse {
	if err := user.IsValidForInitialCreation(); err != nil {
		return err
	}

	token := auth.NewUserToken()
	if err := i.manager.CreateToken(token, user.Id.Hex()); err != nil {
		return err
	}
	user.Token = token.AccessToken
	user.RefreshToken = token.RefreshToken

	hash, err := i.password.Generate(user.Password)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("error creating user: %v", err))
		return response.NewInternalServerError("error creating user")
	}
	user.Password = hash
	user.HasVerifiedPhone = true

	err = i.repository.Create(user)
	if err != nil {
		logger.Error(fmt.Sprintf("value: %v", user), err)
		return i.Service.HandleMongoDbError("user", err)
	}
	return nil
}

func (i *impl) LoginUser(user *user.User) *response.BaseResponse {
	password := user.Password

	err := i.repository.Get(user)
	if err != nil {
		return i.Service.HandleMongoDbError("user", err)
	}

	err = i.password.Compare(user.Password, password)
	if err != nil {
		return response.NewBadRequestError("invalid login credentials")
	}

	token := auth.NewUserToken()
	if err := i.manager.CreateToken(token, user.Id.Hex()); err != nil {
		return err
	}
	user.Token = token.AccessToken
	user.RefreshToken = token.RefreshToken

	return nil
}

func (i *impl) SuspendUser(token string) *response.BaseResponse {
	nUser := new(user.User)
	userId, resp := i.manager.Get(token)
	if resp != nil {
		return resp
	}
	entId, err := entity.StringToErrandId(*userId)
	if err != nil {
		return response.NewBadRequestError(err.Error())
	}

	nUser.Id = entId
	err = i.repository.Get(nUser)
	if err != nil {
		return i.Service.HandleMongoDbError("user", err)
	}

	nUser.IsSuspended = true

	err = i.repository.Update(nUser)
	if err != nil {
		return i.Service.HandleMongoDbError("user", err)
	}

	return nil
}

func (i *impl) LogOut(token string) *response.BaseResponse {
	resp := i.manager.DeleteToken(token)
	if resp != nil {
		return resp
	}

	return nil
}

func (i *impl) GetUserProfile(token string) (*user.User, error) {
	userId, resp := i.manager.Get(token)
	if resp != nil {
		return nil, errors.New(resp.Message)
	}

	nUser, err := i.repository.GetWithId(*userId)
	if err != nil {
		return nil, err
	}
	return nUser, nil
}

func (i *impl) GetUser(userId string) (*user.User, error) {
	nUser, err := i.repository.GetWithId(userId)
	if err != nil {
		return nil, err
	}
	return nUser, nil
}

func (i *impl) UpdateUserPassword(phone string, password string) error {
	nUser, err := i.repository.GetWithPhone(phone)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("user", err).Message)
	}

	hash, err := i.password.Generate(password)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("error creating user: %v", err))
		return errors.New("error creating user")
	}

	nUser.Password = hash
	err = i.repository.Update(nUser)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("user", err).Message)
	}

	return nil
}

func (i *impl) GetNotifications(token string) ([]notification.Notification, error) {
	userId, resp := i.manager.Get(token)
	if resp != nil {
		return nil, errors.New(resp.Message)
	}

	nNotifications, err := i.NotificationRepo.GetAllNotifications(*userId)
	if err != nil {
		return nil, errors.New(i.Service.HandleMongoDbError("notification", err).Message)
	}

	return nNotifications, nil
}

func (i *impl) ResetUserPassword(user *user.User) *response.BaseResponse {
	return nil
}

func (i *impl) DeleteUserAccount(user *user.User) *response.BaseResponse {
	return nil
}
