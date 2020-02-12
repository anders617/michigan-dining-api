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
	"google.golang.org/grpc/credentials"
)

//
// A basic client implementation for testing purposes
//

func main() {
	address := flag.String("address", "michigan-dining-api.tendiesti.me:80", "The address of the mdining server to connect to.")
	useCredentials := flag.Bool("use_credentials", false, "Whether to use tls credentials or not when connecting to the server.")
	flag.Parse()
	glog.Infof("Connecting...")
	credentialOpt := grpc.WithInsecure()
	if *useCredentials {
		credentialOpt = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	}
	conn, err := grpc.Dial(*address, credentialOpt, grpc.WithBlock())
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
	allReply, err := c.GetAll(context.Background(), &pb.AllRequest{})
	if err != nil {
		glog.Fatalf("Could not call GetAll: %s", err)
	}
	glog.Infof("AllReply: %v", allReply)
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
