package web

import (
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-common/result"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/gin-gonic/gin"
	"go/types"
)

func health(gc *gin.Context, dc *dgctx.DgContext, req *types.Nil) *result.Result[types.Nil] {
	dglogger.Infof(dc, "health check.")

	return result.Success(types.Nil{})
}
