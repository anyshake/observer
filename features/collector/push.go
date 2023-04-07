package collector

import (
	"context"
	"encoding/base64"
	"encoding/json"

	pb "com.geophone.observer/common/grpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
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

	respb, err := proto.Marshal(res)
	if err != nil {
		options.Status.Failures++
		options.Status.Queued--
		options.OnErrorCallback(err)
		return
	}

	message := new(pb.ResponseMessage)
	err = proto.Unmarshal(respb, message)
	if err != nil {
		options.Status.Failures++
		options.Status.Queued--
		options.OnErrorCallback(err)
		return
	}

	respb, err = protojson.Marshal(message)
	if err != nil {
		options.Status.Failures++
		options.Status.Queued--
		options.OnErrorCallback(err)
		return
	}

	var response interface{}
	err = json.Unmarshal(respb, &response)
	if err != nil {
		options.Status.Failures++
		options.Status.Queued--
		options.OnErrorCallback(err)
		return
	}

	respb, err = base64.StdEncoding.DecodeString(
		response.(map[string]interface{})["Data"].(string),
	)
	if err != nil {
		options.Status.Failures++
		options.Status.Queued--
		options.OnErrorCallback(err)
	}

	err = json.Unmarshal(respb, &response)
	if err != nil {
		options.Status.Failures++
		options.OnErrorCallback(err)
		return
	} else {
		options.Status.Pushed++
		options.OnCompleteCallback(response)
	}

	options.Status.Queued--
}
