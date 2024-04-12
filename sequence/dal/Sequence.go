package sequence

import (
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/ttypes"
)

var SequenceFields = struct {
	SeqName   string
	CurValue  string
	Step      string
	CacheSize string
	GmtUpdate string
}{
	"seq_name",
	"cur_value",
	"step",
	"cache_size",
	"gmt_update",
}

var SequenceMeta = &daog.TableMeta[Sequence]{
	Table: "sequence",
	Columns: []string{
		"seq_name",
		"cur_value",
		"step",
		"cache_size",
		"gmt_update",
	},
	AutoColumn: "",
	LookupFieldFunc: func(columnName string, ins *Sequence, point bool) any {
		if "seq_name" == columnName {
			if point {
				return &ins.SeqName
			}
			return ins.SeqName
		}
		if "cur_value" == columnName {
			if point {
				return &ins.CurValue
			}
			return ins.CurValue
		}
		if "step" == columnName {
			if point {
				return &ins.Step
			}
			return ins.Step
		}
		if "cache_size" == columnName {
			if point {
				return &ins.CacheSize
			}
			return ins.CacheSize
		}
		if "gmt_update" == columnName {
			if point {
				return &ins.GmtUpdate
			}
			return ins.GmtUpdate
		}

		return nil
	},
}

var SequenceDao daog.QuickDao[Sequence] = &struct {
	daog.QuickDao[Sequence]
}{
	daog.NewBaseQuickDao(SequenceMeta),
}

type Sequence struct {
	SeqName   string
	CurValue  int64
	Step      int32
	CacheSize int32
	GmtUpdate ttypes.NormalDatetime
}
