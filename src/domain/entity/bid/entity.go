package bid

import (
	"DX/src/domain/entity"
	"DX/src/domain/entity/haggle"
	"time"
)

type Bid struct {
	Id          entity.DatabaseId `json:"id" bson:"_id"`
	Description string            `json:"description" bson:"description"`
	ErrandId    string            `json:"errand_id" bson:"errand_id"`
	Runner      string            `json:"runner" bson:"runner"`
	BidState    State             `json:"-" bson:"bid_state"`
	State       string            `json:"state" bson:"state"`
	Haggles     []haggle.Haggle   `json:"haggles" bson:"haggles"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at"`
}

func NewBid(errandId string) *Bid {
	currTime := time.Now()
	return &Bid{
		Id:        entity.NewDatabaseId(),
		CreatedAt: currTime,
		UpdatedAt: currTime,
		ErrandId:  errandId,
		Haggles:   []haggle.Haggle{},
	}
}

func NewOfflineBid(errandId, runnerId string, amount int64) *Bid {
	currTime := time.Now()
	return &Bid{
		Id:          entity.NewDatabaseId(),
		Description: "This bid is automatically created by the admin for an offline user",
		State:       Accepted.Id(),
		BidState:    Accepted,
		Haggles: []haggle.Haggle{
			haggle.NewOfflineHaggle(amount),
		},
		CreatedAt: currTime,
		UpdatedAt: currTime,
		ErrandId:  errandId,
		Runner:    runnerId,
	}
}

func FromPayload(data map[string]interface{}, errandId string) (*Bid, error) {
	cTime := time.Now()
	nBid := &Bid{
		Id:        entity.NewDatabaseId(),
		CreatedAt: cTime,
		UpdatedAt: cTime,
		ErrandId:  errandId,
		Haggles:   []haggle.Haggle{},
		State:     Open.Id(),
		BidState:  Open,
	}

	if description, ok := data["description"].(string); ok {
		nBid.Description = description
	}

	return nBid, nil
}

func (b *Bid) Accept() {
	b.UpdatedAt = time.Now()
	b.BidState = Accepted
	b.State = Accepted.Id()
}

func (b *Bid) Reject() {
	b.UpdatedAt = time.Now()
	b.BidState = Rejected
	b.State = Rejected.Id()
}

func (b *Bid) IsAccepted() bool {
	return b.BidState == Accepted
}
