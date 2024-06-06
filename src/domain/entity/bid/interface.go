package bid

import "DX/src/domain/entity/haggle"

type reader interface {
	GetBid(string) (*Bid, error)
	GetBidsBetween(string, string, string) ([]Bid, error)
	GetAllBids(string) ([]Bid, error)
	GetBidWithRunnerAndErrandId(string, string) (*Bid, error)
}

type writer interface {
	CreateBid(*Bid) error
	UpdateBid(*Bid) error
	UpdateBidHaggle(string, *haggle.Haggle) error
	DeleteBidsFor(string) error
}

type Repository interface {
	reader
	writer
}
