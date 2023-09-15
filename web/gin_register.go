package web

import (
	"github.com/darwinOrg/go-common/result"
	"github.com/darwinOrg/go-web/wrapper"
	"github.com/gin-gonic/gin"
	"go/types"
)

func RegisterAll(engine *gin.Engine) {
	healthGroup := engine.Group("/health")
	wrapper.Get(&wrapper.RequestHolder[types.Nil, *result.Result[types.Nil]]{
		RouterGroup:  healthGroup,
		RelativePath: "",
		NonLogin:     true,
		BizHandler:   health,
	})

	registerApi(engine)
}

func registerApi(engine *gin.Engine) {
	apiGroup := engine.Group("/internal/id")

	wrapper.Post(&wrapper.RequestHolder[NextIdReq, *result.Result[*NextIdVO]]{
		RouterGroup:  apiGroup,
		RelativePath: "/next-id",
		NonLogin:     true,
		BizHandler:   NextId,
	})

	wrapper.Post(&wrapper.RequestHolder[NextIdsReq, *result.Result[*NextIdsVO]]{
		RouterGroup:  apiGroup,
		RelativePath: "/next-ids",
		NonLogin:     true,
		BizHandler:   NextIds,
	})
}
