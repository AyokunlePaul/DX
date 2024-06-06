package handler

import (
	"DX/src/domain/usecase/init_data"
	"DX/src/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Init interface {
	GetInitData(*gin.Context)
	VerifyHash(*gin.Context)
}

type impl struct {
	init_data.UseCase
}

func NewInitHandler(UseCase init_data.UseCase) Init {
	return &impl{
		UseCase: UseCase,
	}
}

func (i *impl) GetInitData(ctx *gin.Context) {
	data, err := i.UseCase.GetInitData()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("data fetched", data))
}

func (i *impl) VerifyHash(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}
