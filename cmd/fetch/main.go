package main

import (
	"flag"
	"sync"

	mc "github.com/MichiganDiningAPI/api/mdining/mdiningclient"
	dc "github.com/MichiganDiningAPI/db/dynamoclient"
	"github.com/MichiganDiningAPI/util/containers"
	"github.com/MichiganDiningAPI/internal/processing/mdiningprocessing"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

func main() {
	flag.Parse()

	mdining := mc.New()

	dynamoclient := dc.New()
	dynamoclient.CreateTablesIfNotExists()

	dh, e := mdining.GetDiningHallList()
	if e != nil {
		glog.Fatalf("Failed to get dining hall list %s", e)
	}
	dynamoclient.PutProtoBatch(&dc.DiningHallsTableName, 
		util.AsSliceType(dh.DiningHalls, []proto.Message{}).([]proto.Message))
	menus, err := mdining.GetAllMenus(dh)
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
