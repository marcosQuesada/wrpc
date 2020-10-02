package cmd

import (
	"fmt"
	"github.com/marcosQuesada/wrpc/pkg/bufconn"
	pb "github.com/marcosQuesada/wrpc/pkg/route_guide"
	"github.com/marcosQuesada/wrpc/pkg/route_guide/routeguide"
	"github.com/marcosQuesada/wrpc/pkg/ws"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

var serverAddr string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "gRPC server in top of websocket transport",
	Long:  `gRPC server in top of websocket transport.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")

		listener := bufconn.Listen()
		srv := ws.NewServer(listener)

		var opts = []grpc.ServerOption{
			grpc.StreamInterceptor(listener.StreamInterceptor),
			grpc.UnaryInterceptor(listener.UnaryInterceptor),
		}

		grpcServer := grpc.NewServer(opts...)
		routeguide.RegisterRouteGuideService(grpcServer, pb.NewServer().Svc())
		go grpcServer.Serve(listener)

		http.HandleFunc("/ws", srv.Handler)
		log.Fatal(http.ListenAndServe(serverAddr, nil))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&serverAddr, "addr", "a", "localhost:8080", "Remote Host")
}
