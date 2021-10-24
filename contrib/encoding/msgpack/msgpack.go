package msgpack

import (
	"github.com/xqk/good/encoding"
	"github.com/vmihailenco/msgpack/v5"
)

// Name is the name registered for the msgpack compressor.
const Name = "msgpack"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with msgpack.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}

func (codec) Name() string {
	return Name
}
