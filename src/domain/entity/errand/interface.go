package errand

import (
	"DX/src/domain/entity/bid"
	"DX/src/domain/entity/haggle"
	"DX/src/domain/entity/timeline"
)

type Writer interface {
	Create(*Errand) error
	Update(*Errand) error
	AddBidToErrand(string, string, *bid.Bid) error
	AcceptBid(string, string, string, int64, timeline.Update) error
	RejectBid(string, string) error
	StartErrand(string, string, timeline.Update) error
	ResetErrandBids(string, string) error
	UpdateBidHaggle(string, string, *haggle.Haggle) error
	UpdateTimeline(string, string, timeline.Update) error
	AssignErrandToSender(string, string, string) error
	AssignErrandToRunner(string, string, string) error
	AssignErrandToOfflineRunner(string, string, string, *bid.Bid) error
	RunnerComplete(string, string) error
	SenderComplete(string, string) error
	Delete(string) error
}

type Reader interface {
	Get(string) (*Errand, error)
	GetDraft(string) (*Errand, error)
	GetFor(string) ([]Errand, error)
	GetBidForUser(string, string) error
	GetAllMarketErrands() ([]Errand, error)
	Search(string) ([]string, error)
	GetAll() ([]Errand, error)
	GetAllDraftErrands() ([]Errand, error)
	GetAllCompletedErrands() ([]Errand, error)
	GetAllCancelledErrands() ([]Errand, error)
	GetAllActiveErrands() ([]Errand, error)
	GetAllAbandonedErrands() ([]Errand, error)
}

type Repository interface {
	Reader
	Writer
}
