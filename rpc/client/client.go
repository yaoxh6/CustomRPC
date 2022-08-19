package client
import (
	"context"
	"github.com/yaoxh6/CustomRPC/rpc"
	"github.com/yaoxh6/CustomRPC/rpc/transport"
)
type Client interface {
	NewRequest(ctx context.Context, remoteFuncName string, isSyncRequest bool, params []interface{}, opts ...Option) (string, []byte, error)
	CallFunction(ctx context.Context, remoteFuncName string, requestId string, packedRequest []byte, opts ...Option) ([]byte, error)
	CallProcedure(ctx context.Context, remoteFuncName string, packedRequest []byte, opts ...Option) error
}

type CustomClient struct {
	c *rpc.CustomService
	options Options
	serviceName string
}

func NewCustomClient(h *rpc.CustomService, serviceName string, opts ...Option) *CustomClient {
	options := NewOptions()
	for _, o := range opts {
		o(&options)
	}
	return &CustomClient{
		c:           h,
		serviceName: serviceName,
	}
}
//func generateRequestData(remoteFuncName string, requestId string, params []interface{}, options Options) ([]byte, error) {
//	var err error
//	var enc archiver.Encoder
//	enc.Init()
//
//	err = enc.Marshal(remoteFuncName)
//	if err != nil {
//		return nil, errors.Wrapf(err, "error encoding remote function name: %s", remoteFuncName)
//	}
//	for i, param := range params {
//		if len(requestId) > 0 && options.CallOptions.ReqIdShift == i {
//			err = enc.Marshal(requestId)
//			if err != nil {
//				return nil, errors.Wrapf(err, "error encoding request id @%d", options.CallOptions.ReqIdShift)
//			}
//		}
//
//		err = enc.Marshal(param)
//		if err != nil {
//			return nil, errors.Wrapf(err, "error encoding request param #%d: %v", i, param)
//		}
//	}
//	if len(requestId) > 0 && options.CallOptions.ReqIdShift >= len(params) {
//		err = enc.Marshal(requestId)
//		if err != nil {
//			return nil, errors.Wrapf(err, "error encoding request id @%d", options.CallOptions.ReqIdShift)
//		}
//	}
//
//	return enc.Save()
//}

func generateRequestData(remoteFuncName string, requestId string, params []interface{}, options Options) ([]byte, error) {
	code := rpc.JsonCodec{}
	data := rpc.RequestData{
		RemoteFuncName: remoteFuncName,
		RequestId: requestId,
		Params: params,
	}
	return code.Encode(data)
}

func (h *CustomClient) NewRequest(ctx context.Context, remoteFuncName string, isSyncRequest bool, params []interface{}, opts ...Option) (string, []byte, error) {
	var err error
	options := h.options
	for _, o := range opts {
		o(&options)
	}

	var requestId string
	if isSyncRequest {
		requestId = h.c.NewRequestId()
	}

	data, err := generateRequestData(remoteFuncName, requestId, params, options)
	return requestId, data, err
}

func (h *CustomClient) CallFunction(ctx context.Context, remoteFuncName string, requestId string, packedRequest []byte, opts ...Option) ([]byte, error) {
	options := h.options

	pak := transport.Package{
		ServiceName: h.serviceName,
		Data: packedRequest,
	}
	for _, o := range opts {
		o(&options)
	}

	req := rpc.NewCustomRequest(requestId, packedRequest)
	resp := h.c.HangOnRequest(ctx, requestId, req, &pak)

	if resp.GetError() != nil {
		return nil, resp.GetError()
	}
	return resp.GetData().Data, nil
}

func (h *CustomClient) CallProcedure(ctx context.Context, remoteFuncName string, packedRequest []byte, opts ...Option) error {
	options := h.options

	pak := transport.Package{
		ServiceName: h.serviceName,
		Data: packedRequest,
	}
	for _, o := range opts {
		o(&options)
	}
	return h.c.Send(&pak)
}
