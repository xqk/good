package mongox

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/xqk/good/pkg/application"
	"github.com/xqk/good/pkg/ilog"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func init() {
	// govern.RegisterStatSnapper("mongo", Stats)
	// govern.RegisterConfSnapper("mongo", Configs)
	http.HandleFunc("/debug/mongo/stats", func(w http.ResponseWriter, r *http.Request) {
		type mongoStatus struct {
			application.RuntimeStats
			Mongos map[string]interface{} `json:"mongos"`
		}
		var rets = mongoStatus{
			RuntimeStats: application.NewRuntimeStats(),
			Mongos:       make(map[string]interface{}, 0),
		}
		Range(func(name string, cc *mongo.Client) bool {
			rets.Mongos[name] = map[string]interface{}{
				"numberSessionsInProgress": cc.NumberSessionsInProgress(),
			}
			return true
		})

		_ = jsoniter.NewEncoder(w).Encode(rets)
	})
}

var _logger = ilog.GoodLogger.With(ilog.FieldMod("mongodb"))
