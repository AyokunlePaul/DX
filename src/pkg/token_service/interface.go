package token_service

import (
	"DX/src/domain/entity"
	"DX/src/pkg/response"
)

type Service interface {
	Create(int64, int64, string, string, string) (string, string, *response.BaseResponse)
	VerifyAndExtract(string, bool) (*entity.ErrandJwtToken, *response.BaseResponse)
	CheckTokenValidity(string, bool) *response.BaseResponse
	GetTokenMetaData(string, bool) (string, *response.BaseResponse)
}
