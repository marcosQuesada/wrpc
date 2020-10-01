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
	"github.com/marcosQuesada/wrpc/pkg/bufconn"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/route_guide/routeguide"
	"log"
	"net"
	"time"
)

// xserverCmd represents the xserver command
var xserverCmd = &cobra.Command{
	Use:   "xserver",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("xserver called")

		buffSpace := 10
		lis := bufconn.Listen(buffSpace)

		go func(){
			time.Sleep(time.Millisecond * 10)

			var opts = []grpc.DialOption{
				grpc.WithInsecure(),
				grpc.WithBlock(),
				grpc.WithContextDialer(func(ctx context.Context, host string) (net.Conn, error) {
					return lis.Dial()
				}),
			}

			conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), opts...)
			if err != nil {
				log.Fatalf("fail to dial: %v", err)
			}
			defer func(){
				conn.Close()
				log.Printf("Close connection \n")
			}()

			client := pb.NewRouteGuideClient(conn)

			log.Printf("printFeature \n")
			// Looking for a valid feature
			printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})

			log.Printf("printFeature \n")
			// Feature missing.
			printFeature(client, &pb.Point{Latitude: 0, Longitude: 0})

			log.Printf("runRecordRoute \n")
			runRecordRoute(client)

			log.Printf("runRouteChat \n")
			runRouteChat(client);
			log.Printf("DONE \n")
		}()

		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)
		pb.RegisterRouteGuideServer(grpcServer, newServer())
		grpcServer.Serve(lis)
	},
}

func init() {
	rootCmd.AddCommand(xserverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// xserverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// xserverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
