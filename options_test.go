package good

import (
	"context"
	"github.com/stretchr/testify/assert"
	xlog "github.com/xqk/good/log"
	"github.com/xqk/good/registry"
	"github.com/xqk/good/transport"
	"log"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestID(t *testing.T) {
	o := &options{}
	v := "123"
	ID(v)(o)
	assert.Equal(t, v, o.id)
}

func TestName(t *testing.T) {
	o := &options{}
	v := "abc"
	Name(v)(o)
	assert.Equal(t, v, o.name)
}

func TestVersion(t *testing.T) {
	o := &options{}
	v := "123"
	Version(v)(o)
	assert.Equal(t, v, o.version)
}

func TestMetadata(t *testing.T) {
	o := &options{}
	v := map[string]string{
		"a": "1",
		"b": "2",
	}
	Metadata(v)(o)
	assert.Equal(t, v, o.metadata)
}

func TestEndpoint(t *testing.T) {
	o := &options{}
	v := []*url.URL{
		{Host: "example.com"},
		{Host: "foo.com"},
	}
	Endpoint(v...)(o)
	assert.Equal(t, v, o.endpoints)
}

func TestContext(t *testing.T) {
	type ctxKey = struct{}
	o := &options{}
	v := context.WithValue(context.TODO(), ctxKey{}, "b")
	Context(v)(o)
	assert.Equal(t, v, o.ctx)
}

func TestLogger(t *testing.T) {
	o := &options{}
	v := xlog.NewStdLogger(log.Writer())
	Logger(v)(o)
	assert.Equal(t, xlog.NewHelper(v), o.logger)
}

type mockServer struct{}

func (m *mockServer) Start(ctx context.Context) error { return nil }
func (m *mockServer) Stop(ctx context.Context) error  { return nil }

func TestServer(t *testing.T) {
	o := &options{}
	v := []transport.Server{
		&mockServer{}, &mockServer{},
	}
	Server(v...)(o)
	assert.Equal(t, v, o.servers)
}

type mockSignal struct{}

func (m *mockSignal) String() string { return "sig" }
func (m *mockSignal) Signal()        {}

func TestSignal(t *testing.T) {
	o := &options{}
	v := []os.Signal{
		&mockSignal{}, &mockSignal{},
	}
	Signal(v...)(o)
	assert.Equal(t, v, o.sigs)
}

type mockRegistrar struct{}

func (m *mockRegistrar) Register(ctx context.Context, service *registry.ServiceInstance) error {
	return nil
}

func (m *mockRegistrar) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	return nil
}

func TestRegistrar(t *testing.T) {
	o := &options{}
	v := &mockRegistrar{}
	Registrar(v)(o)
	assert.Equal(t, v, o.registrar)
}

func TestRegistrarTimeout(t *testing.T) {
	o := &options{}
	v := time.Duration(123)
	RegistrarTimeout(v)(o)
	assert.Equal(t, v, o.registrarTimeout)
}
