package errand

import (
	"DX/src/domain/entity"
	"DX/src/domain/entity/bid"
	"DX/src/domain/entity/category"
	"DX/src/domain/entity/timeline"
	"errors"
	"strings"
	"time"
)

func New(userId string) *Errand {
	currTime := time.Now()
	return &Errand{
		Id:         entity.NewDatabaseId(),
		UserId:     userId,
		State:      Draft,
		Status:     Draft.Id(),
		CreatedAt:  currTime,
		UpdatedAt:  currTime,
		Step:       -1,
		ModifiedBy: []entity.ModifiedBy{},
		Bids:       []bid.Bid{},
		TotalBids:  0,
		Budget:     0,
	}
}

type Errand struct {
	Id                 entity.DatabaseId   `json:"id" bson:"_id"`
	Description        string              `json:"description,omitempty" bson:"description,omitempty"`
	UserId             string              `json:"user_id" bson:"user_id"`
	RunnerId           string              `json:"runner_id" bson:"runner_id"`
	State              State               `json:"-" bson:"state"`
	Category           *category.Category  `json:"category" bson:"category"`
	CreatedBy          *entity.CreatedBy   `json:"created_by" bson:"created_by"`
	Duration           *Duration           `json:"-" bson:"duration,omitempty"`
	Images             []string            `json:"images,omitempty" bson:"images,omitempty"`
	Audio              []string            `json:"audio,omitempty" bson:"audio,omitempty"`
	Step               int                 `json:"step" bson:"step"`
	TotalBids          int                 `json:"total_bids,omitempty" bson:"total_bids"`
	Bids               []bid.Bid           `json:"bids" bson:"bids"`
	Timeline           *timeline.Timeline  `json:"timeline" bson:"timeline"`
	Restriction        string              `json:"restriction,omitempty" bson:"restriction,omitempty"`
	RestrictBy         Restriction         `json:"-" bson:"restrict_by,omitempty"`
	PickupAddress      *Address            `json:"pickup_address,omitempty" bson:"pickup_address,omitempty"`
	DropOffAddress     *Address            `json:"dropoff_address,omitempty" bson:"dropoff_address,omitempty"`
	Budget             int64               `json:"budget" bson:"budget"`
	Amount             int64               `json:"amount" bson:"amount"` // Amount agreed after bidding is accepted
	Status             string              `json:"status" bson:"status"`
	ExpiryDate         time.Time           `json:"expiry_date" bson:"expiry_date,omitempty"`
	CancellationReason string              `json:"cancellation_reason,omitempty" bson:"cancellation_reason"`
	ModifiedBy         []entity.ModifiedBy `json:"-" bson:"modified_by"`
	CreatedAt          time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at" bson:"updated_at"`
}

type MarketErrand struct {
	Id          entity.DatabaseId `json:"id" bson:"id"`
	Bids        int               `json:"bids" bson:"bids"`
	Description string            `json:"description" bson:"description"`
	Location    Address           `json:"location" bson:"location"`
	User        MarketUser        `json:"user" bson:"user"`
}

type MarketUser struct {
	Id             entity.DatabaseId `json:"id" bson:"_id"`
	FirstName      string            `json:"first_name" bson:"first_name"`
	LastName       string            `json:"last_name" bson:"last_name"`
	Rating         float64           `json:"rating" bson:"rating"`
	ProfilePicture string            `json:"profile_picture" bson:"profile_picture"`
}

func (e *Errand) CanBeUpdated() bool {
	return e.State == Draft || e.State == EditMode
}

func (e *Errand) CanBeBiddedFor() bool {
	return e.State == Open || e.State == Pending
}

func ValidateErrandFromMap(data map[string]interface{}, isAdmin bool) (*Errand, error) {
	var duration Duration
	var nErrand *Errand

	if isAdmin {
		nErrand = New("")
	} else {
		nErrand = &Errand{
			ModifiedBy: []entity.ModifiedBy{},
		}
	}

	if dur, ok := data["duration"].(map[string]interface{}); !ok {
		return nil, errors.New("duration is required")
	} else {
		//TODO handle potential errors here
		duration.Period = dur["period"].(string)
		duration.Value = int64(dur["value"].(float64))

		err := duration.IsValid()
		if err != nil {
			return nil, err
		}

		nErrand.Duration = &duration
	}
	if description, ok := data["description"].(string); ok {
		nErrand.Description = description
	}
	if imgFiles, ok := data["images"].([]interface{}); ok {
		if len(imgFiles) > 3 {
			return nil, errors.New("invalid number of images")
		}
		for _, imgFile := range imgFiles {
			nErrand.Images = append(nErrand.Images, imgFile.(string))
		}
	}
	if audFiles, ok := data["audio"].([]interface{}); ok {
		if len(audFiles) > 1 {
			return nil, errors.New("invalid number of audio files")
		}
		for _, audFile := range audFiles {
			nErrand.Audio = append(nErrand.Audio, audFile.(string))
		}
	}
	if rstn, ok := data["restriction"].(string); ok {
		rstnValue := RestrictionType(rstn)
		if rstnValue == -1 {
			return nil, errors.New("invalid restriction type")
		}
		nErrand.RestrictBy = rstnValue
		nErrand.Restriction = rstn
	}
	if categoryId, ok := data["category"].(string); !ok {
		return nil, errors.New("errand category id is required")
	} else {
		catId, _ := entity.StringToErrandId(categoryId)
		nErrand.Category = &category.Category{
			Id: catId,
		}
	}
	if pickupLocation, ok := data["pickup_location"].(map[string]interface{}); !ok {
		return nil, errors.New("pick-up location is required")
	} else {
		//TODO handle potential errors here
		address := &Address{
			Latitude:  pickupLocation["lat"].(float64),
			Longitude: pickupLocation["lng"].(float64),
		}
		nErrand.PickupAddress = address
	}
	if dropoffLocation, ok := data["dropoff_location"].(map[string]interface{}); ok {
		//TODO handle potential errors here
		address := &Address{
			Latitude:  dropoffLocation["lat"].(float64),
			Longitude: dropoffLocation["lng"].(float64),
		}
		nErrand.DropOffAddress = address
	}
	if budget, ok := data["budget"].(interface{}); !ok {
		return nil, errors.New("budget is required")
	} else {
		if budget.(float64) <= 0 {
			return nil, errors.New("invalid budget")
		} else {
			nErrand.Budget = int64(budget.(float64))
		}
	}

	return nErrand, nil
}

func (e *Errand) Update(errand *Errand) {
	if errand.Duration != nil && e.Duration != errand.Duration {
		e.Duration = errand.Duration
	}
	if errand.Images != nil {
		e.Images = errand.Images
	}
	if e.Step != 5 && e.Step != errand.Step {
		e.Step = errand.Step
	}
}

func (e *Errand) UpdateForCreation(createdBy *entity.CreatedBy) {
	cTime := time.Now()
	e.UpdatedAt = cTime
	e.State = Open
	e.Step = 4
	e.Status = e.State.Id()
	e.ExpiryDate = e.Duration.ExpiryDate()
	e.ModifiedBy = append(e.ModifiedBy, entity.ModifiedBy{
		Id:   createdBy.Id,
		Date: cTime,
	})
	e.CreatedBy = createdBy
}

func (e *Errand) Cancel(userId string, reason string) {
	cTime := time.Now()
	if strings.TrimSpace(reason) == "" {
		e.CancellationReason = reason
	}
	e.ModifiedBy = append(e.ModifiedBy, entity.ModifiedBy{
		Id:   userId,
		Date: cTime,
	})
	e.State = Cancelled
	e.Status = Cancelled.Id()
	e.UpdatedAt = cTime
}

func (e *Errand) Complete(userId string, source string) {
	cTime := time.Now()
	var state State

	if source == "sender" {
		state = Completed
	} else {
		state = RunnerCompleted
	}
	e.State = state
	e.Status = state.Id()
	e.ModifiedBy = append(e.ModifiedBy, entity.ModifiedBy{
		Id:   userId,
		Date: cTime,
	})
	e.UpdatedAt = cTime
}

func (e *Errand) AcceptBid(runnerId string, amount float64) {
	cTime := time.Now()
	e.ModifiedBy = append(e.ModifiedBy, entity.ModifiedBy{
		Id:   e.UserId,
		Date: cTime,
	})
	e.RunnerId = runnerId
	e.State = Active
	e.Status = Active.Id()
	e.UpdatedAt = cTime
	e.Amount = int64(amount)
}

func (e *Errand) IsValidBidAndRunner(bidId, runnerId string) (*bid.Bid, error) {
	bidFound := false
	for _, cBid := range e.Bids {
		if cBid.Id.Hex() == bidId {
			bidFound = true
			if cBid.Runner != runnerId {
				return nil, errors.New("invalid runner id for selected bid")
			}
			return &cBid, nil
		}
	}
	if !bidFound {
		return nil, errors.New("bid not valid for errand")
	}

	return nil, nil
}

func (e *Errand) HasAcceptedBid() bool {
	for _, cBid := range e.Bids {
		if cBid.IsAccepted() {
			return true
		}
	}

	return false
}

func (e *Errand) IsAdminErrand(userId string) bool {
	return e.CreatedBy.Admin() && e.CreatedBy.Id == userId
}

func (e *Errand) InProgress() bool {
	return e.State == Active
}

func (e *Errand) IsCompleted() bool {
	return e.State == Active
}

type Duration struct {
	Value  int64
	Period string
}

func (d *Duration) IsValid() error {
	if d.Period != "days" && d.Period != "weeks" && d.Period != "hours" {
		return errors.New("invalid duration period")
	}
	if d.Period == "weeks" {
		if d.Value > 10 {
			return errors.New("duration cannot be longer than 10 weeks")
		}
	}
	return nil
}

func (d *Duration) ExpiryDate() time.Time {
	currentTime := time.Now()

	if d.Period == "hours" {
		currentTime.Add(time.Duration(d.Value) * time.Hour)
	}
	if d.Period == "days" {
		currentTime.Add(time.Duration(d.Value) * 24 * time.Hour)
	}
	if d.Period == "weeks" {
		currentTime.Add(time.Duration(d.Value) * 7 * 24 * time.Hour)
	}

	return currentTime
}

type Restriction int //Qualification, Verification, Insurance

type Address struct {
	Latitude  float64 `json:"lat,omitempty" bson:"latitude"`
	Longitude float64 `json:"lng,omitempty" bson:"longitude"`
}
