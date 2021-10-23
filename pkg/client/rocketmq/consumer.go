package rocketmq

import (
	"context"
	"github.com/xqk/good/pkg/ilog"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/xqk/good/pkg/defers"
	"github.com/xqk/good/pkg/istats"
)

type PushConsumer struct {
	rocketmq.PushConsumer
	name string
	ConsumerConfig

	subscribers  map[string]func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)
	interceptors []primitive.Interceptor
	fInfo        FlowInfo
}

func (conf *ConsumerConfig) Build() *PushConsumer {
	name := conf.Name
	if _, ok := _consumers.Load(name); ok {
		ilog.Panic("duplicated load", ilog.String("name", name))
	}

	ilog.Debug("rocketmq's config: ", ilog.String("name", name), ilog.Any("conf", conf))

	pc := &PushConsumer{
		name:           name,
		ConsumerConfig: *conf,
		subscribers:    make(map[string]func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)),
		interceptors:   []primitive.Interceptor{},
		fInfo: FlowInfo{
			FlowInfoBase: istats.NewFlowInfoBase(conf.Shadow.Mode),
			Name:         name,
			Addr:         conf.Addr,
			Topic:        conf.Topic,
			Group:        conf.Group,
			GroupType:    "consumer",
		},
	}
	pc.interceptors = append(pc.interceptors, pushConsumerDefaultInterceptor(pc), pushConsumerMDInterceptor(pc), pushConsumerShadowInterceptor(pc, conf.Shadow))

	_consumers.Store(name, pc)
	return pc
}

func (cc *PushConsumer) Close() error {
	err := cc.Shutdown()
	if err != nil {
		ilog.Warn("consumer close fail", ilog.Any("error", err.Error()))
		return err
	}
	return nil
}

func (cc *PushConsumer) WithInterceptor(fs ...primitive.Interceptor) *PushConsumer {
	cc.interceptors = append(cc.interceptors, fs...)
	return cc
}

func (cc *PushConsumer) Subscribe(topic string, f func(context.Context, *primitive.MessageExt) error) *PushConsumer {
	if _, ok := cc.subscribers[topic]; ok {
		ilog.Panic("duplicated subscribe", ilog.String("topic", topic))
	}
	fn := func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			err := f(ctx, msg)
			if err != nil {
				ilog.Error("consumer message", ilog.String("err", err.Error()), ilog.String("field", cc.name), ilog.Any("ext", msg))
				return consumer.ConsumeRetryLater, err
			}
		}

		return consumer.ConsumeSuccess, nil
	}
	cc.subscribers[topic] = fn
	return cc
}

func (cc *PushConsumer) Start() error {
	// 初始化 PushConsumer
	client, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(cc.Group),
		consumer.WithNameServer(cc.Addr),
		consumer.WithMaxReconsumeTimes(cc.Reconsume),
		consumer.WithInterceptor(cc.interceptors...),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: cc.AccessKey,
			SecretKey: cc.SecretKey,
		}),
	)
	cc.PushConsumer = client

	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: "",
	}
	if cc.ConsumerConfig.SubExpression != "*" {
		selector.Expression = cc.ConsumerConfig.SubExpression
	}

	for topic, fn := range cc.subscribers {
		if err := cc.PushConsumer.Subscribe(topic, selector, fn); err != nil {
			return err
		}
	}

	if err != nil || client == nil {
		ilog.Panic("create consumer",
			ilog.FieldName(cc.name),
			ilog.FieldExtMessage(cc.ConsumerConfig),
			ilog.Any("error", err),
		)
	}

	if cc.Enable {
		if err := client.Start(); err != nil {
			ilog.Panic("start consumer",
				ilog.FieldName(cc.name),
				ilog.FieldExtMessage(cc.ConsumerConfig),
				ilog.Any("error", err),
			)
			return err
		}
		// 在应用退出的时候，保证注销
		defers.Register(cc.Close)
	}

	return nil
}
