package ecode

// 统一 Err Kind
const (
	// ErrKindUnmarshalConfigErr ...
	ErrKindUnmarshalConfigErr = "unmarshal config err"
	// ErrKindRegisterErr ...
	ErrKindRegisterErr = "register err"
	// ErrKindUriErr ...
	ErrKindUriErr = "uri err"
	// ErrKindRequestErr ...
	ErrKindRequestErr = "request err"
	// ErrKindFlagErr ...
	ErrKindFlagErr = "flag err"
	// ErrKindListenErr ...
	ErrKindListenErr = "listen err"
	// ErrKindAny ...
	ErrKindAny = "any"
)

// 统一 Msg信息，模块内信息唯一
const (
	// MsgRegisterParseConfigErr1 ...
	MsgRegisterParseConfigErr1 = "register parse config err1"
	// MsgRegisterParseConfigErr2 ...
	MsgRegisterParseConfigErr2 = "register parse config err2"
	// MsgRegisterETCDErr ...
	MsgRegisterETCDErr = "register service err"
	// MsgRegisterETCDOk ...
	MsgRegisterETCDOk = "register service ok"
	// MsgDeleteParseProviderUriErr ...
	MsgDeleteParseProviderUriErr = "delete parse provider uri err"
	// MsgDeleteParseConfiguratorsUriErr ...
	MsgDeleteParseConfiguratorsUriErr = "delete parse configurators uri err"
	// MsgUpdateParseProviderUriErr ...
	MsgUpdateParseProviderUriErr = "update parse provider uri err"
	// MsgUpdateParseConfiguratorsUriErr ...
	MsgUpdateParseConfiguratorsUriErr = "update parse configurators uri err"
	// MsgUpdateParseConfigErr ...
	MsgUpdateParseConfigErr = "update parse config err"
	// MsgUpdateResolverOk ...
	MsgUpdateResolverOk = "update resolver ok"
	// MsgWatchRequestErr ...
	MsgWatchRequestErr = "watch request err"
	// MsgDeregisterETCDOk ...
	MsgDeregisterETCDOk = "deregister etcd service ok"
	// MsgDeregisterETCDErr ...
	MsgDeregisterETCDErr = "deregister etcd service err"
	// MsgRegistryResolverOk ...
	MsgRegistryResolverOk = "resolver build ok"
	// MsgRegistryResolverNow ...
	MsgRegistryResolverNow = "resolver now"
	// MsgRegistryResolverClose ...
	MsgRegistryResolverClose = "resolver close"

	// config msg
	MsgConfigParseFlagPanic = "parse flag panic"
	// MsgConfigLoadFromRemoteDataSourcePanic ...
	MsgConfigLoadFromRemoteDataSourcePanic = "load from remote datasource panic"
	// MsgConfigLoadFromRemoteDataSourceOK ...
	MsgConfigLoadFromRemoteDataSourceOK = "load from remote datasource ok"
	// MsgConfigLoadFromFilePanic ...
	MsgConfigLoadFromFilePanic = "load from file panic"
	// MsgConfigLoadFromFileOK ...
	MsgConfigLoadFromFileOK = "load from file ok"

	// app msg
	MsgAppStartServerOk = "server start"
	// MsgAppStartGovernorOk ...
	MsgAppStartGovernorOk = "governor start"

	// proc msg
	MsgProcSetPanic = "set max procs panic"
	// MsgProcSetOk ...
	MsgProcSetOk = "set max procs ok"

	// grpc server
	MsgGrpcServerNewErr = "new grpc server err"
	// MsgGrpcServerRecover ...
	MsgGrpcServerRecover = "grpc server recover"

	// client mysql
	MsgClientMysqlOpenStart = "client mysql start"
	// MsgClientMysqlOpenPanic ...
	MsgClientMysqlOpenPanic = "mysql open panic"
	// MsgClientMysqlPingPanic ...
	MsgClientMysqlPingPanic = "mysql ping panic"
)

// 统一模块信息
const (
	// ModConfig ...
	ModConfig = "config"
	// ModApp ...
	ModApp = "app"
	// ModProc ...
	ModProc = "proc"
	// ModGrpcServer ...
	ModGrpcServer = "server.grpc"
	// ModRegistryETCD ...
	ModRegistryETCD = "registry.etcd"
	// ModClientETCD ...
	ModClientETCD = "client.etcd"
	// ModClientGrpc ...
	ModClientGrpc = "client.grpc"
	// ModClientMySQL ...
	ModClientMySQL = "client.mysql"
	// ModXcronETCD ...
	ModXcronETCD = "xcron.etcd"
)
