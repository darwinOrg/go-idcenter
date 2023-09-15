package web

import (
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-common/result"
	"github.com/gin-gonic/gin"
	"go-idcenter/sequence"
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
	id, err := sequence.NextId(dc, nextIdReq.SeqName)
	if err != nil {
		return result.FailByError[*NextIdVO](err)
	}

	return result.Success[*NextIdVO](&NextIdVO{NextId: id})
}

func NextIds(gc *gin.Context, dc *dgctx.DgContext, nextIdsReq *NextIdsReq) *result.Result[*NextIdsVO] {
	ids, err := sequence.NextIds(dc, nextIdsReq.SeqName, nextIdsReq.Count)
	if err != nil {
		return result.FailByError[*NextIdsVO](err)
	}

	return result.Success[*NextIdsVO](&NextIdsVO{NextIds: ids})
}
