package main

import (
	"flag"

	pb "github.com/MichiganDiningAPI/api/proto"
	dc "github.com/MichiganDiningAPI/cmd/fetch/dynamoclient"
	"github.com/MichiganDiningAPI/util/io"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()

	dynamoclient := dc.New()
	dynamoclient.CreateTables()
	var mockDiningHalls pb.DiningHalls
	var mockItems pb.Items
	if util.ReadProtoFromFile("api/proto/sample/dininghalls.proto.txt", &mockDiningHalls) != nil {
		glog.Fatalf("Failed to read dining hall proto")
	}
	if util.ReadProtoFromFile("api/proto/sample/items.proto.txt", &mockItems) != nil {
		glog.Fatalf("Failed to read items proto")
	}
	dynamoclient.PutProto(&dc.DiningHallsTableName, mockDiningHalls.DiningHalls[0])
	tendies, exists := mockItems.Items["apple braised swiss chard"]
	if !exists {
		glog.Fatalf("Failed to find item in map")
	}
	dynamoclient.PutProto(&dc.ItemsTableName, tendies)
	var item pb.Item
	dynamoclient.GetProto(dc.ItemsTableName, tendies.Name, &item)
	glog.Infof("Result of Get: %v", item)
}
