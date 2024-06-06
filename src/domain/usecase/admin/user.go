package admin

import (
	"DX/src/domain/entity/auth"
	"DX/src/domain/entity/user"
	"DX/src/pkg/error_service"
	"DX/src/utils/logger"
	"errors"
	"fmt"
	"strings"
	"time"
)

type UserUseCase interface {
	GetUser(string) (*user.User, error)
	GetAllUsers(string) ([]user.User, error)
	CreateNewUser(string, *user.User) error
	UpdateUser(string, string, *user.User) error
	SuspendUser(string, string) error
	SuspendUsers(string, []string) error
	RestoreUser(string, string) error
	DeleteUser(string, string) error
	DeleteUsers(string, []string) error
}

type userImpl struct {
	auth.Manager
	user.Repository
	error_service.Service
}

func NewUserUseCase(manager auth.Manager, repository user.Repository, service error_service.Service) UserUseCase {
	return &userImpl{
		Manager:    manager,
		Repository: repository,
		Service:    service,
	}
}

func (i *userImpl) CreateNewUser(userToken string, user *user.User) error {
	adminUserId, resp := i.Manager.Get(userToken)
	if resp != nil {
		return errors.New(resp.Message)
	}

	token := auth.NewUserToken()
	if resp := i.Manager.CreateToken(token, user.Id.Hex()); resp != nil {
		return errors.New(resp.Message)
	}
	user.Token = token.AccessToken
	user.RefreshToken = token.RefreshToken
	user.AdminId = *adminUserId

	err := i.Repository.Create(user)
	if err != nil {
		logger.Error(fmt.Sprintf("value: %v", user), err)
		return errors.New(i.HandleMongoDbError("user", err).Message)
	}
	return nil
}

func (i *userImpl) UpdateUser(token string, userId string, user *user.User) error {
	adminUserId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	user.UpdatedAt = time.Now()
	oUser, err := i.Repository.GetWithId(userId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("user", err).Message)
	}

	oUser.UpdateUserDataForAdmin(user, *adminUserId)

	err = i.Repository.Update(oUser)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("user", err).Message)
	}

	return nil
}

func (i *userImpl) GetUser(userId string) (*user.User, error) {
	nUser, err := i.Repository.GetWithId(userId)
	if err != nil {
		return nil, errors.New(i.Service.HandleMongoDbError("user", err).Message)
	}

	return nUser, nil
}

func (i *userImpl) GetAllUsers(filter string) (users []user.User, err error) {
	if strings.TrimSpace(filter) == "suspended" {
		users, err = i.Repository.GetAllSuspendedUsers()
	} else if strings.TrimSpace(filter) == "deleted" {
		users, err = i.Repository.GetAllDeletedUsers()
	} else {
		users, err = i.Repository.GetAllUsers()
	}
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (i *userImpl) SuspendUser(token string, userId string) error {
	adminUserId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	err := i.Repository.Suspend(userId, *adminUserId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("user", err).Message)
	}

	return nil
}

func (i *userImpl) SuspendUsers(token string, ids []string) error {
	adminUserId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	if err := i.Repository.SuspendMany(*adminUserId, ids); err != nil {
		return errors.New(i.Service.HandleMongoDbError("users", err).Message)
	}

	return nil
}

func (i *userImpl) RestoreUser(token string, userId string) error {
	adminUserId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	err := i.Repository.Restore(userId, *adminUserId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("user", err).Message)
	}

	return nil
}

func (i *userImpl) DeleteUser(token string, userId string) error {
	adminUserId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	err := i.Repository.Delete(*adminUserId, userId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("user", err).Message)
	}

	return nil
}

func (i *userImpl) DeleteUsers(token string, users []string) error {
	adminUserId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	if err := i.Repository.DeleteMany(*adminUserId, users); err != nil {
		return errors.New(i.Service.HandleMongoDbError("users", err).Message)
	}

	return nil
}
