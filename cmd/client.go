package cmd

import (
	"context"
	"fmt"
	pb "github.com/marcosQuesada/wrpc/pkg/route_guide"
	"github.com/marcosQuesada/wrpc/pkg/route_guide/routeguide"
	"github.com/marcosQuesada/wrpc/pkg/ws"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/url"

	"github.com/spf13/cobra"
)

var host string

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "gRPC client in top of websockets",
	Long:  `gRPC client in top of websockets`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("client called")

		var opts = []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithContextDialer(func(ctx context.Context, host string) (net.Conn, error) {
				uri := url.URL{Scheme: "ws", Host: host, Path: "/ws"}
				return ws.NewClient(uri)
			}),
		}

		conn, err := grpc.Dial(host, opts...)
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()

		client := routeguide.NewRouteGuideClient(conn)

		// Looking for a valid feature
		pb.PrintFeature(client, &routeguide.Point{Latitude: 409146138, Longitude: -746188906})

		// Feature missing.
		pb.PrintFeature(client, &routeguide.Point{Latitude: 0, Longitude: 0})

		// Looking for features between 40, -75 and 42, -73.
		pb.PrintFeatures(client, &routeguide.Rectangle{
			Lo: &routeguide.Point{Latitude: 400000000, Longitude: -750000000},
			Hi: &routeguide.Point{Latitude: 420000000, Longitude: -730000000},
		})

		// RecordRoute
		pb.RunRecordRoute(client)

		// RouteChat
		pb.RunRouteChat(client)
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVarP(&host, "host", "s", "localhost:8080", "Remote Host")
}
