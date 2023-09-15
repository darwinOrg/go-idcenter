package sequence

type cachedSeqInfo struct {
	seqName string // 序列名称
	curId   int64  // 缓存ID当前值
	maxId   int64  // 缓存ID最大值 (要注意该值为开区间不能使用)
	step    int    // 缓存ID步长
}

func (sc *cachedSeqInfo) nextIds(count int) []int64 {
	c := sc.maxId - sc.curId
	if c <= 0 {
		return nil
	}
	size := int(c)
	if count < size {
		size = count
	}
	cur := sc.curId
	sc.curId += int64(size)

	retArr := make([]int64, 0, size)
	for i := cur; i < sc.curId; i++ {
		retArr = append(retArr, i)
	}

	return retArr
}

func (sc *cachedSeqInfo) refresh(curId int64, maxId int64, step int) {
	sc.curId = curId
	sc.maxId = maxId
	sc.step = step
}
