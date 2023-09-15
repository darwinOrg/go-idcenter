package test

import (
	"com.startrek/go-idcenter/sequence"
	dgctx "github.com/darwinOrg/go-common/context"
	"testing"
)

func TestNextId(t *testing.T) {
	ctx := &dgctx.DgContext{
		TraceId: "12314",
	}

	id, err := sequence.NextId(ctx, "media")
	t.Logf("id: %d,err: %v", id, err)

	ids, err := sequence.NextIds(ctx, "media", 300)
	t.Logf("ids: %v,err: %v", ids, err)
}

func TestNextIds(t *testing.T) {
	ctx := &dgctx.DgContext{
		TraceId: "12314",
	}

	ids, err := sequence.NextIds(ctx, "media", 300)
	t.Logf("ids: %v,err: %v", ids, err)
}
