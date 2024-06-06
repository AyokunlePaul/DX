package entity

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DatabaseId = primitive.ObjectID
type DefaultID = uuid.UUID
type ErrandJwtToken = jwt.Token
type SearchResult = elastic.SearchResult

type By string

const (
	ByAdmin      By = "admin"
	ByUser       By = "user"
	BySuperAdmin By = "super_admin"
)

type ModifiedBy struct {
	Id   string
	Date time.Time
}

type CreatedBy struct {
	By By
	Id string
}

func CreatedByAdmin(adminId string) *CreatedBy {
	return &CreatedBy{
		Id: adminId,
		By: ByAdmin,
	}
}

func CreatedByUser(userId string) *CreatedBy {
	return &CreatedBy{
		Id: userId,
		By: ByUser,
	}
}

func (c CreatedBy) Admin() bool {
	return c.By == ByAdmin
}

func NewDatabaseId() DatabaseId {
	return primitive.NewObjectIDFromTimestamp(time.Now())
}

func StringToErrandId(idHex string) (DatabaseId, error) {
	return primitive.ObjectIDFromHex(idHex)
}

func NewDefaultId() DefaultID {
	id, _ := uuid.NewRandom()
	return id
}

func StringToDefaultId(uuidString string) (DefaultID, error) {
	return uuid.Parse(uuidString)
}
