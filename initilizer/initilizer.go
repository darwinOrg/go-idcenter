package initilizer

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/rolandhe/daog"
	setting "go-idcenter/settings"
)

var GlobalDatasource = initDatasource()

func initDatasource() daog.Datasource {
	setting.InitSetting()

	poolConf := setting.GetDbPoolConf()

	dbConf := &daog.DbConf{
		DbUrl:    setting.GetMysqlUrl(),
		Size:     poolConf.Size,
		Life:     poolConf.Life,
		IdleCons: poolConf.IdleCons,
		IdleTime: poolConf.IdleTime,
		LogSQL:   true,
	}
	datasource, err := daog.NewDatasource(dbConf)
	if err != nil {
		panic(err)
	}
	return datasource
}
