package collector

import (
	"context"
	"encoding/json"

	pb "com.geophone.observer/helper/grpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func PushCollection(conn *grpc.ClientConn, grpc pb.CollectorClient, options *CollectorOptions) {
	d, err := json.Marshal(options.Message)
	if err != nil {
		options.OnErrorCallback(err)
		return
	}

	options.Status.Queued++
	r, err := grpc.ToDatabase(context.Background(),
		&pb.ClientMessage{
			Data: []byte(d),
		})
	if err != nil {
		options.Status.Queued--
		options.OnErrorCallback(err)
		return
	}

	rb, err := proto.Marshal(r)
	if err != nil {
		options.Status.Queued--
		options.OnErrorCallback(err)
		return
	}

	var response interface{}
	err = json.Unmarshal(rb, &response)
	if err != nil {
		options.OnErrorCallback(err)
	} else {
		options.OnCompleteCallback(response)
	}

	options.Status.Queued--
}
