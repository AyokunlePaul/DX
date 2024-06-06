package handler

import (
	"DX/src/domain/usecase/wallet"
	"DX/src/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Wallet interface {
	PaystackWebhook(*gin.Context)
	MakePayment(*gin.Context)
	GetBalance(*gin.Context)
	GetWallet(*gin.Context)
	MakeWithdrawal(*gin.Context)
}

type walletImpl struct {
	wallet.UseCase
}

func NewWalletHandler(useCase wallet.UseCase) Wallet {
	return &walletImpl{
		UseCase: useCase,
	}
}

func (w *walletImpl) MakePayment(ctx *gin.Context) {
	var amount float64
	var paymentType string
	var ok bool

	var payload Payload
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}
	if amount, ok = payload["amount"].(float64); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("amount is required"))
		return
	}
	if paymentType, ok = payload["type"].(string); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("type is required"))
		return
	}
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	if paymentType == "credit" {
		err = w.UseCase.AddMoney(token, int64(amount))
	} else {
		err = w.UseCase.RemoveMoney(token, int64(amount))
	}
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("transaction successful", nil))
}

func (w *walletImpl) GetBalance(ctx *gin.Context) {
	balance, err := w.UseCase.GetBalance()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("balance fetched", balance))
}

func (w *walletImpl) GetWallet(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	nWallet, err := w.GetWalletFor(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("wallet fetched successfully", nWallet))
}

func (w *walletImpl) MakeWithdrawal(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (w *walletImpl) PaystackWebhook(ctx *gin.Context) {

}
