package sequence

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dgerr "github.com/darwinOrg/go-common/enums/error"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/rolandhe/daog"
	txrequest "github.com/rolandhe/daog/tx"
	"go-idcenter/initilizer"
	sequence "go-idcenter/sequence/dal"
	"sync"
)

var cacheMap = make(map[string]*SequenceCache)

var cacheLockMap = &sync.Map{}

func NextId(dc *dgctx.DgContext, seqName string) (int64, error) {
	ids, err := NextIds(dc, seqName, 1)
	if err != nil {
		dglogger.Errorf(dc, "NextId NextIds err:%v", err)
		return 0, err
	}

	return ids[0], nil
}

func NextIds(dc *dgctx.DgContext, seqName string, count int) ([]int64, error) {
	if count < 1 {
		return nil, dgerr.ARGUMENT_NOT_VALID
	}

	// 加锁 释放锁
	lock := getCacheLock(seqName)
	lock.Lock()
	defer lock.Unlock()

	cache, err := getCacheCreateIfAbsent(dc, seqName)
	if err != nil {
		dglogger.Errorf(dc, "NextIds getCacheCreateIfAbsent err:%v", err)
		return nil, err
	}
	// 如果正常获取了ID 则直接返回
	ids := cache.nextIds(count)
	if len(ids) == count {
		return ids, nil
	}

	// 如果没获取到足够的ID 则需要重新初始化缓存
	more := count - len(ids)
	err = refreshCache(dc, cache, more)
	if err != nil {
		dglogger.Errorf(dc, "NextIds refreshCache err:%v", err)
		return nil, err
	}

	moreIds := cache.nextIds(more)
	var totalIds []int64
	totalIds = append(totalIds, ids...)
	totalIds = append(totalIds, moreIds...)
	return totalIds, nil
}

func getCacheLock(seqName string) *sync.Mutex {
	value, ok := cacheLockMap.Load(seqName)
	if ok {
		return value.(*sync.Mutex)
	}

	lock := &sync.Mutex{}
	newLock, _ := cacheLockMap.LoadOrStore(seqName, lock)
	return newLock.(*sync.Mutex)
}

func getCacheCreateIfAbsent(dc *dgctx.DgContext, seqName string) (*SequenceCache, error) {
	val, ok := cacheMap[seqName]
	if ok {
		return val, nil
	}

	dglogger.Infof(dc, "SequenceCache未找到,seqName=%s", seqName)
	newCache := &SequenceCache{seqName: seqName}
	err := refreshCache(dc, newCache, 0)
	if err != nil {
		dglogger.Errorf(dc, "getCacheCreateIfAbsent refreshCache err:%v", err)
		return nil, err
	}

	cacheMap[seqName] = newCache
	return newCache, nil
}

func refreshCache(dc *dgctx.DgContext, cache *SequenceCache, more int) error {
	dglogger.Infof(dc, "SequenceCache需要刷新,cache=%v", cache)

	tcCreate := func() (*daog.TransContext, error) {
		return daog.NewTransContext(initilizer.GlobalDatasource, txrequest.RequestWrite, dc.TraceId)
	}

	return daog.AutoTrans(tcCreate, func(tc *daog.TransContext) error {
		mc := daog.NewMatcher()
		mc.Eq(sequence.SequenceFields.SeqName, cache.seqName)
		exist, err := sequence.SequenceDao.QueryOneMatcherForUpdate(tc, mc, false)
		if err != nil {
			return err
		}

		if exist == nil {
			return dgerr.RECORD_NOT_EXISTS
		}

		curId := exist.CurValue
		cacheSize := exist.CacheSize
		maxId, err := nextValue(tc, exist, int(cacheSize)+more)
		if err != nil {
			dglogger.Errorf(dc, "refreshCache nextValue err:%v", err)
			return err
		}

		cache.refresh(curId, maxId, int(exist.Step))

		return nil
	})

}

func nextValue(tc *daog.TransContext, seq *sequence.Sequence, size int) (int64, error) {
	seqName := seq.SeqName
	curValue := seq.CurValue

	// 计算目标值
	step := seq.Step
	newValue := curValue + int64(int32(size)*step)

	mf := daog.NewModifier()
	mf.Add(sequence.SequenceFields.CurValue, newValue)

	mc := daog.NewMatcher()
	mc.Eq(sequence.SequenceFields.SeqName, seqName)

	// 获取成功
	affectRow, err := sequence.SequenceDao.UpdateByModifier(tc, mf, mc)
	if err != nil {
		return 0, err
	}

	if affectRow < 1 {
		return 0, dgerr.SimpleDgError("nextValue UpdateByModifier err")
	}

	return newValue, nil
}
