package client
import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/yaoxh6/CustomRPC/rpc"
	"github.com/yaoxh6/CustomRPC/rpc/transport"
)
type Client interface {
	NewRequest(ctx context.Context, remoteFuncName string, isSyncRequest bool, params []interface{}) (string, []byte, error)
	CallFunction(ctx context.Context, remoteFuncName string, requestId string, packedRequest []byte) ([]byte, error)
	CallProcedure(ctx context.Context, remoteFuncName string, packedRequest []byte) error
}

type CustomClient struct {
	c *rpc.CustomService
	serviceName string
}

func NewCustomClient(h *rpc.CustomService, serviceName string) *CustomClient {
	return &CustomClient{
		c:           h,
		serviceName: serviceName,
	}
}

func generateRequestData(remoteFuncName string, requestId string, params []interface{}, options Options) ([]byte, error) {
	var err error

	remoteFuncNameBytes, err := json.Marshal(remoteFuncName)
	if err != nil {
		return nil, errors.Wrapf(err, "error encoding remote function name: %s", remoteFuncName)
	}
	for i, param := range params {
		if len(requestId) > 0 && options.CallOptions.ReqIdShift == i {
			err = enc.Marshal(requestId)
			if err != nil {
				return nil, errors.Wrapf(err, "error encoding request id @%d", options.CallOptions.ReqIdShift)
			}
		}

		err = enc.Marshal(param)
		if err != nil {
			return nil, errors.Wrapf(err, "error encoding request param #%d: %v", i, param)
		}
	}
	if len(requestId) > 0 && options.CallOptions.ReqIdShift >= len(params) {
		err = enc.Marshal(requestId)
		if err != nil {
			return nil, errors.Wrapf(err, "error encoding request id @%d", options.CallOptions.ReqIdShift)
		}
	}

	return enc.Save()
}

func (h *CustomClient) NewRequest(ctx context.Context, remoteFuncName string, isSyncRequest bool, params []interface{}) (string, []byte, error) {
	var err error

	var requestId string
	if isSyncRequest {
		requestId = h.c.NewRequestId()
	}

	data, err := generateRequestData(remoteFuncName, requestId, params)
	return requestId, data, err
}

func (h *CustomClient) CallFunction(ctx context.Context, remoteFuncName string, requestId string, packedRequest []byte) ([]byte, error) {
	pak := transport.Package{
		ServiceName: h.serviceName,
		Data: packedRequest,
	}

	req := rpc.NewCustomRequest(requestId, packedRequest)
	resp := h.c.HangOnRequest(ctx, requestId, req, &pak)

	if resp.GetError() != nil {
		return nil, resp.GetError()
	}
	return resp.GetData().Data, nil
}

func (h *CustomClient) CallProcedure(ctx context.Context, remoteFuncName string, packedRequest []byte) error {
	pak := transport.Package{
		Remote: transport.RemoteInfo{
			Type:    h.serviceType,
			Service: h.serviceName,
		},
		Data: packedRequest,
	}

	return h.c.Send(&pak)
}
