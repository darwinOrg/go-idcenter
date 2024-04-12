package web

import (
	"com.startrek/go-idcenter/sequence"
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-common/result"
	"github.com/darwinOrg/go-web/wrapper"
	"github.com/gin-gonic/gin"
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

func RegisterAll(engine *gin.Engine) {
	rg := engine.Group("/internal/v1/id")

	wrapper.Get(&wrapper.RequestHolder[NextIdReq, *result.Result[*NextIdVO]]{
		RouterGroup:  rg,
		RelativePath: "/next-id",
		NonLogin:     true,
		BizHandler: func(_ *gin.Context, dc *dgctx.DgContext, nextIdReq *NextIdReq) *result.Result[*NextIdVO] {
			id, err := sequence.NextId(dc, nextIdReq.SeqName)
			if err != nil {
				return result.FailByError[*NextIdVO](err)
			}

			return result.Success[*NextIdVO](&NextIdVO{NextId: id})
		},
	})

	wrapper.Get(&wrapper.RequestHolder[NextIdsReq, *result.Result[*NextIdsVO]]{
		RouterGroup:  rg,
		RelativePath: "/next-ids",
		NonLogin:     true,
		BizHandler: func(_ *gin.Context, dc *dgctx.DgContext, nextIdsReq *NextIdsReq) *result.Result[*NextIdsVO] {
			ids, err := sequence.NextIds(dc, nextIdsReq.SeqName, nextIdsReq.Count)
			if err != nil {
				return result.FailByError[*NextIdsVO](err)
			}

			return result.Success[*NextIdsVO](&NextIdsVO{NextIds: ids})
		},
	})
}
