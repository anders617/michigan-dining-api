package main

import (
	"context"
	"flag"

	pb "github.com/anders617/mdining-proto/proto/mdining"
	"github.com/golang/glog"
	"google.golang.org/grpc"
)

//
// A basic client implementation for testing purposes
//

func main() {
	address := flag.String("address", "michigan-dining-api.herokuapp.com:80", "The address of the mdining server to connect to.")
	flag.Parse()
	glog.Infof("Connecting...")
	conn, err := grpc.Dial(*address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		glog.Fatalf("Could not dial %s: %s", address, err)
	}
	glog.Infof("Connected")
	defer conn.Close()
	c := pb.NewMDiningClient(conn)
	diningHallsReply, err := c.GetDiningHalls(context.Background(), &pb.DiningHallsRequest{})
	if err != nil {
		glog.Fatalf("Could not call GetDiningHalls: %s", err)
	}
	glog.Infof("DiningHallsReply: %v", diningHallsReply)
}
