package rpc

import (
	"context"
	"fmt"
	log "github.com/hyahm/golog"
	"github.com/pkg/errors"
	"github.com/yaoxh6/CustomRPC/rpc/transport"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type ContextKey string

const (
	ContextRequestPackage = ContextKey("CONTEXT_REQUEST_PACKAGE")
	ContextCustomService    = ContextKey("CONTEXT_CUSTOM_SERVICE")
)

type ServiceHandler interface {
	Name() string
	HandleRPC(context.Context, string) ([]byte, error)
}

type CustomService struct {
	requestId      int64
	ctx            context.Context
	cancel         context.CancelFunc
	trans          transport.Transport
	serviceHandler ServiceHandler
	suspendMap     sync.Map // map[string]*CustomRequest
	d				Codec
}

func (h *CustomService) NewRequestId() string {
	requestId := atomic.AddInt64(&h.requestId, 1)
	return fmt.Sprintf("#%d", requestId)
}

func (h *CustomService) initService() error {
	h.suspendMap = sync.Map{}
	return nil
}

func (h *CustomService) handleSuspendedRequest(requestId string, pak *transport.Package) error {
	if raw, hasRequest := h.suspendMap.LoadAndDelete(requestId); hasRequest {
		request := raw.(*CustomRequest)
		request.ResumeExecution(pak)
		return nil
	}
	return errors.Errorf("skip non-existing request: %s", requestId)
}

func (h *CustomService) handleRPC(rpcName string, pak *transport.Package) ([]byte, error) {
	ctxNew := context.WithValue(h.ctx, ContextRequestPackage, pak)
	return h.serviceHandler.HandleRPC(ctxNew, rpcName)
}

func (h *CustomService) internalHandle(pak *transport.Package) {
	//err := d.Load(pak.Data)
	//if err != nil {
	//	log.Errorf("load failed. ctx:%+v, err:%+v, data:%v", h.ctx, err, pak.Data)
	//	return
	//}

	var rpcName string
	err := h.d.Decode(pak.Data, &rpcName)
	if err != nil {
		log.Errorf("unmarshal failed. ctx:%+v, err:%+v", h.ctx, err)
		//data := d.Peek()
		log.Errorf("unmarshal failed data:%v", pak.Data)
		//log.Infof("unmarshal failed. ", data)
		return
	}

	if len(rpcName) > 0 && rpcName[0] == byte('#') {
		err = h.handleSuspendedRequest(rpcName, pak)
		if err != nil {
			log.Debugf("handle suspended request failed:%+v", err)
		}
	} else {
		rspBin, err := h.handleRPC(rpcName, pak)
		if err == nil && len(rspBin) > 0 {
			var sendPak = *pak
			sendPak.Data = rspBin
			err = h.Send(&sendPak)
		}
		if err != nil {
			log.Debugf("handle rpc request failed:%+v", err)
		}
	}
}

func (h *CustomService) internalRecv() (int, error) {
	pak, err := h.trans.Recv()
	if err != nil {
		return 0, errors.Wrap(err, "Recv failed")
	}
	if pak == nil {
		return 0, nil
	}

	go func(pak *transport.Package) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("panic message:%v", r)

				var buf [1024 * 1024]byte
				n := runtime.Stack(buf[:], false)
				log.Error(string(buf[:n]))
				panic(r)
			}
		}()

		h.internalHandle(pak)
	}(pak)
	return 1, nil
}

func (h *CustomService) listenAndServe() error {
	for {
		select {
		case <-h.ctx.Done():
			return h.ctx.Err()
		default:
		}

		recvCnt, err := h.internalRecv()
		if err != nil {
			log.Errorf("internalRecv err: %+v", err)
		}

		if recvCnt == 0 {
			time.Sleep(time.Millisecond)
		}
	}
}

func (h *CustomService) waitForEventReady() error {
	// TODO: This is definitely Poor but Effective, expected to be replaced by Signal Chan
	time.Sleep(1 * time.Second)
	return nil
}

func (h *CustomService) StartEventLoop() error {
	go func() {
		err := h.listenAndServe()
		if err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Fatalf("serve crash: %+v", err)
			}
		}
	}()
	return h.waitForEventReady()
}

func (h *CustomService) Close(ch chan struct{}) error {
	if ch == nil {
		ch = make(chan struct{}, 1)
	}

	pid := os.Getpid()
	log.Warnf("process:%d, service:HiveService, closing ...", pid)
	h.cancel()

	log.Warnf("process:%d, service:HiveService, closed", pid)
	ch <- struct{}{}
	return nil
}

func (h *CustomService) Serve() error {
	defer log.Sync()

	err := h.StartEventLoop()
	if err != nil {
		log.Fatalf("serve crash: %+v", err)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)
	s := <-ch
	log.Warnf("recv signal:", s)

	c := make(chan struct{}, 1)
	return h.Close(c)
}

func (h *CustomService) HangOnRequest(ctx context.Context, requestId string, msg *CustomRequest, pak *transport.Package) *CustomRespond {
	h.suspendMap.Store(requestId, msg)
	defer h.suspendMap.Delete(requestId)

	err := h.Send(pak)
	if err != nil {
		var rsp CustomRespond
		rsp.pak = nil
		rsp.err = errors.Wrap(err, "send failed")
		return &rsp
	}

	rsp := msg.WaitComplete(ctx)
	return rsp
}

func (h *CustomService) Send(pak *transport.Package) error {
	return h.trans.Send(pak)
}
