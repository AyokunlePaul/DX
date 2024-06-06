package errand

import (
	"DX/src/domain/entity"
	"DX/src/domain/entity/auth"
	"DX/src/domain/entity/bid"
	"DX/src/domain/entity/category"
	"DX/src/domain/entity/errand"
	"DX/src/domain/entity/haggle"
	"DX/src/domain/entity/notification"
	"DX/src/domain/entity/timeline"
	"DX/src/domain/entity/user"
	"DX/src/domain/entity/wallet"
	"DX/src/pkg/error_service"
	"DX/src/utils/logger"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type impl struct {
	auth.Manager
	errand.Repository
	UserRepo user.Repository
	error_service.Service
	NotificationRepo   notification.Repository
	CategoryRepository category.Repository
	ErrandRepo         errand.Repository
	WalletRepo         wallet.Repository
}

func NewUseCase(
	manager auth.Manager,
	repository errand.Repository,
	userRepo user.Repository,
	service error_service.Service,
	notificationRepo notification.Repository,
	categoryRepository category.Repository,
	errandRepo errand.Repository,
	walletRepo wallet.Repository,
) UseCase {
	return &impl{
		Manager:            manager,
		Repository:         repository,
		UserRepo:           userRepo,
		Service:            service,
		NotificationRepo:   notificationRepo,
		CategoryRepository: categoryRepository,
		ErrandRepo:         errandRepo,
		WalletRepo:         walletRepo,
	}
}

func (i *impl) CreateDraftErrand(token string) (*errand.Errand, error) {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return nil, errors.New(resp.Message)
	}

	draft, err := i.Repository.GetDraft(*userId)
	if err != nil {
		errMessage := i.Service.HandleMongoDbError("errand", err).Message
		return nil, errors.New(errMessage)
	}

	if draft == nil {
		draft = errand.New(*userId)
		err = i.Repository.Create(draft)
		if err != nil {
			return nil, errors.New(i.Service.HandleMongoDbError("errand", err).Message)
		}
	}

	return draft, nil
}

func (i *impl) UpdateErrand(token string, errand *errand.Errand) error {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	oErrand, err := i.Repository.Get(errand.Id.Hex())
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}

	if *userId != oErrand.UserId {
		return errors.New("user not authorized to update errand")
	}

	oErrand.Update(errand)

	err = i.Repository.Update(oErrand)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}

	return nil
}

func (i *impl) CreateErrand(token string, errandId string, nErrand *errand.Errand) error {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	oErrand, err := i.Repository.Get(errandId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}

	if *userId != oErrand.UserId {
		return errors.New("user not authorized to create this errand")
	}

	if !oErrand.CanBeUpdated() {
		return errors.New("errand can't be updated")
	}
	if balance, err := i.WalletRepo.GetBalance(*userId); err != nil {
		return errors.New(i.Service.HandleMongoDbError("wallet", err).Message)
	} else {
		if nErrand.Budget > balance {
			return errors.New("insufficient funds. kindly top up your wallet")
		}
	}

	if nCategory, err := i.CategoryRepository.Get(nErrand.Category.Id.Hex()); err != nil {
		return errors.New(i.Service.HandleMongoDbError("category", err).Message)
	} else {
		if nErrand.DropOffAddress == nil && nCategory.Type == "task" {
			return errors.New("drop-off location is required for tasks")
		}
		nErrand.Category = nCategory
	}
	nErrand.UpdateForCreation(entity.CreatedByUser(*userId))
	nErrand.UserId = oErrand.UserId
	nErrand.Id = oErrand.Id
	nErrand.CreatedAt = oErrand.CreatedAt
	nErrand.Timeline = timeline.NewTimeline(errandId)

	debitTxn := wallet.NewDebitTransaction(*userId, "Errand creation", nErrand.Id.Hex(), nErrand.Budget)
	if err = i.WalletRepo.CreateTransaction(debitTxn); err != nil {
		return errors.New(i.Service.HandleMongoDbError("wallet", err).Message)
	}

	err = i.Repository.Update(nErrand)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("nErrand", err).Message)
	}

	return nil
}

func (i *impl) GetErrand(token string, errandId string) (*errand.Errand, error) {
	_, resp := i.Manager.Get(token)
	if resp != nil {
		return nil, errors.New(resp.Message)
	}

	nErrand, err := i.Repository.Get(errandId)
	if err != nil {
		return nil, err
	}

	return nErrand, nil
}

func (i *impl) AcceptBid(token, errandId, bidId, runnerId string, amount float64) error {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	nErrand, err := i.Repository.Get(errandId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}
	// Check if user can accept bid
	if nErrand.UserId != *userId {
		return errors.New("user not authorized to accept bid")
	}
	if _, err = nErrand.IsValidBidAndRunner(bidId, runnerId); err != nil {
		return err
	}
	if nErrand.HasAcceptedBid() {
		return errors.New("user already accepted a bid for this errand")
	}

	timelineMessage := "Bid accepted"
	update := timeline.NewUpdate(timelineMessage, timeline.BidAccepted, entity.Sender.Id())
	err = i.Repository.AcceptBid(errandId, bidId, *userId, int64(amount), update)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}

	bidNotification := notification.NewBidAcceptedNotification(*userId, bidId)
	err = i.NotificationRepo.SendNotification(bidNotification)
	if err != nil {
		logger.Error("Failed to send notifications", err)
	}

	return nil
}

func (i *impl) RejectBid(token, errandId, bidId string) error {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	nErrand, err := i.Repository.Get(errandId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}
	// Check if user can accept bid
	if nErrand.UserId != *userId {
		return errors.New("user not authorized to reject bid")
	}

	err = i.Repository.RejectBid(errandId, bidId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("bid", err).Message)
	}

	return nil
}

func (i *impl) CancelErrand(token, errandId, reason string) error {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	oErrand, err := i.Repository.Get(errandId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}

	if oErrand.UserId != *userId {
		return errors.New("user not authorized to cancel errand")
	}

	oErrand.Cancel(*userId, reason)
	err = i.Repository.Update(oErrand)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}

	return nil
}

func (i *impl) CompleteErrand(token, errandId, source string) error {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	oErrand, err := i.Repository.Get(errandId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}

	if source == "sender" {
		if oErrand.UserId != *userId {
			return errors.New("user not authorized to complete errand")
		}
		if err = i.Repository.SenderComplete(errandId, *userId); err != nil {
			return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
		} else {
			if err = i.UserRepo.CompleteErrand(oErrand.RunnerId); err != nil {
				return errors.New(i.Service.HandleMongoDbError("user", err).Message)
			}
		}
		txn := wallet.NewCreditTransaction(oErrand.RunnerId, "Completed errand", oErrand.Id.Hex(), oErrand.Amount)
		if err = i.WalletRepo.CreateTransaction(txn); err != nil {
			return errors.New(i.Service.HandleMongoDbError("wallet", err).Message)
		}
		go i.sendNotification(notification.NewSenderErrandCompletedNotification(oErrand.RunnerId, errandId))
	} else {
		if oErrand.RunnerId != *userId {
			return errors.New("user not authorized to complete errand")
		}
		if err = i.Repository.RunnerComplete(errandId, *userId); err != nil {
			return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
		}
		go i.sendNotification(notification.NewRunnerErrandCompletedNotification(oErrand.UserId, errandId))
	}

	//TODO handle notifications here
	return nil
}

func (i *impl) AcceptContract(token, errandId, bidId string) error {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	oErrand, err := i.Repository.Get(errandId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}

	if !oErrand.HasAcceptedBid() {
		return errors.New("can't accept contract for errand with no accepted bid")
	}
	if _, err = oErrand.IsValidBidAndRunner(bidId, *userId); err != nil {
		return err
	}

	update := timeline.NewUpdate("Errand contract accepted", timeline.ErrandStarted, entity.Runner.Id())
	err = i.Repository.StartErrand(errandId, *userId, update)
	if err != nil {
		logger.Error("unable to start errand", err)
		return errors.New("unable to start errand")
	}

	go i.sendNotification(notification.NewErrandStartedNotification(oErrand.UserId, errandId))

	return nil
}

func (i *impl) RejectContract(token string, errandId string, bidId string) error {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	oErrand, err := i.Repository.Get(errandId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}

	if !oErrand.HasAcceptedBid() {
		return errors.New("can't reject contract for errand with no accepted bid")
	}
	if _, err = oErrand.IsValidBidAndRunner(bidId, *userId); err != nil {
		return err
	}

	err = i.Repository.ResetErrandBids(errandId, *userId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("bid", err).Message)
	}
	go i.sendNotification(notification.NewBidProposalRejectedNotification(oErrand.UserId, errandId))

	return nil
}

func (i *impl) GetAllErrands(token string) ([]errand.Errand, error) {
	// TODO Use the id to optimize marketplace response sent to user
	_, resp := i.Manager.Get(token)
	if resp != nil {
		return nil, errors.New(resp.Message)
	}

	errands, err := i.Repository.GetAllMarketErrands()
	if err != nil {
		return nil, err
	}

	return errands, nil
}

func (i *impl) GetErrandsFor(token string) ([]errand.Errand, error) {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return nil, errors.New(resp.Message)
	}

	errands, err := i.Repository.GetFor(*userId)
	if err != nil {
		return nil, err
	}

	return errands, nil
}

func (i *impl) BidForErrand(token string, bid *bid.Bid, haggle *haggle.Haggle) error {
	// Get user id
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	// Get errand
	nErrand, err := i.Repository.Get(bid.ErrandId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}
	// Check if this is the errand creator
	if nErrand.UserId == *userId {
		return errors.New("sender not allowed to bid for own errand")
	}
	if !nErrand.CanBeBiddedFor() {
		return errors.New("errand no longer available for bidding")
	}

	// TODO Check if runner is qualified to bid for the errand
	// Check if runner has an existing bid for this errand
	err = i.Repository.GetBidForUser(bid.ErrandId, *userId)
	// The query should return mongo.ErrNoDocuments
	if err != mongo.ErrNoDocuments {
		return errors.New("user already has a active bid for errand")
	}

	// Create bid for errand
	bid.Haggles = append(bid.Haggles, *haggle)
	bid.Runner = *userId
	err = i.Repository.AddBidToErrand(bid.ErrandId, *userId, bid)
	if err != nil {
		return errors.New(i.HandleMongoDbError("bid", err).Message)
	}

	err = i.NotificationRepo.SendNotification(notification.NewBidNotification(nErrand.UserId, nErrand.Id.Hex()))
	if err != nil {
		logger.Error("Failed to send notifications", err)
	}

	return nil
}

func (i *impl) UpdateErrandBid(token string, errandId string, bidId string, haggle *haggle.Haggle) error {
	var runnerId string
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	nErrand, err := i.Repository.Get(errandId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}

	if haggle.FromSender() {
		if nErrand.UserId != *userId {
			return errors.New("sender not authorized to update this bid")
		}
	} else {
		if cBid, err := nErrand.IsValidBidAndRunner(bidId, *userId); err != nil {
			return err
		} else {
			runnerId = cBid.Runner
		}
	}

	err = i.Repository.UpdateBidHaggle(errandId, bidId, haggle)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("haggle", err).Message)
	}

	if haggle.FromSender() {
		err = i.NotificationRepo.SendNotification(notification.NewHaggleNotification(runnerId, bidId))
	} else {
		err = i.NotificationRepo.SendNotification(notification.NewHaggleNotification(nErrand.UserId, bidId))
	}
	if err != nil {
		logger.Error("Failed to send notifications", err)
	}

	return nil
}

func (i *impl) RequestErrandTimelineUpdate(token string, errandId string) error {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	nErrand, err := i.Repository.Get(errandId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}

	if nErrand.UserId != *userId {
		return errors.New("user not authorized to request update for this errand")
	}
	if !nErrand.InProgress() {
		return errors.New("errand not in progress")
	}

	timelineMessage := "Update request"
	err = i.Repository.UpdateTimeline(errandId, *userId, timeline.NewUpdate(timelineMessage, timeline.SenderRequest, entity.Sender.Id()))
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("timeline", err).Message)
	}

	// TODO update for other notification types
	title := "New update request"
	message := "An update has been requested for the errand you have in progress"
	err = i.NotificationRepo.SendNotification(notification.NewErrandUpdateRequestNotification(nErrand.RunnerId, errandId, title, message))
	if err != nil {
		logger.Error("Failed to send notifications", err)
	}

	return nil
}

func (i *impl) UpdateErrandTimeline(token, errandId, message string) error {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	nErrand, err := i.Repository.Get(errandId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}

	if nErrand.RunnerId != *userId {
		return errors.New("runner not authorized to update this errand timeline")
	}

	err = i.Repository.UpdateTimeline(errandId, *userId, timeline.NewUpdate(message, timeline.RunnerUpdate, entity.Runner.Id()))
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("timeline", err).Message)
	}

	title := "Errand timeline update"
	notificationMessage := "Errand runner has provided a new update for your errand"
	err = i.NotificationRepo.SendNotification(notification.NewErrandUpdateRequestNotification(nErrand.UserId, errandId, title, notificationMessage))
	if err != nil {
		logger.Error("Failed to send notifications", err)
	}

	return nil
}

func (i *impl) RateUser(token, runnerId, errandId string, rating int64) error {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	oErrand, err := i.ErrandRepo.Get(errandId)
	if err != nil {
		return errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}

	if oErrand.UserId != *userId {
		return errors.New("user not authorized to update errand")
	}
	if !oErrand.IsCompleted() {
		return errors.New("can only rate user for a completed errand")
	}
	if oErrand.RunnerId != runnerId {
		return errors.New("invalid runner id")
	}

	if err = i.UserRepo.RateUser(runnerId, rating); err != nil {
		return errors.New(i.Service.HandleMongoDbError("user", err).Message)
	}

	return nil
}

func (i *impl) sendNotification(notification notification.Notification) {
	err := i.NotificationRepo.SendNotification(notification)
	if err != nil {
		logger.Error("Failed to send notifications", err)
	}
}
