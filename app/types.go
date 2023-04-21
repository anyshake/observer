package app

import (
	"database/sql"

	pb "com.geophone.observer/common/grpc"
	"com.geophone.observer/features/collector"
)

type ServerOptions struct {
	Gzip         int
	Cors         bool
	Version      string
	ApiPrefix    string
	WebPrefix    string
	ConnPostgres *sql.DB
	ConnGRPC     *pb.CollectorClient
	Message      *collector.Message
	Status       *collector.Status
}
