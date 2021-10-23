package registry

import (
	"context"
	"github.com/xqk/good/pkg/ilog"
	"github.com/xqk/good/pkg/server"
)

// Nop registry, used for local development/debugging
type Local struct{}

// ListServices ...
func (n Local) ListServices(ctx context.Context, s string, s2 string) ([]*server.ServiceInfo, error) {
	panic("implement me")
}

// WatchServices ...
func (n Local) WatchServices(ctx context.Context, s string, s2 string) (chan Endpoints, error) {
	panic("implement me")
}

// RegisterService ...
func (n Local) RegisterService(ctx context.Context, si *server.ServiceInfo) error {
	ilog.Info("register service locally", ilog.FieldMod("registry"), ilog.FieldName(si.Name), ilog.FieldAddr(si.Label()))
	return nil
}

// UnregisterService ...
func (n Local) UnregisterService(ctx context.Context, si *server.ServiceInfo) error {
	ilog.Info("unregister service locally", ilog.FieldMod("registry"), ilog.FieldName(si.Name), ilog.FieldAddr(si.Label()))
	return nil
}

// Close ...
func (n Local) Close() error { return nil }

// Close ...
func (n Local) Kind() string { return "local" }
