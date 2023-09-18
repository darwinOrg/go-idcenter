package web

import (
	"com.startrek/go-idcenter/sequence"
	dgctx "github.com/darwinOrg/go-common/context"
	dgerr "github.com/darwinOrg/go-common/enums/error"
	"github.com/darwinOrg/go-common/result"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NextIdReq struct {
	SeqName string `json:"seqName" form:"seqName"  binding:"required"`
}

type NextIdsReq struct {
	SeqName string `json:"seqName" form:"seqName"  binding:"required"`
	Count   int    `json:"count" form:"count"  binding:"required"`
}

type NextIdVO struct {
	NextId int64 `json:"nextId" form:"nextId"`
}

type NextIdsVO struct {
	NextIds []int64 `json:"nextIds" form:"nextIds"`
}

func NextId(gc *gin.Context, dc *dgctx.DgContext, nextIdReq *NextIdReq) *result.Result[*NextIdVO] {
	internalCall := isInternalCall(gc)
	if !internalCall {
		gc.AbortWithStatus(http.StatusNotFound)
		return result.FailByError[*NextIdVO](dgerr.ILLEGAL_OPERATION)
	}

	id, err := sequence.NextId(dc, nextIdReq.SeqName)
	if err != nil {
		return result.FailByError[*NextIdVO](err)
	}

	return result.Success[*NextIdVO](&NextIdVO{NextId: id})
}

func NextIds(gc *gin.Context, dc *dgctx.DgContext, nextIdsReq *NextIdsReq) *result.Result[*NextIdsVO] {
	internalCall := isInternalCall(gc)
	if !internalCall {
		gc.AbortWithStatus(http.StatusNotFound)
		return result.FailByError[*NextIdsVO](dgerr.ILLEGAL_OPERATION)
	}

	ids, err := sequence.NextIds(dc, nextIdsReq.SeqName, nextIdsReq.Count)
	if err != nil {
		return result.FailByError[*NextIdsVO](err)
	}

	return result.Success[*NextIdsVO](&NextIdsVO{NextIds: ids})
}

// isInternalCall 是否内部访问
func isInternalCall(gc *gin.Context) bool {
	internalV := gc.GetHeader("internal_val")
	return internalV == "61057154-c8f4-40af-bf9a-8c85c0c3ac93"
}
