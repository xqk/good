// Package pkg
// @Description: 外部应用程序可以使用的库代码
package pkg

import (
	"fmt"
	"github.com/xqk/good/pkg/constant"
	"github.com/xqk/good/pkg/util/icolor"
	"github.com/xqk/good/pkg/util/itime"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const goodVersion = "0.1.0"

var (
	startTime string
	goVersion string
)

// build info
/*

 */
var (
	appName         string
	appID           string
	hostName        string
	buildAppVersion string
	buildUser       string
	buildHost       string
	buildStatus     string
	buildTime       string
)

func init() {
	if appName == "" {
		appName = os.Getenv(constant.EnvAppName)
		if appName == "" {
			appName = filepath.Base(os.Args[0])
		}
	}

	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	hostName = name
	startTime = itime.TS.Format(time.Now())
	SetBuildTime(buildTime)
	goVersion = runtime.Version()
	InitEnv()
}

// Name gets application name.
func Name() string {
	return appName
}

//SetName set app anme
func SetName(s string) {
	appName = s
}

//AppID get appID
func AppID() string {
	if appID == "" {
		return "1234567890" //default appid when APP_ID Env var not set
	}
	return appID
}

//SetAppID set appID
func SetAppID(s string) {
	appID = s
}

//AppVersion get buildAppVersion
func AppVersion() string {
	return buildAppVersion
}

//appVersion not defined
// func SetAppVersion(s string) {
// 	appVersion = s
// }

//GoodVersion get goodVersion
func GoodVersion() string {
	return goodVersion
}

//BuildTime get buildTime
func BuildTime() string {
	return buildTime
}

//BuildUser get buildUser
func BuildUser() string {
	return buildUser
}

//BuildHost get buildHost
func BuildHost() string {
	return buildHost
}

//SetBuildTime set buildTime
func SetBuildTime(param string) {
	buildTime = strings.Replace(param, "--", " ", 1)
}

// HostName get host name
func HostName() string {
	return hostName
}

//StartTime get start time
func StartTime() string {
	return startTime
}

//GoVersion get go version
func GoVersion() string {
	return goVersion
}

func LogDir() string {
	// LogDir gets application log directory.
	logDir := AppLogDir()
	if logDir == "" {
		if appPodIP != "" && appPodName != "" {
			// k8s 环境
			return fmt.Sprintf("/home/www/logs/applogs/%s/%s/", Name(), appPodName)
		}
		return fmt.Sprintf("/home/www/logs/applogs/%s/%s/", Name(), appInstance)
	}
	return fmt.Sprintf("%s/%s/%s/", logDir, Name(), appInstance)
}

//
// PrintVersion
// @Description: 打印格式化后的版本信息
//
func PrintVersion() {
	fmt.Printf("%-8s]> %-30s => %s\n", "good", icolor.Red("name"), icolor.Blue(appName))
	fmt.Printf("%-8s]> %-30s => %s\n", "good", icolor.Red("appID"), icolor.Blue(appID))
	fmt.Printf("%-8s]> %-30s => %s\n", "good", icolor.Red("region"), icolor.Blue(AppRegion()))
	fmt.Printf("%-8s]> %-30s => %s\n", "good", icolor.Red("zone"), icolor.Blue(AppZone()))
	fmt.Printf("%-8s]> %-30s => %s\n", "good", icolor.Red("appVersion"), icolor.Blue(buildAppVersion))
	fmt.Printf("%-8s]> %-30s => %s\n", "good", icolor.Red("goodVersion"), icolor.Blue(goodVersion))
	fmt.Printf("%-8s]> %-30s => %s\n", "good", icolor.Red("buildUser"), icolor.Blue(buildUser))
	fmt.Printf("%-8s]> %-30s => %s\n", "good", icolor.Red("buildHost"), icolor.Blue(buildHost))
	fmt.Printf("%-8s]> %-30s => %s\n", "good", icolor.Red("buildTime"), icolor.Blue(BuildTime()))
	fmt.Printf("%-8s]> %-30s => %s\n", "good", icolor.Red("buildStatus"), icolor.Blue(buildStatus))
}
