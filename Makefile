custom:
	cd protoc-tool && \
	go build && \
	move protoc-gen-custom.exe ./protoc-gen-custom/
	./protoc-tool/protoc/bin/protoc --plugin=protoc-gen-custom=protoc-tool/protoc-gen-custom/protoc-gen-custom.exe --custom_out=./examples/helloworld/helloworld ./examples/helloworld/helloworld/*.proto

go:
	./protoc-tool/protoc/bin/protoc --plugin=protoc-gen-go=protoc-tool/protoc-gen-go/protoc-gen-go.exe --go_out=./examples/helloworld/helloworld ./examples/helloworld/helloworld/*.proto
	./protoc-tool/protoc/bin/protoc --plugin=protoc-gen-go-grpc=cmd/protoc-gen-go-grpc/protoc-gen-go-grpc.exe --go-grpc_out=./examples/helloworld/helloworld ./examples/helloworld/helloworld/*.proto