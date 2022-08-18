package rpc

import (
	"fmt"
	"sync"
	"context"
	"sync/atomic"
)

type ServiceHandler interface {
	Name() string
	HandleRPC(context.Context, string) ([]byte, error)
}

type CustomService struct {
	requestId      int64
	ctx            context.Context
	cancel         context.CancelFunc
	serviceHandler ServiceHandler
	suspendMap     sync.Map // map[string]*HiveRequest
}

func (h *CustomService) NewRequestId() string {
	requestId := atomic.AddInt64(&h.requestId, 1)
	return fmt.Sprintf("#%d", requestId)
}