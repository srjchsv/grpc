gen:
	@protoc  --go_out=. --go-grpc_out=. --proto_path=proto proto/*.proto
clean:
	rm pd/*.go
run:
	go run main.go

test:
	go test -cover -race ./...

server:
	go run cmd/server/main.go -port 8080

client:
	go run cmd/client/main.go -address localhost:8080
