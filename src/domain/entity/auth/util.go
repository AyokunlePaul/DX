package auth

import (
	"DX/src/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func NewUserToken() *Token {
	errandToken := new(Token)
	errandToken.AccessUuid = entity.NewDefaultId().String()
	errandToken.RefreshUuid = entity.NewDefaultId().String()
	errandToken.AccessTokenExpires = primitive.NewDateTimeFromTime(time.Now().Add(6 * 4 * 24 * 7 * time.Hour)) //6 months
	//errandToken.RefreshTokenExpires = time.Now().Add(48 * time.Hour).Unix() //2 days

	return errandToken
}
