module github.com/yaoxh6/CustomRPC/examples/helloworld

go 1.18

require (
	github.com/golang/protobuf v1.5.2
	github.com/yaoxh6/CustomRPC/rpc v0.0.0-20220822081750-da79c3733135
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.1
)

replace (
	github.com/yaoxh6/CustomRPC/rpc => ../../rpc
)
