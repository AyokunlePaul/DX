package auth

import "go.mongodb.org/mongo-driver/bson/primitive"

type Token struct {
	AccessToken         string             `json:"access_token,omitempty" bson:"access_token"`
	RefreshToken        string             `json:"-,omitempty" bson:"refresh_token,omitempty"`
	AccessUuid          string             `json:"access_uuid,omitempty" bson:"access_uuid"`
	RefreshUuid         string             `json:"refresh_uuid,omitempty" bson:"refresh_uuid,omitempty"`
	AccessTokenExpires  primitive.DateTime `json:"access_token_expires,omitempty" bson:"access_token_expires"`
	RefreshTokenExpires int64              `json:"refresh_token_expires,omitempty" bson:"refresh_token_expires,omitempty"`
	UserId              string             `json:"user_id" bson:"user_id"`
}
