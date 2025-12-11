package church

//go:generate protoc -I./server/proto --go_opt=module=github.com/abisalde/go-showcase --go_out=. --go-grpc_opt=module=github.com/abisalde/go-showcase --go-grpc_out=. ./server/proto/church/church.proto
