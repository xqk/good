package filter

import (
	"context"

	"github.com/xqk/good/selector"
)

// Version is verion filter.
func Version(version string) selector.Filter {
	return func(_ context.Context, nodes []selector.Node) []selector.Node {
		filters := make([]selector.Node, 0, len(nodes))
		for _, n := range nodes {
			if n.Version() == version {
				filters = append(filters, n)
			}
		}
		return filters
	}
}
