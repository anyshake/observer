package app

import (
	pb "com.geophone.observer/common/grpc"
	"com.geophone.observer/features/collector"
	"github.com/go-redis/redis/v8"
)

type ServerOptions struct {
	Gzip      int
	Cors      bool
	Version   string
	ApiPrefix string
	WebPrefix string
	ConnRedis *redis.Client
	ConnGRPC  *pb.CollectorClient
	Message   *collector.Message
	Status    *collector.Status
}
