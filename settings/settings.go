package setting

import (
	"fmt"
	dgsys "github.com/darwinOrg/go-common/sys"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var appConf map[string]any

type DbPoolConf struct {
	Size     int
	Life     int
	IdleCons int
	IdleTime int
}

type AppInfo struct {
	Name        string
	Port        int
	MonitorPort int
}

func InitSetting() {
	profile := os.Getenv("profile")

	confRoot := GetConfRoot()
	confFile := confRoot + "/app.yml"

	if profile != "" {
		confFile = fmt.Sprintf("%s/app-%s.yml", confRoot, profile)
	}

	log.Printf("use profile:%s, conf: %s\n", profile, confFile)

	buff, err := os.ReadFile(confFile)
	if err != nil {
		log.Println(err)
		return
	}

	if dgsys.IsFormalProfile() {
		// 替换配置文件中的环境变量
		envMap := make(map[string]string)
		for _, env := range os.Environ() {
			pair := strings.SplitN(env, "=", 2)
			envMap[pair[0]] = pair[1]
		}

		re := regexp.MustCompile(`\${(.*?)}`)
		buff = []byte(re.ReplaceAllStringFunc(string(buff), func(match string) string {
			envName := match[2 : len(match)-1] // 移除 "${" 和 "}"
			envValue := envMap[envName]
			return envValue
		}))
	}

	if err := yaml.Unmarshal(buff, &appConf); err != nil {
		log.Println(err)
		return
	}
}

func GetConfRoot() string {
	confRoot := "./conf"
	testConfRoot := os.Getenv("test.conf.root")
	if testConfRoot != "" {
		return testConfRoot
	}
	return confRoot
}

func GetMysqlUrl() string {
	mysqlConf := getMapInfo(appConf, "mysql")

	return convertString(mysqlConf["url"])
}

func GetDbPoolConf() *DbPoolConf {
	mysqlConf := getMapInfo(appConf, "mysql")

	return &DbPoolConf{
		Size:     mysqlConf["size"].(int),
		Life:     mysqlConf["life"].(int),
		IdleCons: mysqlConf["idleCons"].(int),
		IdleTime: mysqlConf["idleTime"].(int),
	}
}

func GetAppInfo() *AppInfo {
	v, ok := appConf["app"]
	if !ok {
		return nil
	}

	mp := v.(map[string]any)
	return &AppInfo{
		Name:        GetFromMapString(mp, "name"),
		Port:        getFromMapInt(mp, "port"),
		MonitorPort: getFromMapInt(mp, "monitor"),
	}
}

func GetFromMapString(mpValue map[string]any, key string) string {
	v, ok := mpValue[key]
	if !ok {
		return ""
	}
	return v.(string)
}

func getFromMapInt(mpValue map[string]any, key string) int {
	v, ok := mpValue[key]
	if !ok {
		return 0
	}
	return v.(int)
}

func getMapInfo(conf map[string]any, key string) map[string]any {
	v, ok := conf[key]
	if !ok {
		return nil
	}
	return v.(map[string]any)
}

func convertString(raw any) string {
	v, ok := raw.(int)
	if ok {
		return strconv.Itoa(v)
	}

	return raw.(string)
}
