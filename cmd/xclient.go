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
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/route_guide/routeguide"
	"log"
	"net"

	"github.com/spf13/cobra"
)

// xclientCmd represents the xclient command
var xclientCmd = &cobra.Command{
	Use:   "xclient",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("xclient called")

		var opts = []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithContextDialer(func(ctx context.Context, host string) (net.Conn, error) {

				return nil, nil
			}),
		}

		conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), opts...)
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()


		client := pb.NewRouteGuideClient(conn)

		// Looking for a valid feature
		printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})

		// Feature missing.
		printFeature(client, &pb.Point{Latitude: 0, Longitude: 0})

		// Looking for features between 40, -75 and 42, -73.
		printFeatures(client, &pb.Rectangle{
			Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
			Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
		})

		// RecordRoute
		runRecordRoute(client)

		// RouteChat
		runRouteChat(client)
	},
}

func init() {
	rootCmd.AddCommand(xclientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// xclientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// xclientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
