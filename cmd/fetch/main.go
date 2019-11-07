package main

import (
	"flag"
	"sync"

	mc "github.com/MichiganDiningAPI/api/mdining/mdiningclient"
	dc "github.com/MichiganDiningAPI/db/dynamoclient"
	"github.com/MichiganDiningAPI/internal/util/containers"
	"github.com/MichiganDiningAPI/internal/processing/mdiningprocessing"
	pb "github.com/anders617/mdining-proto/proto/mdining"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

func main() {
	flag.Parse()

	mdining := mc.New()

	dynamoclient := dc.New()
	dynamoclient.CreateTablesIfNotExists()

	diningHallsByCampus, err := mdining.GetDiningHallList()
	if err != nil {
		glog.Fatalf("Failed to get dining hall list %s", err)
	}
	diningHallsList := []proto.Message{}
	for campus, diningHalls := range *diningHallsByCampus {
		glog.Infof("Received campus: %s", campus)
		diningHallsList = append(diningHallsList, util.AsSliceType(diningHalls.DiningHalls, []proto.Message{}).([]proto.Message)...)
	}
	dynamoclient.PutProtoBatch(&dc.DiningHallsTableName, diningHallsList)
	menus := []*pb.Menu{}
	for _, diningHalls := range *diningHallsByCampus {
		m, err := mdining.GetAllMenus(diningHalls)
		if err != nil {
			glog.Fatalf("Failed to get menus %s", err)
		}
		menus = append(menus, *m...)
	}
	menusProtoSlice := util.AsSliceType(menus, []proto.Message{}).([]proto.Message)
	if err != nil {
		glog.Fatalf("Failed to get menus %s", err)
	}
	glog.Infof("Menus count: %d", len(menusProtoSlice))
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		dynamoclient.PutProtoBatch(&dc.MenuTableName, menusProtoSlice)
		wg.Done()
	}()
	foodsSlice, err := mdiningprocessing.MenusToFoods(&menusProtoSlice)
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
