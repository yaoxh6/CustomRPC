package transport

import (
	"bufio"
	"fmt"
	"net"
)

type Package struct {
	ServiceName string
	Data []byte
}

type Option func(t Transport) error

type Transport interface {
	Setup(network string, address string) error
	Recv() (*Package, error)
	Send(pak *Package, opts ...Option) error
}

type SimpleTransport struct {
	Client net.Listener
	Conn net.Conn
}

func (s *SimpleTransport) Setup (network string, address string) error {
	var err error
	s.Client, err = net.Listen(network, address)
	if err != nil {
		return err
	}
	return nil
}

func (s *SimpleTransport) Recv() (*Package, error)  {
	var err error
	s.Conn, err = s.Client.Accept() // 监听客户端的连接请求
	if err != nil {
		fmt.Println("Accept() failed, err: ", err)
		return nil, err
	}
	reader := bufio.NewReader(s.Conn)
	var buf [128]byte
	n, err := reader.Read(buf[:]) // 读取数据
	if err != nil {
		fmt.Println("read from client failed, err: ", err)
		return nil, err
	}
	recvStr := buf[:n]
	//fmt.Println("send client data：", recvStr)
	//conn.Write([]byte(recvStr))
	res := &Package {
		ServiceName: "",
		Data:        recvStr,
	}
	return res, nil
}

func (s *SimpleTransport) Send (pak *Package, opts ...Option) error {
	_, err := s.Conn.Write(pak.Data)
	if err != nil {
		return err
	}
	return nil
}