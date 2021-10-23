package iattr

import "errors"

// Attributes 属性
type Attributes struct {
	m map[interface{}]interface{}
}

var (
	// ErrInvalidKVPairs 无效的key
	ErrInvalidKVPairs = errors.New("invalid kv pairs")
)

// New 创建一个属性结构对象
func New(kvs ...interface{}) *Attributes {
	if len(kvs)%2 != 0 {
		panic(ErrInvalidKVPairs)
	}
	a := &Attributes{m: make(map[interface{}]interface{}, len(kvs)/2)}
	for i := 0; i < len(kvs)/2; i++ {
		a.m[kvs[i*2]] = kvs[i*2+1]
	}
	return a
}

// WithValues 附加一组属性值
func (a *Attributes) WithValues(kvs ...interface{}) *Attributes {
	if len(kvs)%2 != 0 {
		panic(ErrInvalidKVPairs)
	}
	n := &Attributes{m: make(map[interface{}]interface{}, len(a.m)+len(kvs)/2)}
	for k, v := range a.m {
		n.m[k] = v
	}
	for i := 0; i < len(kvs)/2; i++ {
		n.m[kvs[i*2]] = kvs[i*2+1]
	}
	return n
}

// Value 获取属性的某个值
func (a *Attributes) Value(key interface{}) interface{} {
	return a.m[key]
}
