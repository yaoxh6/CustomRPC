module github.com/yaoxh6/CustomRPC/examples/helloworld

go 1.18

require (
	github.com/golang/protobuf v1.5.2
	github.com/hyahm/golog v0.0.0-20220530023613-47d2a9f33e27
	github.com/pkg/errors v0.9.1
	github.com/yaoxh6/CustomRPC/rpc v0.0.0-20220822081750-da79c3733135
	google.golang.org/grpc v1.48.0
)

require (
	github.com/fatih/color v1.13.0 // indirect
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/yaoxh6/CustomRPC/rpc => ../../rpc
