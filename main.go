package main

import (
	setting "com.startrek/go-idcenter/settings"
	"com.startrek/go-idcenter/web"
	"fmt"
	"github.com/darwinOrg/go-monitor"
	"github.com/darwinOrg/go-web/wrapper"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	appInfo := setting.GetAppInfo()

	monitor.Start(appInfo.Name, appInfo.MonitorPort)

	engine := createEngine()
	checkErr(engine.Run(fmt.Sprintf(":%d", appInfo.Port)), "start go idcenter server error")
}

func createEngine() *gin.Engine {
	engine := wrapper.DefaultEngine()
	engine.SetTrustedProxies(nil)
	web.RegisterAll(engine)
	return engine
}
func checkErr(err error, msg string) {
	if err == nil {
		return
	}
	log.Printf("ERROR: %s: %s\n", msg, err)
	os.Exit(1)
}
