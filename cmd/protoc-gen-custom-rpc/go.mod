module github.com/yaoxh6/CustomRPC/cmd/protoc-gen-custom-rpc

go 1.16

require (
	github.com/golang/protobuf v1.5.2
	google.golang.org/genproto v0.0.0-20220801145646-83ce21fca29f
	google.golang.org/protobuf v1.28.0
)

replace github.com/yaoxh6/CustomRPC/rpc => ../../rpc
