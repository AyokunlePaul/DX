package auth

import (
	"DX/src/pkg/response"
)

type reader interface {
	Get(string) (*string, *response.BaseResponse)
}

type writer interface {
	CreateToken(*Token, string) *response.BaseResponse
	RefreshToken(*Token, string, string) *response.BaseResponse
	DeleteToken(string) *response.BaseResponse
}

type Repository interface {
	reader
	writer
}

type Manager interface {
	Repository
}
