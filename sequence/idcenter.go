package sequence

import (
	dgctx "github.com/darwinOrg/go-common/context"
)

var cacheLockMap = &concurrentMultiLock{}

func NextId(dc *dgctx.DgContext, seqName string) (int64, error) {
	return cacheLockMap.getLockByName(seqName).getNextOneId(dc)
}

func NextIds(dc *dgctx.DgContext, seqName string, count int) ([]int64, error) {
	return cacheLockMap.getLockByName(seqName).batchIds(dc, count)
}
