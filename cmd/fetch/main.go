package main

import (
	"flag"
	"sync"

	pb "github.com/MichiganDiningAPI/api/proto"
	dc "github.com/MichiganDiningAPI/cmd/fetch/dynamoclient"
	mc "github.com/MichiganDiningAPI/cmd/fetch/mdiningclient"
	"github.com/MichiganDiningAPI/util/io"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()

	mdining := mc.New()

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


	dh, e := mdining.GetDiningHallList()
	if e != nil {
		glog.Fatalf("Failed to get dining hall list %s", e)
	}
	wg := sync.WaitGroup{}
	for _, dininghall := range dh.DiningHalls {
		reply, err := mdining.GetMenuBase(dininghall)
		if err != nil {
			continue
		}
		wg.Add(1)
		go func() {
			for _, menuBase := range reply.Menu {
				r, _ := mdining.GetMenuDetails(dininghall, menuBase)
				glog.Infof("%v", r)
			}
			wg.Done()
		}()
		//dynamoclient.PutProto(&dc.DiningHallsTableName, dininghall)
	}
	wg.Wait()
}
