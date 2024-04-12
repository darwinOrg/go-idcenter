package sequence

import (
	"com.startrek/go-idcenter/initilizer"
	sequence "com.startrek/go-idcenter/sequence/dal"
	dgctx "github.com/darwinOrg/go-common/context"
	dgerr "github.com/darwinOrg/go-common/enums/error"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/rolandhe/daog"
	txrequest "github.com/rolandhe/daog/tx"
	"sync"
)

type concurrentMultiLock struct {
	sync.Map
}

type cachedNamedSeq struct {
	sync.Mutex
	cachedSeqInfo
}

func (cl *concurrentMultiLock) getLockByName(seqName string) *cachedNamedSeq {
	if v, ok := cl.Load(seqName); ok {
		return v.(*cachedNamedSeq)
	}
	seq := &cachedNamedSeq{}
	seq.seqName = seqName
	act, loaded := cl.LoadOrStore(seqName, seq)
	if loaded {
		return act.(*cachedNamedSeq)
	}
	return seq
}

func (cs *cachedNamedSeq) getNextOneId(dc *dgctx.DgContext) (int64, error) {
	ids, err := cs.batchIds(dc, 1)
	if err != nil {
		return 0, err
	}
	return ids[0], nil
}

func (cs *cachedNamedSeq) batchIds(dc *dgctx.DgContext, count int) ([]int64, error) {
	cs.Lock()
	defer cs.Unlock()

	ids := cs.nextIds(count)
	if len(ids) == count {
		return ids, nil
	}

	// 如果没获取到足够的ID 则需要重新初始化缓存
	more := count - len(ids)
	err := refreshCacheFromDB(dc, &cs.cachedSeqInfo, more)
	if err != nil {
		dglogger.Errorf(dc, "NextIds refreshCacheFromDB err:%v", err)
		return nil, err
	}

	moreIds := cs.nextIds(more)
	totalIds := make([]int64, 0, count)
	if len(ids) > 0 {
		totalIds = append(totalIds, ids...)
	}
	totalIds = append(totalIds, moreIds...)
	return totalIds, nil
}

func refreshCacheFromDB(dc *dgctx.DgContext, seq *cachedSeqInfo, more int) error {
	dglogger.Infof(dc, "SequenceCache需要刷新,cache=%s", seq.seqName)

	tcCreate := func() (*daog.TransContext, error) {
		return daog.NewTransContext(initilizer.GlobalDatasource, txrequest.RequestWrite, dc.TraceId)
	}

	return daog.AutoTrans(tcCreate, func(tc *daog.TransContext) error {
		mc := daog.NewMatcher()
		mc.Eq(sequence.SequenceFields.SeqName, seq.seqName)
		exist, err := sequence.SequenceDao.QueryOneMatcherForUpdate(tc, mc, false)
		if err != nil {
			return err
		}

		if exist == nil {
			return dgerr.RECORD_NOT_EXISTS
		}

		curId := exist.CurValue
		cacheSize := exist.CacheSize
		maxId, err := nextValueByDB(tc, exist, int(cacheSize)+more)
		if err != nil {
			dglogger.Errorf(dc, "refreshCacheFromDB nextValueByDB err:%v", err)
			return err
		}

		seq.refresh(curId, maxId, int(exist.Step))

		return nil
	})

}
func nextValueByDB(tc *daog.TransContext, seq *sequence.Sequence, size int) (int64, error) {
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
		return 0, dgerr.SimpleDgError("nextValueByDB UpdateByModifier err")
	}

	return newValue, nil
}
