custom:
	go build && \
	move protoc-gen-custom.exe ./protoc-gen-custom/
	./protoc/bin/protoc --plugin=protoc-gen-custom=protoc-gen-custom/protoc-gen-custom.exe --custom_out=./custom-out ./proto/*.proto
go:
	./protoc/bin/protoc --plugin=protoc-gen-go=protoc-gen-go/protoc-gen-go.exe --go_out=./go-out ./proto/*.proto