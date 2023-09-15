package sequence

import "math"

type SequenceCache struct {
	seqName string // 序列名称
	curId   int64  // 缓存ID当前值
	maxId   int64  // 缓存ID最大值 (要注意该值为开区间不能使用)
	step    int    // 缓存ID步长
}

func (sc *SequenceCache) nextIds(count int) []int64 {
	c := sc.maxId - sc.curId
	size := math.Min(float64(c), float64(count))
	cur := sc.curId
	sc.curId += int64(size)

	var retArr []int64
	for i := cur; i < sc.curId; i++ {
		retArr = append(retArr, i)
	}

	return retArr
}

func (sc *SequenceCache) refresh(curId int64, maxId int64, step int) {
	sc.curId = curId
	sc.maxId = maxId
	sc.step = step
}
