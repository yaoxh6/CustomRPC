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
	"log"

	pb "github.com/yaoxh6/CustomRPC/examples/helloworld/helloworld"
	"github.com/yaoxh6/CustomRPC/rpc"
)

//var (
//	port = flag.Int("port", 8888, "The server port")
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

func main() {
	h, err := InitServerEnv()
	if err != nil {
		log.Fatal(err)
		return
	}
	pb.RegisterGreeterServer(h, &s2s{})
	if err := rpc.Serve(h); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
