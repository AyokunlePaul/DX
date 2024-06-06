package timeline

import (
	"DX/src/domain/entity"
	"time"
)

type Timeline struct {
	Id        entity.DatabaseId `json:"id" bson:"id"`
	ErrandId  string            `json:"errand_id" bson:"errand_id"`
	Updates   []Update          `json:"updates" bson:"updates"`
	CreatedAt time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" bson:"updated_at"`
}

type Update struct {
	Id         entity.DatabaseId `json:"id" bson:"id"`
	Message    string            `json:"message" bson:"message"`
	Type       string            `json:"type" bson:"type"`
	ErrandType Type              `json:"-" bson:"errand_type"`
	Source     string            `json:"source" bson:"source"`
	CreatedAt  time.Time         `json:"created_at" bson:"created_at"`
}

type Type int

const (
	ErrandCreated Type = iota
	BidAccepted
	ErrandStarted
	RunnerUpdate
	SenderRequest
	ErrandCancelled
	ErrandCompleted
)

func NewUpdate(message string, updateType Type, source string) Update {
	return Update{
		Id:         entity.NewDatabaseId(),
		Message:    message,
		Type:       updateType.Id(),
		ErrandType: updateType,
		CreatedAt:  time.Now(),
		Source:     source,
	}
}

func NewTimeline(errandId string) *Timeline {
	cTime := time.Now()
	return &Timeline{
		Id:        entity.NewDatabaseId(),
		CreatedAt: cTime,
		UpdatedAt: cTime,
		Updates: []Update{
			NewUpdate("Errand created", ErrandCreated, entity.Sender.Id()),
		},
		ErrandId: errandId,
	}
}

func (t Type) String() string {
	if t == ErrandCreated {
		return "Errand Created"
	}
	if t == BidAccepted {
		return "Bid Accepted"
	}
	if t == ErrandStarted {
		return "Errand Started"
	}
	if t == RunnerUpdate {
		return "Sender Request"
	}
	if t == SenderRequest {
		return "Runner Update"
	}
	if t == ErrandCancelled {
		return "Errand Cancelled"
	}
	if t == ErrandCompleted {
		return "Errand Completed"
	}
	return ""
}

func (t Type) Id() string {
	if t == ErrandCreated {
		return "errand-created"
	}
	if t == BidAccepted {
		return "bid-accepted"
	}
	if t == ErrandStarted {
		return "errand-started"
	}
	if t == RunnerUpdate {
		return "runner-update"
	}
	if t == SenderRequest {
		return "sender-request"
	}
	if t == ErrandCancelled {
		return "errand-cancelled"
	}
	if t == ErrandCompleted {
		return "errand-completed"
	}
	return ""
}
