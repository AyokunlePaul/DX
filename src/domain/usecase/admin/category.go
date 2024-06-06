package admin

import (
	"DX/src/domain/entity/auth"
	"DX/src/domain/entity/category"
	"DX/src/domain/entity/file"
	"DX/src/domain/entity/user"
	"DX/src/pkg/error_service"
	"errors"
	"strings"
)

type CategoryUseCase interface {
	CreateNewCategory(string, string, string, []*file.File) error
	UpdateCategory(string, string, string) error
	GetAllCategories(int) ([]category.Category, error)
	DeleteCategory(string, string) error
}

type categoryImpl struct {
	FileManager        file.Repository
	AuthManager        auth.Manager
	UserManager        user.Repository
	CategoryRepository category.Repository
	error_service.Service
}

func NewCategoryUseCase(FileManager file.Repository, AuthManager auth.Manager, UserManager user.Repository,
	CategoryRepo category.Repository, Service error_service.Service) CategoryUseCase {
	return &categoryImpl{
		FileManager:        FileManager,
		AuthManager:        AuthManager,
		UserManager:        UserManager,
		Service:            Service,
		CategoryRepository: CategoryRepo,
	}
}

func (c *categoryImpl) CreateNewCategory(token, name, categoryType string, icon []*file.File) error {
	adminUserId, resp := c.AuthManager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	iconFile := icon[0]
	err := c.FileManager.UploadCategoryIcon(*adminUserId, iconFile)
	if err != nil {
		return errors.New(c.Service.HandleGoogleStorageError(err).Message)
	}

	nCategory := category.New(name, *adminUserId, categoryType, iconFile.UploadedUrl)

	err = c.CategoryRepository.Create(nCategory)
	if err != nil {
		return errors.New(c.Service.HandleMongoDbError("category", err).Message)
	}

	return nil
}

func (c *categoryImpl) UpdateCategory(token, categoryId, name string) error {
	adminUserId, resp := c.AuthManager.Get(token)
	if resp != nil {
		return errors.New(resp.Message)
	}

	err := c.CategoryRepository.Update(categoryId, name, strings.ToLower(name), *adminUserId)
	if err != nil {
		return errors.New(c.Service.HandleMongoDbError("category", err).Message)
	}

	return nil
}

func (c *categoryImpl) GetAllCategories(limit int) (categories []category.Category, err error) {
	if limit == 0 {
		categories, err = c.CategoryRepository.GetAll()
	} else {
		categories, err = c.CategoryRepository.GetTopCategories(limit)
	}

	if err != nil {
		return nil, errors.New(c.Service.HandleMongoDbError("categories", err).Message)
	}

	return categories, nil
}

func (c *categoryImpl) DeleteCategory(_ string, categoryId string) error {
	err := c.CategoryRepository.Delete(categoryId)
	if err != nil {
		return errors.New(c.Service.HandleMongoDbError("category", err).Message)
	}

	return nil
}
