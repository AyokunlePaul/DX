package admin

import (
	"DX/src/api/handler"
	"DX/src/domain/entity/user"
	adminUseCase "DX/src/domain/usecase/admin"
	"DX/src/pkg/response"
	"DX/src/utils/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type User interface {
	GetAllUsers(*gin.Context)
	GetUser(*gin.Context)
	UpdateUser(*gin.Context)
	SuspendUser(*gin.Context)
	SuspendUsers(*gin.Context)
	RestoreUser(*gin.Context)
	DeleteUser(*gin.Context)
	DeleteUsers(*gin.Context)
	CreateUser(*gin.Context)
	AssignUserErrand(*gin.Context)
}

type impl struct {
	adminUseCase.UserUseCase
}

func NewAdminUserHandler(UseCase adminUseCase.UserUseCase) User {
	return &impl{
		UserUseCase: UseCase,
	}
}

func (i *impl) CreateUser(ctx *gin.Context) {
	var payload handler.Payload
	err := ctx.ShouldBind(&payload)
	if err != nil {
		logger.Error("CreateUser::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("invalid request data"))
		return
	}

	nUser, err := user.CreateForAdmin(payload)
	if err != nil {
		logger.Error("CreateUser::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	err = i.UserUseCase.CreateNewUser(token, nUser)
	if err != nil {
		logger.Error("CreateUser::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewOkResponse("user created", nil))
}

func (i *impl) SuspendUser(ctx *gin.Context) {
	userId := ctx.Param("id")

	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	err := i.UserUseCase.SuspendUser(token, userId)
	if err != nil {
		logger.Error("SuspendUser::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("user successfully updated", nil))
}

func (i *impl) SuspendUsers(ctx *gin.Context) {
	var payload handler.Payload
	var users []string
	err := ctx.ShouldBind(&payload)
	if err != nil {
		logger.Error("SuspendUsers::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("invalid request data"))
		return
	}

	if tUsers, ok := payload["users"].([]interface{}); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("users list is required"))
		return
	} else {
		for _, cUser := range tUsers {
			users = append(users, cUser.(string))
		}
	}

	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	err = i.UserUseCase.SuspendUsers(token, users)
	if err != nil {
		logger.Error("SuspendUsers::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("users suspended successfully", nil))
}

func (i *impl) RestoreUser(ctx *gin.Context) {
	userId := ctx.Param("id")

	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	err := i.UserUseCase.RestoreUser(token, userId)
	if err != nil {
		logger.Error("RestoreUser::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("user successfully updated", nil))
}

func (i *impl) UpdateUser(ctx *gin.Context) {
	var payload handler.Payload
	err := ctx.ShouldBind(&payload)
	if err != nil {
		logger.Error("UpdateUser::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("invalid request data"))
		return
	}

	userId := ctx.Param("id")
	nUser, err := user.UpdateForAdmin(payload)
	if err != nil {
		logger.Error("UpdateUser::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	err = i.UserUseCase.UpdateUser(token, userId, nUser)
	if err != nil {
		logger.Error("UpdateUser::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("user successfully updated", nil))
}

func (i *impl) GetUser(ctx *gin.Context) {
	userId := ctx.Param("id")

	nUser, err := i.UserUseCase.GetUser(userId)
	if err != nil {
		logger.Error("GetUser::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("user fetched successfully", nUser))
}

func (i *impl) GetAllUsers(ctx *gin.Context) {
	filter := ctx.Query("type")
	users, err := i.UserUseCase.GetAllUsers(filter)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("users fetched", users))
}

func (i *impl) DeleteUser(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	userId := ctx.Param("id")

	err := i.UserUseCase.DeleteUser(token, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("user deleted successfully", nil))
}

func (i *impl) DeleteUsers(ctx *gin.Context) {
	var payload handler.Payload
	var users []string
	err := ctx.ShouldBind(&payload)
	if err != nil {
		logger.Error("DeleteUsers::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("invalid request data"))
		return
	}

	if tUsers, ok := payload["users"].([]interface{}); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("users list is required"))
		return
	} else {
		for _, cUser := range tUsers {
			users = append(users, cUser.(string))
		}
	}

	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	err = i.UserUseCase.DeleteUsers(token, users)
	if err != nil {
		logger.Error("DeleteUsers::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("users deleted successfully", nil))
}

func (i *impl) AssignUserErrand(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}
