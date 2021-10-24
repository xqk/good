package filter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xqk/good/registry"
	"github.com/xqk/good/selector"
)

func TestVersion(t *testing.T) {
	f := Version("v2.0.0")
	var nodes []selector.Node
	nodes = append(nodes, selector.NewNode(
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:9090",
			Name:      "helloworld",
			Version:   "v1.0.0",
			Endpoints: []string{"http://127.0.0.1:9090"},
		}))

	nodes = append(nodes, selector.NewNode(
		"127.0.0.2:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.2:9090",
			Name:      "helloworld",
			Version:   "v2.0.0",
			Endpoints: []string{"http://127.0.0.2:9090"},
		}))

	n := f(context.Background(), nodes)
	assert.Equal(t, 1, len(n))
	assert.Equal(t, "127.0.0.2:9090", n[0].Address())
}
