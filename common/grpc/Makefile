.PHONY: proto

proto:
	rm -rf *.pb.go
	protoc --go_out=. --go-grpc_out=. config.proto