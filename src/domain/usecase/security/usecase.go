package security

import (
	"DX/src/domain/entity/auth"
	"DX/src/domain/entity/security"
	"DX/src/domain/entity/user"
	"DX/src/pkg/error_service"
	"DX/src/utils"
	"errors"
)

type UseCase interface {
	GetSecurity(string) (string, error)
	GetSecurityWithPhone(string) (string, error)
	SaveSecurity(string, *security.Security) error
	Validate(string, string) error
}

type impl struct {
	auth.Manager
	security.Repository
	UserRepo user.Repository
	error_service.Service
}

func NewUseCase(manager auth.Manager, repository security.Repository, service error_service.Service, UserRepo user.Repository) UseCase {
	return &impl{
		Manager:    manager,
		Repository: repository,
		Service:    service,
		UserRepo:   UserRepo,
	}
}

func (i *impl) GetSecurity(token string) (string, error) {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return "", errors.New("unable to get security question")
	}

	sec, err := i.Repository.Get(*userId)
	if err != nil {
		mongoResp := i.Service.HandleMongoDbError("security question", err)
		return "", errors.New(mongoResp.Message)
	}

	return sec.Question, nil
}

func (i *impl) GetSecurityWithPhone(phone string) (string, error) {
	if !utils.IsValidPhoneNumber(phone) {
		return "", errors.New("invalid phone number")
	}
	sec, err := i.Repository.Get(phone)
	if err != nil {
		mongoResp := i.Service.HandleMongoDbError("security question", err)
		return "", errors.New(mongoResp.Message)
	}

	return sec.Question, nil
}

func (i *impl) SaveSecurity(token string, sec *security.Security) error {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New("unable to save security questions")
	}

	nUser, err := i.UserRepo.GetWithId(*userId)
	if err != nil {
		mongoResp := i.Service.HandleMongoDbError("security question", err)
		return errors.New(mongoResp.Message)
	}

	sec.UserId = nUser.Id.Hex()
	sec.PhoneNumber = nUser.PhoneNumber

	err = i.Repository.Create(sec)
	if err != nil {
		mongoResp := i.Service.HandleMongoDbError("security question", err)
		return errors.New(mongoResp.Message)
	}
	return nil
}

func (i *impl) Validate(phone string, answer string) error {
	if !utils.IsValidPhoneNumber(phone) {
		return errors.New("invalid phone number")
	}
	sec, err := i.Repository.Get(phone)
	if err != nil {
		mongoResp := i.Service.HandleMongoDbError("security question", err)
		return errors.New(mongoResp.Message)
	}

	if sec.Answer != answer {
		return errors.New("wrong answer to security question")
	}

	return nil
}
