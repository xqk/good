package idebug

import (
	"fmt"

	"github.com/tidwall/pretty"
	"github.com/xqk/good/pkg/util/icolor"
	"github.com/xqk/good/pkg/util/istring"
)

// DebugObject ...
func PrintObject(message string, obj interface{}) {
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

// DebugBytes ...
func DebugBytes(obj interface{}) string {
	return string(pretty.Color(pretty.Pretty([]byte(istring.Json(obj))), pretty.TerminalStyle))
}

// PrintKV ...
func PrintKV(key string, val string) {
	if !IsDevelopmentMode() {
		return
	}
	fmt.Printf("%-50s => %s\n", icolor.Red(key), icolor.Green(val))
}

// PrettyKVWithPrefix ...
func PrintKVWithPrefix(prefix string, key string, val string) {
	if !IsDevelopmentMode() {
		return
	}
	fmt.Printf("%-8s]> %-30s => %s\n", prefix, icolor.Red(key), icolor.Blue(val))
}

// PrintMap ...
func PrintMap(data map[string]interface{}) {
	if !IsDevelopmentMode() {
		return
	}
	for key, val := range data {
		fmt.Printf("%-20s : %s\n", icolor.Red(key), fmt.Sprintf("%+v", val))
	}
}
