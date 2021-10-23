package governor

import (
	"encoding/json"
	"github.com/xqk/good/pkg/util/istring"
	"log"
	"net/http"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/xqk/good/pkg"
	"github.com/xqk/good/pkg/conf"
)

func init() {
	conf.OnLoaded(func(c *conf.Configuration) {
		log.Print("hook config, init runtime(governor)")

	})

	registerHandlers()
}

func registerHandlers() {
	HandleFunc("/configs", func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		if r.URL.Query().Get("pretty") == "true" {
			encoder.SetIndent("", "    ")
		}
		encoder.Encode(conf.Traverse("."))
	})

	HandleFunc("/debug/config", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write(istring.PrettyJSONBytes(conf.Traverse(".")))
	})

	HandleFunc("/debug/env", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_ = jsoniter.NewEncoder(w).Encode(os.Environ())
	})

	HandleFunc("/build/info", func(w http.ResponseWriter, r *http.Request) {
		serverStats := map[string]string{
			"name":        pkg.Name(),
			"appID":       pkg.AppID(),
			"appMode":     pkg.AppMode(),
			"appVersion":  pkg.AppVersion(),
			"goodVersion": pkg.GoodVersion(),
			"buildUser":   pkg.BuildUser(),
			"buildHost":   pkg.BuildHost(),
			"buildTime":   pkg.BuildTime(),
			"startTime":   pkg.StartTime(),
			"hostName":    pkg.HostName(),
			"goVersion":   pkg.GoVersion(),
		}
		_ = jsoniter.NewEncoder(w).Encode(serverStats)
	})
}
