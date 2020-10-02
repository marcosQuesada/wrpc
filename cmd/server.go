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
	pb "google.golang.org/grpc/examples/route_guide/routeguide"
	"log"
	"net"
	"net/http"
	"sync"
)

var saddr string

var upgrader = websocket.Upgrader{} // use default options

// serverCmd represents the xserver command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")

		buffSpace := 10
		lis := &foo{listener: bufconn.Listen(buffSpace)}

		var opts = []grpc.ServerOption{
			grpc.StreamInterceptor(lis.streamInterceptor),
			grpc.UnaryInterceptor(lis.unaryInterceptor),
		}
		grpcServer := grpc.NewServer(opts...)
		pb.RegisterRouteGuideServer(grpcServer, newServer())
		go grpcServer.Serve(lis.listener)

		http.HandleFunc("/ws", lis.handler)
		log.Fatal(http.ListenAndServe(saddr, nil))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&saddr, "addr", "a", "localhost:8080", "Remote Host")
}

type foo struct {
	listener *bufconn.Listener
	mutex    sync.Mutex
}

func (f *foo) handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("On Handle \n")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	log.Printf("Upgraded \n")

	var conn net.Conn = ws.NewConn(c)
	inBound, outBound := net.Pipe()

	err = f.listener.Handle(outBound)
	if err != nil {
		log.Print("dial:", err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			_, data, err := c.ReadMessage() //@TODO: Must be conn.Read
			if err != nil {
				log.Println("Error ReadMessage:", err)
				_ = inBound.Close()
				_ = outBound.Close()

				break
			}

			_, err = inBound.Write(data)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}()

	for {
		rsp := make([]byte, 1024) //@TODO: HERE
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

func (f *foo) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Printf("Unary Interceptor begin \n")
	defer log.Printf("Unary Interceptor done \n")

	return handler(ctx, req)
}

func (f *foo) streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Printf("Stream Interceptor begin \n")
	defer log.Printf("Stream Interceptor done \n")

	spew.Dump(info)
	return handler(srv, ss)
}
