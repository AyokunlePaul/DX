package category

import (
	"DX/src/domain/entity"
	"strings"
	"time"
)

type Category struct {
	Id         entity.DatabaseId   `json:"id" bson:"_id"`
	ImageUrl   string              `json:"image_url" bson:"image_url"`
	Identifier string              `json:"identifier" bson:"identifier"`
	Name       string              `json:"name" bson:"name"`
	Type       string              `json:"type" bson:"type"`
	CreatedBy  string              `json:"created_by" bson:"created_by"`
	ModifiedBy []entity.ModifiedBy `json:"modified_by" bson:"modified_by"`
	CreatedAt  time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at" bson:"updated_at"`
}

func New(name, id, categoryType, iconUrl string) *Category {
	cTime := time.Now()

	return &Category{
		Id:         entity.NewDatabaseId(),
		ModifiedBy: []entity.ModifiedBy{},
		CreatedBy:  id,
		ImageUrl:   iconUrl,
		Type:       categoryType,
		CreatedAt:  cTime,
		UpdatedAt:  cTime,
		Name:       name,
		Identifier: strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "-")),
	}
}
