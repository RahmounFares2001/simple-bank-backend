package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/RahmounFares2001/simple-bank-backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if !server.validAccount(ctx, req.FromAccountId, req.Currency) {
		return
	}
	if !server.validAccount(ctx, req.ToAccountId, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountId,
		ToAccountID:   req.ToAccountId,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, result)
}

// check currency
func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountId)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}

// get transfers

type getTransferRequest struct {
	TransferID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTransfer(ctx *gin.Context) {
	var req getTransferRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	transfer, err := server.store.GetTransfer(ctx, req.TransferID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, transfer)
}

// get list transfers
type getListTransferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"min=1"`
}

type transferPaginationRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1"`
}

func (server *Server) getListTransfers(ctx *gin.Context) {
	var jsonReq getListTransferRequest
	var queryReq transferPaginationRequest

	// Bind JSON body parameters
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Bind query parameters
	if err := ctx.ShouldBindQuery(&queryReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Custom validation: Ensure at least one account ID is provided and valid
	if jsonReq.FromAccountID < 1 && jsonReq.ToAccountID < 1 {
		err := fmt.Errorf("from_account_id or to_account_id must be provided and greater than 0")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Prepare arguments for querying the database
	arg := db.ListTransfersParams{
		FromAccountID: jsonReq.FromAccountID,
		ToAccountID:   jsonReq.ToAccountID,
		Limit:         queryReq.PageSize,
		Offset:        (queryReq.PageID - 1) * queryReq.PageSize,
	}

	// Query the database
	transfersList, err := server.store.ListTransfers(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Successful response
	ctx.JSON(http.StatusOK, transfersList)
}
