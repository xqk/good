package grpc

import (
	"context"
	"github.com/xqk/good/pkg/ecode"
	"github.com/xqk/good/pkg/ilog"
	"time"

	"google.golang.org/grpc"
)

func newGRPCClient(config *Config) *grpc.ClientConn {
	var ctx = context.Background()
	var dialOptions = config.dialOptions
	logger := config.logger.With(
		ilog.FieldMod("client.grpc"),
		ilog.FieldAddr(config.Address),
	)
	// 默认配置使用block
	if config.Block {
		if config.DialTimeout > time.Duration(0) {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, config.DialTimeout)
			defer cancel()
		}

		dialOptions = append(dialOptions, grpc.WithBlock())
	}

	if config.KeepAlive != nil {
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(*config.KeepAlive))
	}

	dialOptions = append(dialOptions, grpc.WithBalancerName(config.BalancerName))

	cc, err := grpc.DialContext(ctx, config.Address, dialOptions...)

	if err != nil {
		if config.OnDialError == "panic" {
			logger.Panic("dial grpc server", ilog.FieldErrKind(ecode.ErrKindRequestErr), ilog.FieldErr(err))
		} else {
			logger.Error("dial grpc server", ilog.FieldErrKind(ecode.ErrKindRequestErr), ilog.FieldErr(err))
		}
	}
	logger.Info("start grpc client")
	return cc
}
