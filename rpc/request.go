package rpc

import (
	"context"
	"github.com/yaoxh6/CustomRPC/rpc/transport"
	"time"
	"github.com/pkg/errors"
)

const (
	DefaultCustomRequestTimeout = 1 * time.Second
)

type RequestData struct {
	RemoteFuncName string `json:"remote_func_name"`
	RequestId string `json:"request_id"`
	Params []interface{} `json:"params"`
}

type CustomRequest struct {
	requestId   string
	data        []byte
	suspendLock chan transport.Package
	options     RequestOptions
}

type CustomRespond struct {
	pak *transport.Package
	err error
}

type RequestOption func(o *RequestOptions)
type RequestOptions struct {
	Timeout time.Duration
}

func Timeout(timeout time.Duration) RequestOption {
	return func(o *RequestOptions) {
		o.Timeout = timeout
		if o.Timeout < 0 {
			// In this case, timeout would be immediately fired
			o.Timeout = 0
		}
	}
}

func NewCustomRequest(requestId string, data []byte, opts ...RequestOption) *CustomRequest {
	h := &CustomRequest{
		requestId:   requestId,
		data:        data,
		suspendLock: make(chan transport.Package),
	}
	h.options = RequestOptions{Timeout: DefaultCustomRequestTimeout}
	for _, o := range opts {
		o(&h.options)
	}
	return h
}

func (r *CustomRequest) WaitComplete(ctx context.Context) *CustomRespond {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, r.options.Timeout)
	defer cancel()

	select {
	case <-ctxWithTimeout.Done():
		return NewCustomRespond(nil, errors.Wrap(ctxWithTimeout.Err(), "request timeout"))
	case resp := <-r.suspendLock:
		return NewCustomRespond(&resp, nil)
	}
}

func (r *CustomRequest) ResumeExecution(pak *transport.Package) {
	r.suspendLock <- *pak
}

func NewCustomRespond(pak *transport.Package, err error) *CustomRespond {
	return &CustomRespond{
		pak: pak,
		err: err,
	}
}

func (r *CustomRespond) GetError() error {
	return r.err
}

func (r *CustomRespond) GetData() *transport.Package {
	return r.pak
}
