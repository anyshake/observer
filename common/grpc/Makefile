.PHONY: proto

proto:
	rm -rf *.pb.go
	protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:. *.proto