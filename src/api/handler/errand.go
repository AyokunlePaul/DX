package handler

import (
	"DX/src/domain/entity/bid"
	errandEntity "DX/src/domain/entity/errand"
	fileUtil "DX/src/domain/entity/file"
	"DX/src/domain/entity/haggle"
	errandUseCase "DX/src/domain/usecase/errand"
	"DX/src/domain/usecase/file"
	"DX/src/pkg/response"
	"DX/src/utils"
	"DX/src/utils/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Errand interface {
	GetDraftErrand(*gin.Context)
	UploadErrandFiles(*gin.Context)
	CreateErrand(*gin.Context)
	UpdateErrand(*gin.Context)
	CancelErrand(*gin.Context)
	CompleteErrand(*gin.Context)
	GetErrand(*gin.Context)
	FetchAllErrands(*gin.Context)
	BidForErrand(*gin.Context)
	UpdateBidForErrand(*gin.Context)
	RespondToBid(*gin.Context)
	RequestForUpdate(*gin.Context)
	PostUpdate(*gin.Context)
	StartErrand(*gin.Context)
	RejectErrandContract(*gin.Context)
}

type errand struct {
	errandUseCase.UseCase
	FileUseCase file.UseCase
}

func NewErrandHandler(useCase errandUseCase.UseCase, fileUseCase file.UseCase) Errand {
	return &errand{
		UseCase:     useCase,
		FileUseCase: fileUseCase,
	}
}

func (e *errand) GetDraftErrand(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	nErrand, err := e.UseCase.CreateDraftErrand(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewOkResponse("success", nErrand))
}

func (e *errand) UploadErrandFiles(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	reqType := ctx.PostForm("type")
	errandId := ctx.PostForm("errand_id")
	files := form.File["files"]

	logger.Info(fmt.Sprintf("length of uploaded files: %v", len(files)))

	if len(files) > 3 || len(files) < 1 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("invalid request"))
		return
	}
	if strings.TrimSpace(reqType) == "" || !utils.Contains(fileUtil.ReqTypes, reqType) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("request type is invalid"))
		return
	}
	if strings.TrimSpace(errandId) == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("errand id is mandatory"))
		return
	}

	listReq := fileUtil.NewListRequest(reqType, files)

	resp, err := e.FileUseCase.UploadFiles(token, errandId, listReq)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewOkResponse("files uploaded successfully", resp))
}

func (e *errand) UpdateErrand(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	var payload Payload
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	nErrand, err := errandEntity.ValidateErrandFromMap(payload, false)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	err = e.UseCase.UpdateErrand(token, nErrand)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("errand created", nErrand))
}

func (e *errand) CancelErrand(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	errandId := ctx.Param("id")

	var payload Payload
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	reason := payload["reason"].(string)
	err = e.UseCase.CancelErrand(token, errandId, reason)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("errand cancelled successfully", nil))
}

func (e *errand) CompleteErrand(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	errandId := ctx.Param("id")
	var source string
	var ok bool

	var payload Payload
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}
	if source, ok = payload["source"].(string); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("source is required"))
		return
	} else {
		if source != "sender" && source != "runner" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("invalid source type"))
			return
		}
	}

	err = e.UseCase.CompleteErrand(token, errandId, source)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("errand completed successfully", nil))
}

func (e *errand) CreateErrand(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	var payload Payload
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	nErrand, err := errandEntity.ValidateErrandFromMap(payload, false)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}
	errandId := ctx.Param("id")
	if strings.TrimSpace(errandId) == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("errand id is required"))
		return
	}

	err = e.UseCase.CreateErrand(token, errandId, nErrand)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("errand created", nErrand))
}

func (e *errand) BidForErrand(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	var payload Payload
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}
	errandId := ctx.Param("id")

	nBid, err := bid.FromPayload(payload, errandId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	nHaggle, err := haggle.FromPayload(payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	err = e.UseCase.BidForErrand(token, nBid, nHaggle)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("bid placed", nil))
}

func (e *errand) UpdateBidForErrand(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	var payload Payload
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}
	errandId := ctx.Param("id")
	bidId := ctx.Param("bid_id")

	nHaggle, err := haggle.FromPayload(payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	err = e.UseCase.UpdateErrandBid(token, errandId, bidId, nHaggle)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("bid placed", nil))
}

func (e *errand) StartErrand(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	var payload Payload
	var errandId, bidId string
	var ok bool
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}
	if errandId, ok = payload["errand_id"].(string); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("errand id is required"))
		return
	}
	if bidId, ok = payload["bid_id"].(string); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("bid id is required"))
		return
	}

	err = e.UseCase.AcceptContract(token, errandId, bidId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("errand started successfully", nil))
}

func (e *errand) RejectErrandContract(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	errandId := ctx.Param("id")
	bidId := ctx.Param("bid_id")

	err := e.UseCase.RejectContract(token, errandId, bidId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("errand contract rejected successfully", nil))
}

func (e *errand) RespondToBid(ctx *gin.Context) {
	var amount float64
	var ok bool
	var responseType, runnerId string
	var payload Payload
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}
	if amount, ok = payload["amount"].(float64); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("amount is required"))
		return
	} else {
		if amount < 0.0 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("amount is invalid"))
			return
		}
	}
	if responseType, ok = payload["response"].(string); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("response type is required"))
		return
	} else {
		if responseType != "accept" && responseType != "reject" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("invalid response type specified"))
			return
		}
	}
	if runnerId, ok = payload["runner_id"].(string); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("runner id is required"))
		return
	}

	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	errandId := ctx.Param("id")
	bidId := ctx.Param("bid_id")

	if responseType == "accept" {
		err = e.UseCase.AcceptBid(token, errandId, bidId, runnerId, amount)
	} else if responseType == "reject" {
		err = e.UseCase.RejectBid(token, errandId, bidId)
	}
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	if responseType == "accept" {
		ctx.JSON(http.StatusOK, response.NewOkResponse("bid accepted successfully", nil))
	} else {
		ctx.JSON(http.StatusOK, response.NewOkResponse("bid rejected successfully", nil))
	}
}

func (e *errand) RequestForUpdate(ctx *gin.Context) {
	errandId := ctx.Param("id")

	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]
	err := e.UseCase.RequestErrandTimelineUpdate(token, errandId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("update requested", nil))
}

func (e *errand) PostUpdate(ctx *gin.Context) {
	var payload Payload
	var errandId, message string
	var ok bool

	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}
	if errandId, ok = payload["errand_id"].(string); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("errand id is required"))
		return
	}
	if message, ok = payload["message"].(string); !ok {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError("update message is required"))
		return
	}
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	err = e.UseCase.UpdateErrandTimeline(token, errandId, message)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("update posted successfully", nil))
}

func (e *errand) GetErrand(ctx *gin.Context) {
	errandId := ctx.Param("id")
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	nErrand, err := e.UseCase.GetErrand(token, errandId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("errand details fetched", nErrand))
}

func (e *errand) FetchAllErrands(ctx *gin.Context) {
	token := strings.Split(ctx.GetHeader("Authorization"), " ")[1]

	errands, err := e.GetAllErrands(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewBadRequestError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewOkResponse("errands fetched successfully", errands))
}
