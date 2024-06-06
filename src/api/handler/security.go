package handler

import (
	entSecurity "DX/src/domain/entity/security"
	"DX/src/domain/usecase/security"
	"DX/src/pkg/response"
	"DX/src/utils/logger"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Security interface {
	SaveSecurityQuestion(ctx *gin.Context)
	GetSecurityQuestion(ctx *gin.Context)
	VerifySecurityQuestion(ctx *gin.Context)
}

type securityHandler struct {
	security.UseCase
}

func NewSecurityHandler(useCase security.UseCase) Security {
	return &securityHandler{
		UseCase: useCase,
	}
}

func (s *securityHandler) SaveSecurityQuestion(ctx *gin.Context) {
	var payload Payload
	err := ctx.ShouldBind(&payload)
	if err != nil {
		logger.Error("SaveSecurityQuestion::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("invalid request data"))
		return
	}
	sec, err := entSecurity.Create(payload)
	if err != nil {
		logger.Error("SaveSecurityQuestion::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	err = s.UseCase.SaveSecurity(token, sec)
	if err != nil {
		logger.Error("SaveSecurityQuestion::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("security question saved", nil))
}

func (s *securityHandler) GetSecurityQuestion(ctx *gin.Context) {
	phone := ctx.Query("phone_number")
	if strings.TrimSpace(phone) == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("phone number is required"))
		return
	}

	ques, err := s.UseCase.GetSecurityWithPhone(phone)

	if err != nil {
		logger.Error("GetSecurityQuestion::", errors.New(err.Error()))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("security question fetched", entSecurity.Security{
		Question: ques,
	}))
}

func (s *securityHandler) VerifySecurityQuestion(ctx *gin.Context) {
	var payload Payload

	err := ctx.ShouldBind(&payload)
	if err != nil {
		logger.Error("GetSecurityQuestion::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("invalid request data"))
		return
	}

	phone, answer, err := entSecurity.GetAnswer(payload)
	if err != nil {
		logger.Error("VerifySecurityQuestion::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}
	err = s.UseCase.Validate(phone, answer)

	if err != nil {
		logger.Error("VerifySecurityQuestion::", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("answer is correct", nil))
}
