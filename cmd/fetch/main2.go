package main

import (
	"flag"
	"sync"
	"time"

	"github.com/MichiganDiningAPI/api/mdining/mdiningclient2"
	dc "github.com/MichiganDiningAPI/db/dynamoclient"
	"github.com/MichiganDiningAPI/internal/processing/mdiningprocessing"
	util "github.com/MichiganDiningAPI/internal/util/containers"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

func main() {
	var apiKey = flag.String("api_key", "", "API Key for mdining api")
	var numDays = flag.Int("num_days", 7, "Number of days of data (from today) to retrieve.")
	flag.Parse()
	if len(*apiKey) == 0 {
		glog.Errorf("Missing required argument api_key.")
		return
	}
	glog.Infof("Using api key: %s", *apiKey)
	client := mdiningclient2.New(*apiKey)
	dynamoclient := dc.New()
	dynamoclient.CreateTablesIfNotExists()
	oneDay, _ := time.ParseDuration("24h")
	currentDate := time.Now()
	dates := []time.Time{}
	for i := 0; i < *numDays; i++ {
		dates = append(dates, currentDate)
		currentDate = currentDate.Add(oneDay)
	}
	diningHallsByCampus, partialMenus, err := client.GetDiningHallList(dates)
	if err != nil {
		glog.Errorf("Failed to get dining hall list %s", err)
		return
	}
	diningHallsList := []proto.Message{}
	for campus, diningHalls := range *diningHallsByCampus {
		glog.Infof("Received campus: %s", campus)
		diningHallsList = append(diningHallsList, util.AsSliceType(diningHalls.DiningHalls, []proto.Message{}).([]proto.Message)...)
	}
	dynamoclient.PutProtoBatch(&dc.DiningHallsTableName, diningHallsList)
	menus, err := client.GetAllMenus(partialMenus)
	if err != nil {
		glog.Errorf("Failed to get menus %s", err)
		return
	}
	menusProtoSlice := util.AsSliceType(menus, []proto.Message{}).([]proto.Message)
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
