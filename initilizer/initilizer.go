package initilizer

import (
	setting "com.startrek/go-idcenter/settings"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rolandhe/daog"
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
