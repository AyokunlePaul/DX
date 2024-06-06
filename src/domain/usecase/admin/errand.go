package admin

import (
	"DX/src/domain/entity"
	"DX/src/domain/entity/auth"
	"DX/src/domain/entity/bid"
	"DX/src/domain/entity/category"
	"DX/src/domain/entity/errand"
	"DX/src/domain/entity/notification"
	"DX/src/domain/entity/timeline"
	"DX/src/domain/entity/user"
	"DX/src/pkg/error_service"
	"errors"
	"time"
)

type ErrandUseCase interface {
	CreateErrand(string, string, *errand.Errand) error
	AssignErrand(string, string, string) error
	GetAllErrands(string) ([]errand.Errand, error)
}

type errandImpl struct {
	auth.Manager
	errand.Repository
	UserRepo           user.Repository
	NotificationRepo   notification.Repository
	CategoryRepository category.Repository
	error_service.Service
}

func NewErrandUseCase(
	manager auth.Manager,
	repository errand.Repository,
	service error_service.Service,
	userRepo user.Repository,
	notificationRepo notification.Repository,
	categoryRepository category.Repository,
) ErrandUseCase {
	return &errandImpl{
		Manager:            manager,
		Repository:         repository,
		Service:            service,
		UserRepo:           userRepo,
		NotificationRepo:   notificationRepo,
		CategoryRepository: categoryRepository,
	}
}

func (e *errandImpl) CreateErrand(token, userId string, errand *errand.Errand) error {
	adminUserId, resp := e.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	oUser, err := e.UserRepo.GetWithId(userId)
	if err != nil {
		return errors.New(e.Service.HandleMongoDbError("user", err).Message)
	}
	if !oUser.IsOffline() {
		return errors.New("admin can only create errand for offline users")
	}

	if nCategory, err := e.CategoryRepository.Get(errand.Category.Id.Hex()); err != nil {
		return errors.New(e.Service.HandleMongoDbError("category", err).Message)
	} else {
		if errand.DropOffAddress == nil && nCategory.Type == "task" {
			return errors.New("drop-off location is required for tasks")
		}
		errand.Category = nCategory
	}
	errand.CreatedAt = time.Now()
	errand.UpdateForCreation(entity.CreatedByAdmin(*adminUserId))
	errand.Timeline = timeline.NewTimeline(errand.Id.Hex())
	errand.UserId = userId

	err = e.Repository.Update(errand)
	if err != nil {
		return errors.New(e.Service.HandleMongoDbError("nErrand", err).Message)
	}

	return nil
}

func (e *errandImpl) AssignErrand(token, errandId, userId string) error {
	adminUserId, resp := e.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	nErrand, err := e.Repository.Get(errandId)
	if err != nil {
		return errors.New(e.Service.HandleMongoDbError("errand", err).Message)
	}

	if !nErrand.IsAdminErrand(*adminUserId) {
		return errors.New("admin not authorized to assign this errand")
	}

	cUser, err := e.UserRepo.GetWithId(userId)
	if err != nil {
		return errors.New(e.Service.HandleMongoDbError("user", err).Message)
	}

	if !cUser.IsOffline() {
		return errors.New("admin can only assign bid to offline user")
	}
	nBid := bid.NewOfflineBid(errandId, userId, nErrand.Budget)
	if err = e.Repository.AssignErrandToOfflineRunner(*adminUserId, errandId, userId, nBid); err != nil {
		return errors.New(e.Service.HandleMongoDbError("errand", err).Message)
	}

	return nil
}

func (e *errandImpl) GetAllErrands(errandType string) ([]errand.Errand, error) {
	var err error
	var errands []errand.Errand
	if errandType == "open" {
		errands, err = e.Repository.GetAllMarketErrands()
	} else if errandType == "cancelled" {
		errands, err = e.Repository.GetAllCompletedErrands()
	} else if errandType == "active" {
		errands, err = e.Repository.GetAllActiveErrands()
	} else if errandType == "draft" {
		errands, err = e.Repository.GetAllDraftErrands()
	} else {
		errands, err = e.Repository.GetAll()
	}

	if err != nil {
		return nil, errors.New(e.Service.HandleMongoDbError("errands", nil).Message)
	}

	return errands, nil
}
