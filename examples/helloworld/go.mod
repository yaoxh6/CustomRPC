module github.com/yaoxh6/CustomRPC/examples/helloworld

go 1.18

require (
	github.com/golang/protobuf v1.5.2
	github.com/hyahm/golog v0.0.0-20220530023613-47d2a9f33e27
	github.com/pkg/errors v0.9.1
	github.com/yaoxh6/CustomRPC/rpc v0.0.0-20220822081750-da79c3733135
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/yaoxh6/CustomRPC/rpc => ../../rpc
