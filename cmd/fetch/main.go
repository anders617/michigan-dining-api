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
	dynamoclient.CreateTablesIfNotExists()

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
