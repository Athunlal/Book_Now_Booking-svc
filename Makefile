proto:
	protoc --go_out=. --go-grpc_out=. pkg/pb/booking.proto

protoUser:
	protoc --go_out=. --go-grpc_out=. pkg/usermodule/pb/user.proto

run:
	go run cmd/api/main.go


	