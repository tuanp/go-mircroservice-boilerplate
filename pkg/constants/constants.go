package constants

import "time"

const (
	GrpcPort   = "GRPC_PORT"
	HttpPort   = "HTTP_PORT"
	ConfigPath = "CONFIG_PATH"
	JaegerHost = "JAEGER_HOST"
	JaegerPort = "JAEGER_PORT"
	RedisAddr  = "REDIS_ADDR"
	Yaml       = "yaml"
	Json       = "json"

	GRPC     = "GRPC"
	METHOD   = "METHOD"
	NAME     = "NAME"
	METADATA = "METADATA"
	REQUEST  = "REQUEST"
	REPLY    = "REPLY"
	TIME     = "TIME"
)

const (
	MaxHeaderBytes       = 1 << 20
	StackSize            = 1 << 10 // 1 KB
	BodyLimit            = "2M"
	ReadTimeout          = 15 * time.Second
	WriteTimeout         = 15 * time.Second
	GzipLevel            = 5
	WaitShotDownDuration = 3 * time.Second
	EnvDev               = "dev"
	EnvTest              = "test"
	EnvProd              = "prod"
)
