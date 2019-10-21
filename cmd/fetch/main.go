package main

import (
	"flag"
	"sync"

	mc "github.com/MichiganDiningAPI/api/mdining/mdiningclient"
	dc "github.com/MichiganDiningAPI/db/dynamoclient"
	"github.com/MichiganDiningAPI/internal/processing/mdiningprocessing"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

func main() {
	flag.Parse()

	mdining := mc.New()

	dynamoclient := dc.New()
	dynamoclient.CreateTables()
	// var mockDiningHalls pb.DiningHalls
	// var mockItems pb.Items
	// if util.ReadProtoFromFile("api/proto/sample/dininghalls.proto.txt", &mockDiningHalls) != nil {
	// 	glog.Fatalf("Failed to read dining hall proto")
	// }
	// if util.ReadProtoFromFile("api/proto/sample/items.proto.txt", &mockItems) != nil {
	// 	glog.Fatalf("Failed to read items proto")
	// }
	// dynamoclient.PutProto(&dc.DiningHallsTableName, mockDiningHalls.DiningHalls[0])
	// tendies, exists := mockItems.Items["apple braised swiss chard"]
	// if !exists {
	// 	glog.Fatalf("Failed to find item in map")
	// }
	// dynamoclient.PutProto(&dc.ItemsTableName, tendies)
	// var item pb.Item
	// dynamoclient.GetProto(dc.ItemsTableName, map[string]string{dc.NameKey: tendies.Name}, &item)
	// glog.Infof("Result of Get: %v", item)

	dh, e := mdining.GetDiningHallList()
	if e != nil {
		glog.Fatalf("Failed to get dining hall list %s", e)
	}
	menus := make([]proto.Message, 0)
	for _, dininghall := range dh.DiningHalls {
		m, err := mdining.GetMenus(dininghall)
		if err != nil {
			continue
		}
		for _, menu := range *m {
			menus = append(menus, menu)
		}
	}
	glog.Infof("Menus count: %d", len(menus))
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		// dynamoclient.PutProtoBatch(&dc.MenuTableName, menus)
		wg.Done()
	}()
	foodsSlice, err := mdiningprocessing.MenusToFoods(&menus)
	if err != nil {
		glog.Warningf("Could not convert menus to foods %s", err)
	} else {
		wg.Add(1)
		go func() {
			dynamoclient.PutProtoBatch(&dc.FoodTableName, foodsSlice)
			wg.Done()
		}()
	}
	wg.Wait()
}
