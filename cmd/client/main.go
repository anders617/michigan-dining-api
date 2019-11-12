package main

import (
	"context"
	"flag"
	"io"
	"os"
	"time"

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
	res, err := c.GetHearts(context.Background(), &pb.HeartsRequest{Keys: []string{"test", "test2"}})
	if err != nil {
		glog.Fatalf("GetHearts err: %s", err)
	}
	glog.Infof("%v", res)
	stream, err := c.StreamHearts(context.Background(), &pb.HeartsRequest{Keys: []string{"test"}})
	if err != nil {
		glog.Fatalf("Stream error: %s", err)
	}
	go func() {
		for i := 0; i < 10; i++ {
			newCounts, err := c.AddHeart(context.Background(), &pb.HeartsRequest{Keys: []string{"test", "test2"}})
			if err != nil {
				glog.Fatalf("Failed to add heart: %s", err)
			}
			for _, count := range newCounts.Counts {
				glog.Infof("New heart count key: %s count: %d", count.Key, count.Count)
			}
			time.Sleep(time.Second)
		}
		glog.Infof("All done.")
		os.Exit(0)
	}()
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil || len(reply.Counts) == 0 {
			glog.Fatalf("Error receiving: %s", err)
		}
		for _, count := range reply.Counts {
			glog.Infof("HeartCount key: %s count: %d", count.Key, count.Count)
		}
	}
}
