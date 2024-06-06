package handler

import (
	"DX/src/domain/entity/user"
	"DX/src/domain/usecase/authentication"
	errands "DX/src/domain/usecase/errand"
	"DX/src/pkg/response"
	"DX/src/utils/logger"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Authentication interface {
	VerifyPhoneNumber(*gin.Context)
	Login(*gin.Context)
	LogOut(*gin.Context)
	CreateAccount(*gin.Context)
	ForgotPassword(*gin.Context)
	RefreshToken(*gin.Context)
	SuspendUser(*gin.Context)
	Profile(*gin.Context)
	MyErrands(*gin.Context)
	MyNotifications(*gin.Context)
	UpdatePassword(*gin.Context)
	GetUser(*gin.Context)
	RateUser(*gin.Context)
}

type authImpl struct {
	authentication.UseCase
	ErrandUseCase errands.UseCase
}

func NewAuthHandler(useCase authentication.UseCase, ErrandUseCase errands.UseCase) Authentication {
	return &authImpl{
		UseCase:       useCase,
		ErrandUseCase: ErrandUseCase,
	}
}

func (i *authImpl) VerifyPhoneNumber(ctx *gin.Context) {
	var payload Payload

	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("invalid request data"))
		return
	}

	phone, ok := payload["phone_number"].(string)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("phone number is required"))
		return
	}

	if resp := i.ValidatePhoneNumber(phone); resp != nil {
		ctx.AbortWithStatusJSON(resp.Code, resp)
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("phone number is valid", nil))
}

func (i *authImpl) CreateAccount(ctx *gin.Context) {
	var payload Payload
	err := ctx.ShouldBind(&payload)
	if err != nil {
		logger.Error("CreateAccount::", err)
		ctx.JSON(http.StatusBadRequest, response.NewBadRequestError("invalid request data"))
		return
	}
	newUser, err := user.Create(payload)
	if err != nil {
		logger.Error("CreateAccount::", err)
		ctx.JSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	resp := newUser.IsValidForInitialCreation()
	if resp != nil {
		logger.Error("CreateAccount::", errors.New(resp.Message))
		ctx.JSON(http.StatusBadRequest, response.NewBadRequestError(resp.Message))
		return
	}

	resp = i.UseCase.CreateUser(newUser)
	if resp != nil {
		logger.Error("CreateAccount::", errors.New(resp.Message))
		ctx.JSON(http.StatusBadRequest, response.NewBadRequestError(resp.Message))
		return
	}

	newUser.UpdateVerification()
	ctx.JSON(http.StatusCreated, response.NewCreatedResponse("user successfully created", newUser))
}

func (i *authImpl) Login(ctx *gin.Context) {
	var payload Payload
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("invalid request data"))
		return
	}

	newUser, err := user.CreateForLogin(payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	resp := i.UseCase.LoginUser(newUser)
	if resp != nil {
		ctx.AbortWithStatusJSON(resp.Code, resp)
		return
	}
	if newUser.IsSuspended {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.NewBadRequestError("user is suspended"))
		return
	}

	newUser.UpdateVerification()
	ctx.JSON(http.StatusOK, response.NewOkResponse("user login successful", newUser))
}

func (i *authImpl) LogOut(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	resp := i.UseCase.LogOut(token)

	if resp != nil {
		ctx.AbortWithStatusJSON(resp.Code, resp)
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("user logged out successfully", nil))
}

func (i *authImpl) SuspendUser(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	resp := i.UseCase.SuspendUser(token)

	if resp != nil {
		ctx.AbortWithStatusJSON(resp.Code, resp)
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("user suspended successfully", nil))
}

func (i *authImpl) Profile(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	nUser, err := i.UseCase.GetUserProfile(token)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("user profile fetched successfully", nUser))
}

func (i *authImpl) GetUser(ctx *gin.Context) {
	userId := ctx.Param("id")
	if nUser, err := i.UseCase.GetUser(userId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	} else {
		ctx.JSON(http.StatusOK, response.NewOkResponse("user data fetched successfully", nUser))
	}
}

func (i *authImpl) UpdatePassword(ctx *gin.Context) {
	var payload Payload
	var password, phone string
	var ok bool

	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}
	if password, ok = payload["password"].(string); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("password is required"))
		return
	}
	if phone, ok = payload["phone_number"].(string); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("phone number is required"))
		return
	}

	err = i.UseCase.UpdateUserPassword(phone, password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, response.NewOkResponse("password successfully updated", nil))
}

func (i *authImpl) MyErrands(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	nErrands, err := i.ErrandUseCase.GetErrandsFor(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, response.NewOkResponse("errands fetched", nErrands))
}

func (i *authImpl) RateUser(ctx *gin.Context) {
	var payload Payload
	var errandId, runnerId string
	var rating float64
	var ok bool

	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}
	if runnerId, ok = payload["runner_id"].(string); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("runner id is required"))
		return
	}
	if errandId, ok = payload["errand_id"].(string); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("errand id is required"))
		return
	}
	if rating, ok = payload["rating"].(float64); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("rating is required"))
		return
	}

	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	if err = i.ErrandUseCase.RateUser(token, runnerId, errandId, int64(rating)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("user rated successfully", nil))
}

func (i *authImpl) MyNotifications(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	notifications, err := i.UseCase.GetNotifications(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("notifications fetched successfully", notifications))
}

func (i *authImpl) ForgotPassword(ctx *gin.Context) {

}

func (i *authImpl) RefreshToken(ctx *gin.Context) {

}
