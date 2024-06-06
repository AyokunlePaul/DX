package admin

import (
	"DX/src/api/handler"
	"DX/src/domain/entity/category"
	"DX/src/domain/entity/file"
	adminUseCase "DX/src/domain/usecase/admin"
	"DX/src/pkg/response"
	"DX/src/utils"
	"DX/src/utils/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type Category interface {
	CreateCategory(*gin.Context)
	UpdateCategory(*gin.Context)
	GetAllCategories(*gin.Context)
	DeleteCategory(*gin.Context)
	DeleteCategories(*gin.Context)
}

type categoryImpl struct {
	adminUseCase.CategoryUseCase
}

func NewAdminCategoryHandler(useCase adminUseCase.CategoryUseCase) Category {
	return &categoryImpl{
		CategoryUseCase: useCase,
	}
}

func (a *categoryImpl) CreateCategory(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	var files []*file.File

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}
	name := ctx.PostForm("name")
	categoryType := ctx.PostForm("type")

	if utils.IsEmpty(name) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("category name is required"))
		return
	}

	if utils.IsEmpty(categoryType) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("category type is required"))
		return
	} else if categoryType != "task" && categoryType != "service" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(fmt.Sprintf("invalid category type: %s", categoryType)))
		return
	}

	if icon, ok := form.File["icon"]; !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("category icon is required"))
		return
	} else {
		if len(icon) > 1 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("icon file cannot be more than one"))
			return
		}
		files = file.NewListRequest(name, icon)
	}

	err = a.CategoryUseCase.CreateNewCategory(token, name, categoryType, files)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, response.NewOkResponse("category created", nil))
}

func (a *categoryImpl) UpdateCategory(ctx *gin.Context) {
	var payload handler.Payload
	var name string
	var ok bool

	err := ctx.ShouldBind(&payload)
	if err != nil {
		logger.Error("UpdateCategory::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("invalid request data"))
		return
	}

	if name, ok = payload["name"].(string); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("category name is required"))
		return
	}
	categoryId := ctx.Param("id")
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	err = a.CategoryUseCase.UpdateCategory(token, categoryId, name)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("category updated", nil))
}

func (a *categoryImpl) GetAllCategories(ctx *gin.Context) {
	query := ctx.Query("limit")
	var err error
	var categories []category.Category
	if query == "" {
		categories, err = a.CategoryUseCase.GetAllCategories(0)
	} else {
		limit, _ := strconv.Atoi(query)
		categories, err = a.CategoryUseCase.GetAllCategories(limit)
	}

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("categories fetched", categories))
}

func (a *categoryImpl) DeleteCategory(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	categoryId := ctx.Param("id")

	err := a.CategoryUseCase.DeleteCategory(token, categoryId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("category deleted", nil))
}

func (a *categoryImpl) DeleteCategories(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}
