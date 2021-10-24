package apollo

import (
	"fmt"

	"github.com/xqk/good/config"
	"github.com/xqk/good/encoding"
	"github.com/xqk/good/log"

	"github.com/apolloconfig/agollo/v4/storage"
)

type watcher struct {
	out      <-chan []*config.KeyValue
	cancelFn func()
}

type customChangeListener struct {
	in     chan<- []*config.KeyValue
	logger log.Logger
}

func (c *customChangeListener) onChange(
	namespace string, changes map[string]*storage.ConfigChange) []*config.KeyValue {

	kv := make([]*config.KeyValue, 0, 2)
	next := make(map[string]interface{})

	for key, change := range changes {
		resolve(genKey(namespace, key), change.NewValue, next)
	}

	f := format(namespace)
	codec := encoding.GetCodec(f)
	val, err := codec.Marshal(next)
	if err != nil {
		_ = c.logger.Log(log.LevelWarn,
			"msg",
			fmt.Sprintf("apollo could not handle namespace %s: %v", namespace, err),
		)
		return nil
	}
	kv = append(kv, &config.KeyValue{
		Key:    namespace,
		Value:  val,
		Format: f,
	})

	return kv
}

func (c *customChangeListener) OnChange(changeEvent *storage.ChangeEvent) {
	change := c.onChange(changeEvent.Namespace, changeEvent.Changes)
	if len(change) == 0 {
		return
	}

	c.in <- change
}

func (c *customChangeListener) OnNewestChange(changeEvent *storage.FullChangeEvent) {}

func newWatcher(a *apollo, logger log.Logger) (config.Watcher, error) {
	if logger == nil {
		logger = log.DefaultLogger
	}

	changeCh := make(chan []*config.KeyValue)
	a.client.AddChangeListener(&customChangeListener{in: changeCh, logger: logger})

	return &watcher{
		out: changeCh,
		cancelFn: func() {
			close(changeCh)
		},
	}, nil
}

// Next will be blocked until the Stop method is called
func (w *watcher) Next() ([]*config.KeyValue, error) {
	return <-w.out, nil
}

func (w *watcher) Stop() error {
	if w.cancelFn != nil {
		w.cancelFn()
	}

	return nil
}
