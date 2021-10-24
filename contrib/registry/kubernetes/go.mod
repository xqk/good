module github.com/xqk/good/contrib/registry/kubernetes/v2

go 1.16

require (
	github.com/json-iterator/go v1.1.11
	github.com/xqk/good v0.0.0
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
)

replace github.com/xqk/good => ../../../
