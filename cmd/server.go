/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
	"github.com/marcosQuesada/wrpc/pkg/bufconn"
	"github.com/marcosQuesada/wrpc/pkg/ws"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	pb "github.com/marcosQuesada/wrpc/pkg/route_guide"
	"log"
	"net"
	"net/http"
)

var defaultBufSize  = 32 * 1024

var saddr string

var upgrader = websocket.Upgrader{} // use default options

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "gRPC server in top of websocket transport",
	Long: `gRPC server in top of websocket transport.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")

		listener := bufconn.Listen()
		srv := newServer(listener)

		var opts = []grpc.ServerOption{
			grpc.StreamInterceptor(srv.streamInterceptor),
			grpc.UnaryInterceptor(srv.unaryInterceptor),
		}

		grpcServer := grpc.NewServer(opts...)
		pb.RegisterRouteGuideService(grpcServer, pb.NewServer().Svc())
		go grpcServer.Serve(listener)

		http.HandleFunc("/ws", srv.handler)
		log.Fatal(http.ListenAndServe(saddr, nil))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&saddr, "addr", "a", "localhost:8080", "Remote Host")
}

type server struct {
	listener *bufconn.Listener
}

func newServer(l *bufconn.Listener) *server{
	return &server{
		listener: l,
	}
}

func (f *server) handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade error:", err)
		return
	}

	var conn net.Conn = ws.NewConn(c)
	inBound, outBound := net.Pipe()

	err = f.listener.Handle(outBound)
	if err != nil {
		log.Print("dial error:", err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			data := make([]byte, defaultBufSize)
			n, err := conn.Read(data)
			if err != nil {
				log.Println("Error ReadMessage:", err)
				_ = inBound.Close()
				_ = outBound.Close()

				break
			}

			_, err = inBound.Write(data[:n])
			if err != nil {
				log.Println("inbound write error:", err)
				break
			}
		}
	}()

	for {
		rsp := make([]byte, defaultBufSize)
		n, err := inBound.Read(rsp)
		if err != nil {
			log.Println("readAll:", err)
			break
		}

		_, err = conn.Write(rsp[:n])
		if err != nil {
			log.Println("piped write:", err)
			break
		}

	}
}

func (f *server) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Printf("Unary Interceptor begin \n")
	defer log.Printf("Unary Interceptor done \n")

	return handler(ctx, req)
}

func (f *server) streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Printf("Stream Interceptor begin \n")
	defer log.Printf("Stream Interceptor done \n")

	spew.Dump(info)
	return handler(srv, ss)
}
