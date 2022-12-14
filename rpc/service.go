package rpc

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	log "github.com/hyahm/golog"
	"github.com/pkg/errors"
	"github.com/yaoxh6/CustomRPC/rpc/transport"
)

type ContextKey string

const (
	ContextRequestPackage = ContextKey("CONTEXT_REQUEST_PACKAGE")
	ContextCustomService  = ContextKey("CONTEXT_CUSTOM_SERVICE")
)

type ServiceHandler interface {
	Name() string
	HandleRPC(context.Context, string, Codec, *transport.Package) ([]byte, error)
}

type CustomService struct {
	requestId      int64
	ctx            context.Context
	cancel         context.CancelFunc
	trans          transport.Transport
	serviceHandler ServiceHandler
	d              Codec
}

func (h *CustomService) Register(serviceDesc interface{}, serviceImpl interface{}) error {
	sh, ok := serviceImpl.(ServiceHandler)
	if !ok {
		return errors.Errorf("unsupported serviceImpl type: %+v", serviceImpl)
	}
	if h.serviceHandler != nil {
		return errors.Errorf("serviceHandler already registered")
	}

	h.serviceHandler = sh
	return nil
}

func (h *CustomService) NewRequestId() string {
	requestId := atomic.AddInt64(&h.requestId, 1)
	return fmt.Sprintf("#%d", requestId)
}

func (h *CustomService) handleRPC(rpcName string, pak *transport.Package) ([]byte, error) {
	ctxNew := context.WithValue(h.ctx, ContextRequestPackage, pak)
	return h.serviceHandler.HandleRPC(ctxNew, rpcName, h.d, pak)
}

func (h *CustomService) internalHandle(pak *transport.Package) {
	var param []interface{}
	var rpcName string
	err := h.d.Decode(pak.Data, &param)
	rpcName = param[0].(string)
	if err != nil {
		log.Errorf("unmarshal failed. ctx:%+v, err:%+v", h.ctx, err)
		log.Errorf("unmarshal failed data:%v", pak.Data)
		return
	}

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
			return err
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
				log.Errorf("serve crash: %+v", err)
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
	log.Warnf("process:%d, service:CustomService, closing ...", pid)
	h.cancel()

	log.Warnf("process:%d, service:CustomService, closed", pid)
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

func (h *CustomService) Send(pak *transport.Package) error {
	return h.trans.Send(pak)
}
