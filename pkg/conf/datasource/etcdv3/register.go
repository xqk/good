package etcdv3

import (
	"github.com/xqk/good/pkg/client/etcdv3"
	"github.com/xqk/good/pkg/conf"
	"github.com/xqk/good/pkg/flag"
	"github.com/xqk/good/pkg/ilog"
	"github.com/xqk/good/pkg/util/inet"
)

// DataSourceEtcdv3 defines etcdv3 scheme
const DataSourceEtcdv3 = "etcdv3"

func init() {
	conf.Register(DataSourceEtcdv3, func() conf.DataSource {
		var (
			configAddr = flag.String("config")
		)
		if configAddr == "" {
			ilog.Panic("new apollo dataSource, configAddr is empty")
			return nil
		}
		// configAddr is a string in this format:
		// etcdv3://ip:port?basicAuth=true&username=XXX&password=XXX&key=XXX&certFile=XXX&keyFile=XXX&caCert=XXX&secure=XXX

		urlObj, err := inet.ParseURL(configAddr)
		if err != nil {
			ilog.Panic("parse configAddr error", ilog.FieldErr(err))
			return nil
		}
		etcdConf := etcdv3.DefaultConfig()
		etcdConf.Endpoints = []string{urlObj.Host}
		etcdConf.BasicAuth = urlObj.QueryBool("basicAuth", false)
		etcdConf.Secure = urlObj.QueryBool("secure", false)
		etcdConf.CertFile = urlObj.Query().Get("certFile")
		etcdConf.KeyFile = urlObj.Query().Get("keyFile")
		etcdConf.CaCert = urlObj.Query().Get("caCert")
		etcdConf.UserName = urlObj.Query().Get("username")
		etcdConf.Password = urlObj.Query().Get("password")
		return NewDataSource(etcdConf.MustBuild(), urlObj.Query().Get("key"))
	})
}
