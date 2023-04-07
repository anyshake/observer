package collector

import (
	"context"
	"encoding/json"

	pb "com.geophone.observer/common/grpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func PushMessage(conn *grpc.ClientConn, grpc pb.CollectorClient, options *CollectorOptions) {
	if !options.Enable {
		return
	}

	req, err := json.Marshal(options.Message)
	if err != nil {
		options.OnErrorCallback(err)
		return
	}

	options.Status.Queued++
	res, err := grpc.ToDatabase(context.Background(),
		&pb.RequestMessage{
			Data: req,
		})
	if err != nil {
		options.Status.Failures++
		options.Status.Queued--
		options.OnErrorCallback(err)
		return
	}

	resb, err := proto.Marshal(res)
	if err != nil {
		options.Status.Failures++
		options.Status.Queued--
		options.OnErrorCallback(err)
		return
	}

	for i, v := range resb {
		if v == '\x1d' || v == '\x1e' || v == '\x1f' ||
			v == '\x00' || v == '\x0f' {
			resb = append(resb[:i], resb[i+1:]...)
		}
	}

	var response interface{}
	err = json.Unmarshal(resb, &response)
	if err != nil {
		options.Status.Failures++
		options.OnErrorCallback(err)
	} else {
		options.OnCompleteCallback(response)
	}

	options.Status.Queued--
}
