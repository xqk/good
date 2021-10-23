package idebug

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/tidwall/pretty"
	"github.com/xqk/good/pkg/util/icolor"
	"github.com/xqk/good/pkg/util/istring"

	"github.com/xqk/good/pkg/ilog"
)

var (
	isTestingMode     bool
	isDevelopmentMode = os.Getenv("OX_MODE") == "dev"
)

func init() {
	if isDevelopmentMode {
		ilog.DefaultLogger.SetLevel(ilog.DebugLevel)
		ilog.GoodLogger.SetLevel(ilog.DebugLevel)
	}
}

// IsTestingMode 判断是否在测试模式下
var onceTest = sync.Once{}

// IsTestingMode ...
func IsTestingMode() bool {
	onceTest.Do(func() {
		isTestingMode = flag.Lookup("test.v") != nil
	})

	return isTestingMode
}

// IsDevelopmentMode 判断是否是生产模式
func IsDevelopmentMode() bool {
	return isDevelopmentMode || isTestingMode
}

// IfPanic ...
func IfPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// PrettyJsonPrint ...
func PrettyJsonPrint(message string, obj interface{}) {
	if !IsDevelopmentMode() {
		return
	}
	fmt.Printf("%s => %s\n",
		icolor.Red(message),
		pretty.Color(
			pretty.Pretty([]byte(istring.PrettyJson(obj))),
			pretty.TerminalStyle,
		),
	)
}

// PrettyJsonByte ...
func PrettyJsonByte(obj interface{}) string {
	return string(pretty.Color(pretty.Pretty([]byte(istring.Json(obj))), pretty.TerminalStyle))
}

// PrettyKV ...
func PrettyKV(key string, val string) {
	fmt.Printf("%-50s => %s\n", icolor.Red(key), icolor.Green(val))
}

// PrettyKV ...
func PrettyKVWithPrefix(prefix string, key string, val string) {
	fmt.Printf(prefix+" %-30s => %s\n", icolor.Red(key), icolor.Blue(val))
}

// PrettyMap ...
func PrettyMap(data map[string]interface{}) {
	for key, val := range data {
		fmt.Printf("%-20s : %s\n", icolor.Red(key), fmt.Sprintf("%+v", val))
	}
}

// GetCurrentDirectory ...
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) // 返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		panic(err)
	}
	return strings.Replace(dir, "\\", "/", -1) // 将\替换成/
}
