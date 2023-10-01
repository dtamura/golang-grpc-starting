
ref: https://zenn.dev/hsaki/books/golang-grpc-starting


go get -u google.golang.org/grpc

go mod init golang-grpc-starting
go get -u google.golang.org/grpc
go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

protoc -I=proto --go_out=./genproto/hello --go_opt=paths=source_relative \
	--go-grpc_out=./genproto/hello --go-grpc_opt=paths=source_relative \
	proto/hello.proto


go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest


grpcurl -plaintext -d '{"name": "alice"}' localhost:8080 myapp.GreetingService.Hello