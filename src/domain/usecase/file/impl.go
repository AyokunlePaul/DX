package file

import (
	"DX/src/domain/entity/auth"
	"DX/src/domain/entity/errand"
	"DX/src/domain/entity/file"
	"DX/src/pkg/error_service"
	"errors"
)

type impl struct {
	auth.Manager
	error_service.Service
	file.Repository
	ErrandRepo errand.Repository
}

func NewUseCase(tokenManager auth.Manager, fileRepo file.Repository, errandRepo errand.Repository, service error_service.Service) UseCase {
	return &impl{
		Manager:    tokenManager,
		Repository: fileRepo,
		ErrandRepo: errandRepo,
		Service:    service,
	}
}

func (i *impl) UploadFile(token string, errandId string, payFile *file.File) (string, error) {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return "", errors.New(resp.Message)
	}

	err := i.Create(*userId, errandId, payFile)
	return "", errors.New(i.HandleGoogleStorageError(err).Message)
}

func (i *impl) UploadFiles(token string, errandId string, files []*file.File) ([]string, error) {
	userId, resp := i.Manager.Get(token)
	if resp != nil {
		return nil, errors.New(resp.Message)
	}

	_, err := i.ErrandRepo.Get(errandId)
	if err != nil {
		return nil, errors.New(i.Service.HandleMongoDbError("errand", err).Message)
	}

	err = i.CreateList(*userId, errandId, files)
	if err != nil {
		return nil, errors.New(i.HandleGoogleStorageError(err).Message)
	}

	var urls []string
	for _, nFile := range files {
		urls = append(urls, nFile.UploadedUrl)
	}
	return urls, nil
}
