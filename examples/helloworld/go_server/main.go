/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"bufio"
	"fmt"
	"net"
)

//
//import (
//	"context"
//	"flag"
//	"fmt"
//	"log"
//	"net"
//
//	pb "github.com/yaoxh6/CustomRPC/examples/helloworld/helloworld"
//	"google.golang.org/grpc"
//)
//
//var (
//	port = flag.Int("port", 50051, "The server port")
//)
//
//// server is used to implement helloworld.GreeterServer.
//type server struct {
//	pb.UnimplementedGreeterServer
//}
//
//// SayHello implements helloworld.GreeterServer
//func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
//	log.Printf("Received: %v", in.GetName())
//	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
//}
//
//func main() {
//	flag.Parse()
//	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//	s := grpc.NewServer()
//	pb.RegisterGreeterCustomServer(s, &server{})
//	log.Printf("server listening at %v", lis.Addr())
//	if err := s.Serve(lis); err != nil {
//		log.Fatalf("failed to serve: %v", err)
//	}
//}

func process(conn net.Conn) {
	defer conn.Close() // 关闭连接
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			fmt.Println("read from client failed, err: ", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("收到Client端发来的数据：", recvStr)
		conn.Write([]byte(recvStr)) // 发送数据
	}
}

func main() {
	client, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		fmt.Printf("Error : %+v\n",err)
	}
	for {
		conn, err := client.Accept() // 监听客户端的连接请求
		if err != nil {
			fmt.Println("Accept() failed, err: ", err)
			continue
		}
		go process(conn) // 启动一个goroutine来处理客户端的连接请求
	}
}