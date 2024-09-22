gen:
	protoc --proto_path=proto --go_out=./pb --go-grpc_out=./pb proto/*.proto

run-server:
	go run cmd/server/main.go

.PHONY:
	gen run-server