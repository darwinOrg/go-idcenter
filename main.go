package main

import (
	setting "com.startrek/go-idcenter/settings"
	"com.startrek/go-idcenter/web"
	"fmt"
	"github.com/darwinOrg/go-monitor"
	"github.com/darwinOrg/go-web/wrapper"
)

func main() {
	appInfo := setting.GetAppInfo()
	monitor.Start(appInfo.Name, appInfo.MonitorPort)

	engine := wrapper.DefaultEngine()
	web.RegisterAll(engine)
	_ = engine.Run(fmt.Sprintf(":%d", appInfo.Port))
}
