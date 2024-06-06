package init_data

import (
	"DX/src/domain/entity/category"
	"DX/src/domain/entity/init_data"
	"DX/src/pkg/error_service"
	"errors"
)

type UseCase interface {
	GetInitData() (*init_data.Data, error)
}

type impl struct {
	CategoryRepository category.Repository
	error_service.Service
}

func NewUseCase(CategoryRepository category.Repository) UseCase {
	return &impl{
		CategoryRepository: CategoryRepository,
	}
}

func (i *impl) GetInitData() (*init_data.Data, error) {
	data := &init_data.Data{
		Categories: []category.Category{},
	}
	categories, err := i.CategoryRepository.GetAll()
	if err != nil {
		return nil, errors.New(i.Service.HandleMongoDbError("category", err).Message)
	}

	data.Categories = categories

	return data, nil
}
