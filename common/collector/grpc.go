package collector

import (
	"fmt"

	pb "com.geophone.observer/helper/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func OpenGrpc(host string, port int, tls bool) (*grpc.ClientConn, pb.CollectorClient, error) {
	conn, err := func(host string, port int, tls bool) (*grpc.ClientConn, error) {
		if tls {
			return grpc.Dial(
				fmt.Sprintf("%s:%d", host, port),
				grpc.WithTransportCredentials(
					credentials.NewClientTLSFromCert(nil, host),
				),
			)
		}

		return grpc.Dial(
			fmt.Sprintf("%s:%d", host, port),
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
