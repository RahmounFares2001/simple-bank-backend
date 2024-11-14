package api

import (
	"database/sql"
	"net/http"

	db "github.com/RahmounFares2001/simple-bank-backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

// create entry
type createEntryRequest struct {
	AccountId int64 `json:"account_id" binding:"required,min=1"`
	Amount    int64 `json:"amount" binding:"required,gt=0"`
}

func (server *Server) createEntry(ctx *gin.Context) {
	var req createEntryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// tx: create entry & updata balence
	result, err := server.store.EntryBalenceTx(ctx, db.EntryBalenceTxParams{
		AccountId: req.AccountId,
		Amount:    req.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, result)
}

// get entry
type GetEntryRequest struct {
	EntryId int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getEntry(ctx *gin.Context) {
	var req getAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	entry, err := server.store.GetEntry(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, entry)
}

// get list entry
type GetListEntryRequest struct {
	AccountID int64 `json:"id"`
}
type GetListEntryPagination struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1"`
}

func (server *Server) getListEntries(ctx *gin.Context) {
	var req GetListEntryRequest

	var paginationReq GetListEntryPagination

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	if err := ctx.ShouldBindQuery(&paginationReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	arg := db.ListEntriesParams{
		AccountID: req.AccountID,
		Limit:     paginationReq.PageID,
		Offset:    (paginationReq.PageID - 1) * paginationReq.PageSize,
	}

	entries, err := server.store.ListEntries(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, entries)
}
