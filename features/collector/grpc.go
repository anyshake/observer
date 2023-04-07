package collector

import (
	"fmt"
	"time"

	pb "com.geophone.observer/common/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

func OpenGrpc(host string, port int, tls bool, enable bool) (*grpc.ClientConn, pb.CollectorClient, error) {
	if !enable {
		return nil, nil, nil
	}

	conn, err := func(host string, port int, tls bool) (*grpc.ClientConn, error) {
		if tls {
			return grpc.Dial(
				fmt.Sprintf("%s:%d", host, port),
				grpc.WithKeepaliveParams(keepalive.ClientParameters{
					Time:                10 * time.Second,
					Timeout:             5 * time.Second,
					PermitWithoutStream: true,
				}),
				grpc.WithTransportCredentials(
					credentials.NewClientTLSFromCert(nil, host),
				),
			)
		}

		return grpc.Dial(
			fmt.Sprintf("%s:%d", host, port),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                10 * time.Second,
				Timeout:             5 * time.Second,
				PermitWithoutStream: true,
			}),
			grpc.WithInsecure(),
		)
	}(host, port, tls)
	if err != nil {
		return nil, nil, err
	}

	grpc := pb.NewCollectorClient(conn)
	return conn, grpc, nil
}

func CloseGrpc(conn *grpc.ClientConn) error {
	err := conn.Close()
	if err != nil {
		return err
	}

	return nil
}
