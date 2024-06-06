package admin

import (
	"DX/src/api/handler"
	errandEntity "DX/src/domain/entity/errand"
	"DX/src/domain/usecase/admin"
	"DX/src/domain/usecase/errand"
	"DX/src/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Errand interface {
	CreatErrand(*gin.Context)
	AssignErrandToSender(*gin.Context)
	AssignErrandToRunner(*gin.Context)
	GetAllErrands(*gin.Context)
	GetErrand(*gin.Context)
}

type errandImpl struct {
	admin.ErrandUseCase
	errand.UseCase
}

func NewAdminErrandHandler(adminUseCase admin.ErrandUseCase, useCase errand.UseCase) Errand {
	return &errandImpl{
		ErrandUseCase: adminUseCase,
		UseCase:       useCase,
	}
}

func (e *errandImpl) CreatErrand(ctx *gin.Context) {
	var payload handler.Payload
	var userId string
	var ok bool
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	nErrand, err := errandEntity.ValidateErrandFromMap(payload, true)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}
	if userId, ok = payload["user_id"].(string); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("user id is required"))
		return
	}

	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	err = e.ErrandUseCase.CreateErrand(token, userId, nErrand)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewOkResponse("errand created for user", nil))
}

func (e *errandImpl) AssignErrandToSender(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (e *errandImpl) AssignErrandToRunner(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (e *errandImpl) GetErrand(ctx *gin.Context) {
	errandId := ctx.Param("id")
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	nErrand, err := e.UseCase.GetErrand(token, errandId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("errand fetched", nErrand))
}

func (e *errandImpl) GetAllErrands(ctx *gin.Context) {
	filter := ctx.Query("type")
	errands, err := e.ErrandUseCase.GetAllErrands(filter)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("errands fetched", errands))
}
