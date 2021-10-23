package gorm

import (
	"context"
	"fmt"
	"github.com/xqk/good/pkg/ilog"
	"github.com/xqk/good/pkg/metric"
	"github.com/xqk/good/pkg/trace"
	"github.com/xqk/good/pkg/util/icolor"
	"strconv"
	"time"
)

// Handler ...
type Handler func(*Scope)

// Interceptor ...
type Interceptor func(*DSN, string, *Config) func(next Handler) Handler

func debugInterceptor(dsn *DSN, op string, options *Config) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(scope *Scope) {
			fmt.Printf("%-50s[%s] => %s\n", icolor.Green(dsn.Addr+"/"+dsn.DBName), time.Now().Format("04:05.000"), icolor.Green("Send: "+logSQL(scope.SQL, scope.SQLVars, true)))
			next(scope)
			if scope.HasError() {
				fmt.Printf("%-50s[%s] => %s\n", icolor.Red(dsn.Addr+"/"+dsn.DBName), time.Now().Format("04:05.000"), icolor.Red("Erro: "+scope.DB().Error.Error()))
			} else {
				fmt.Printf("%-50s[%s] => %s\n", icolor.Green(dsn.Addr+"/"+dsn.DBName), time.Now().Format("04:05.000"), icolor.Green("Affected: "+strconv.Itoa(int(scope.DB().RowsAffected))))
			}
		}
	}
}

func metricInterceptor(dsn *DSN, op string, options *Config) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(scope *Scope) {
			beg := time.Now()
			next(scope)
			cost := time.Since(beg)

			// error metric
			if scope.HasError() {
				metric.LibHandleCounter.WithLabelValues(metric.TypeGorm, dsn.DBName+"."+scope.TableName(), dsn.Addr, "ERR").Inc()
				// todo sql语句，需要转换成脱密状态才能记录到日志
				if scope.DB().Error != ErrRecordNotFound {
					options.logger.Error("mysql err", ilog.FieldErr(scope.DB().Error), ilog.FieldName(dsn.DBName+"."+scope.TableName()), ilog.FieldMethod(op))
				} else {
					options.logger.Warn("record not found", ilog.FieldErr(scope.DB().Error), ilog.FieldName(dsn.DBName+"."+scope.TableName()), ilog.FieldMethod(op))
				}
			} else {
				metric.LibHandleCounter.Inc(metric.TypeGorm, dsn.DBName+"."+scope.TableName(), dsn.Addr, "OK")
			}

			metric.LibHandleHistogram.WithLabelValues(metric.TypeGorm, dsn.DBName+"."+scope.TableName(), dsn.Addr).Observe(cost.Seconds())

			if options.SlowThreshold > time.Duration(0) && options.SlowThreshold < cost {
				options.logger.Error(
					"slow",
					ilog.FieldErr(errSlowCommand),
					ilog.FieldMethod(op),
					ilog.FieldExtMessage(logSQL(scope.SQL, scope.SQLVars, options.DetailSQL)),
					ilog.FieldAddr(dsn.Addr),
					ilog.FieldName(dsn.DBName+"."+scope.TableName()),
					ilog.FieldCost(cost),
				)
			}
		}
	}
}

func logSQL(sql string, args []interface{}, containArgs bool) string {
	if containArgs {
		return bindSQL(sql, args)
	}
	return sql
}

func traceInterceptor(dsn *DSN, op string, options *Config) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(scope *Scope) {
			if val, ok := scope.Get("_context"); ok {
				if ctx, ok := val.(context.Context); ok {
					span, _ := trace.StartSpanFromContext(
						ctx,
						"GORM", // TODO this op value is op or GORM
						trace.TagComponent("mysql"),
						trace.TagSpanKind("client"),
					)
					defer span.Finish()

					// 延迟执行 scope.CombinedConditionSql() 避免sqlVar被重复追加
					next(scope)

					span.SetTag("sql.inner", dsn.DBName)
					span.SetTag("sql.addr", dsn.Addr)
					span.SetTag("span.kind", "client")
					span.SetTag("peer.service", "mysql")
					span.SetTag("db.instance", dsn.DBName)
					span.SetTag("peer.address", dsn.Addr)
					span.SetTag("peer.statement", logSQL(scope.SQL, scope.SQLVars, options.DetailSQL))
					return
				}
			}

			next(scope)
		}
	}
}
