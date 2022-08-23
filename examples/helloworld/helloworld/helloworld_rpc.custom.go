// Code generated by protoc-gen-hive. DO NOT EDIT.
// source: examples/greeter/helloworld.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	errors "github.com/pkg/errors"
	math "math"
)

import (
	context "context"
	rpc "github.com/yaoxh6/CustomRPC/rpc"
	client "github.com/yaoxh6/CustomRPC/rpc/client"
	//transport "github.com/yaoxh6/CustomRPC/rpc/transport"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = errors.Errorf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.

// Message definition

type HelloRequest struct {
	Name string `lua:"name"`
}

type HelloReply struct {
	Message string `lua:"message"`
}

// Client API for Greeter service

type GreeterClientProxy interface {
	// Sends a greeting
	SayHello(ctx context.Context, helloRequest *HelloRequest, opts ...client.Option) (*HelloReply, error)
}

type greeterClientProxyImpl struct {
	c client.Client
}

func NewGreeterClientProxy(h *rpc.CustomService) GreeterClientProxy {
	return &greeterClientProxyImpl{
		c: client.NewCustomClient(h, "Greeter"),
	}
}

func (c *greeterClientProxyImpl) SayHello(ctx context.Context, helloRequest *HelloRequest, opts ...client.Option) (*HelloReply, error) {
	//var __d__ archiver.Decoder
	//var __waste__ string // Skip suspend session id
	//var __resp__ []byte
	var err error
	helloReply_ := new(HelloReply)

	//if helloRequest == nil {
	//	return nil, errors.New("parameter helloRequest is nil")
	//}
	//__reqId__, __req__, err := c.c.NewRequest(ctx, "SayHello", true, []interface{}{helloRequest.Name}, opts...)
	//if err != nil {
	//	goto LastReturn
	//}
	//
	//__resp__, err = c.c.CallFunction(ctx, "SayHello", __reqId__, __req__, opts...)
	//if err != nil {
	//	goto LastReturn
	//}
	//__d__.Load(__resp__)
	//_ = __d__.Unmarshal(&__waste__)
	//
	//err = __d__.Unmarshal(&helloReply_.Message)
	//if err != nil {
	//	goto LastReturn
	//}

//LastReturn:
	return helloReply_, err
}

// Server API for Greeter service

const GreeterServer_ServiceName = "Greeter"

type GreeterServer interface {
	// Sends a greeting
	SayHello(ctx context.Context, helloRequest *HelloRequest) (*HelloReply, error)
}

func RegisterGreeterServer(h rpc.Service, svr GreeterServer) error {
	return h.Register(nil, &greeterHandler{svr})
}

type greeterHandler struct {
	GreeterServer
}

func (h *greeterHandler) Name() string {
	return "Greeter"
}

func GreeterServer_Endpoints() []string {
	return []string{
		"SayHello",
	}
}

func (h *greeterHandler) ServiceMeta(rpcName string, metaType int) ([]string, []interface{}, error) {
	switch metaType {
	case 1:
		switch rpcName {
		case "SayHello":
			return []string{
					"*callback*", // callback
					"HelloRequest.Name",
				}, []interface{}{
					"", // callback
					"",
				}, nil
		}
	case 2:
		switch rpcName {
		case "SayHello":
			return []string{
					"HelloReply.Message",
				}, []interface{}{
					"",
				}, nil
		}
	}
	return nil, nil, fmt.Errorf("unknown meta type for '%s': %d", rpcName, metaType)
}

func (h *greeterHandler) HandleRPC(ctx context.Context, rpcName string) ([]byte, error) {
	//var err error
	switch rpcName {
	case "SayHello":
		//var helloRequest HelloRequest
		//var sessionId string
		//err = inData.Unmarshal(&sessionId)
		//if err != nil {
		//	reqFields, reqMeta, _ := h.ServiceMeta(rpcName, 1)
		//	rpc.PrintUnmatchedMeta(reqFields, reqMeta, inData.Peek())
		//	return nil, fmt.Errorf(`rpc failed in [%s]: %s`, rpcName, err.Error())
		//}
		//
		//err = inData.Unmarshal(&helloRequest.Name)
		//if err != nil {
		//	reqFields, reqMeta, _ := h.ServiceMeta(rpcName, 1)
		//	rpc.PrintUnmatchedMeta(reqFields, reqMeta, inData.Peek())
		//	return nil, fmt.Errorf(`rpc failed in [%s]: %s`, rpcName, err.Error())
		//}
		//
		//helloReply_, err := h.SayHello(ctx, &helloRequest)
		//if err != nil {
		//	return nil, fmt.Errorf(`rpc failed in [%s]: %s`, rpcName, err.Error())
		//}
		//
		//var e archiver.Encoder
		//e.Init()
		//
		//err = e.Marshal(sessionId)
		//if err != nil {
		//	return nil, fmt.Errorf(`rpc failed in [%s]: %s`, rpcName, err.Error())
		//}
		//
		//err = e.Marshal(helloReply_.Message)
		//if err != nil {
		//	return nil, fmt.Errorf(`rpc failed in [%s]: %s`, rpcName, err.Error())
		//}
		//
		//return e.Save()
	}
	return nil, fmt.Errorf("unknown rpc name: %s", rpcName)
}

func (h *greeterHandler) SayHello(ctx context.Context, helloRequest *HelloRequest) (*HelloReply, error) {
	return h.GreeterServer.SayHello(ctx, helloRequest)
}
