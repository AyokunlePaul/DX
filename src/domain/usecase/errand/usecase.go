package errand

import (
	"DX/src/domain/entity/bid"
	"DX/src/domain/entity/errand"
	"DX/src/domain/entity/haggle"
)

type UseCase interface {
	GetErrand(string, string) (*errand.Errand, error)
	GetAllErrands(string) ([]errand.Errand, error)
	GetErrandsFor(string) ([]errand.Errand, error)
	UpdateErrand(string, *errand.Errand) error
	CancelErrand(string, string, string) error
	CompleteErrand(string, string, string) error
	CreateErrand(string, string, *errand.Errand) error
	CreateDraftErrand(string) (*errand.Errand, error)
	BidForErrand(string, *bid.Bid, *haggle.Haggle) error
	UpdateErrandBid(string, string, string, *haggle.Haggle) error
	AcceptBid(string, string, string, string, float64) error
	RejectBid(string, string, string) error
	RequestErrandTimelineUpdate(string, string) error
	UpdateErrandTimeline(string, string, string) error
	AcceptContract(string, string, string) error
	RejectContract(string, string, string) error
	RateUser(string, string, string, int64) error
}
