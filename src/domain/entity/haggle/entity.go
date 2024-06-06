package haggle

import (
	"DX/src/domain/entity"
	"errors"
	"time"
)

type Haggle struct {
	Id           entity.DatabaseId `json:"id" bson:"_id"`
	Source       string            `json:"source" bson:"source"`
	HaggleSource entity.Source     `json:"-" bson:"haggle_source"`
	CreatedAt    time.Time         `json:"created_at" bson:"created_at"`
	Amount       int64             `json:"amount" bson:"amount"`
	Description  string            `json:"description,omitempty" bson:"description"`
}

func (h *Haggle) FromSender() bool {
	if h.HaggleSource == entity.Sender {
		return true
	}
	return false
}

func NewOfflineHaggle(amount int64) Haggle {
	return Haggle{
		Id:           entity.NewDatabaseId(),
		Source:       entity.Admin.Id(),
		HaggleSource: entity.Admin,
		Description:  "This haggle is automatically created by the admin for an offline user",
		CreatedAt:    time.Now(),
		Amount:       amount,
	}
}

func FromPayload(data map[string]interface{}) (*Haggle, error) {
	nHaggle := &Haggle{
		Id:        entity.NewDatabaseId(),
		CreatedAt: time.Now(),
	}
	if amount, ok := data["amount"].(float64); !ok {
		return nil, errors.New("bid amount is required")
	} else {
		nHaggle.Amount = int64(amount)
	}
	if source, ok := data["source"].(string); !ok {
		return nil, errors.New("bid source is required")
	} else {
		hSource, err := entity.GetSource(source)
		if err != nil {
			return nil, err
		}
		nHaggle.HaggleSource = hSource
		nHaggle.Source = source
	}
	if description, ok := data["description"].(string); ok {
		nHaggle.Description = description
	}
	return nHaggle, nil
}
