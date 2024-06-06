package auth

import (
	"DX/src/pkg/response"
	"DX/src/pkg/token_service"
)

type manager struct {
	token_service.Service
	Repository
}

func NewManager(service token_service.Service, repository Repository) Manager {
	return &manager{
		Service:    service,
		Repository: repository,
	}
}

func (m *manager) Get(accessToken string) (*string, *response.BaseResponse) {
	uuid, err := m.Service.GetTokenMetaData(accessToken, true)
	if err != nil {
		return nil, err
	}
	return m.Repository.Get(uuid)
}

func (m *manager) CreateToken(token *Token, userId string) *response.BaseResponse {
	accessToken, refreshToken, err := m.Service.Create(int64(token.AccessTokenExpires), token.RefreshTokenExpires, token.AccessUuid, token.RefreshUuid, userId)
	if err != nil {
		return err
	}

	//Update token values
	token.AccessToken = accessToken
	token.RefreshToken = refreshToken

	//Save tokens to redis
	return m.Repository.CreateToken(token, userId)
}

func (m *manager) DeleteToken(accessToken string) *response.BaseResponse {
	uuid, err := m.Service.GetTokenMetaData(accessToken, true)
	if err != nil {
		return err
	}
	return m.Repository.DeleteToken(uuid)
}

func (m *manager) RefreshToken(*Token, string, string) *response.BaseResponse {
	return nil
}
