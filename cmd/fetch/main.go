package main

import (
	"flag"
	"sync"

	pb "github.com/MichiganDiningAPI/api/proto"
	dc "github.com/MichiganDiningAPI/cmd/fetch/dynamoclient"
	mc "github.com/MichiganDiningAPI/cmd/fetch/mdiningclient"
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
		dynamoclient.PutProtoBatch(&dc.MenuTableName, menus)
		wg.Done()
	}()
	foods := make(map[string]proto.Message)
	for _, menu := range menus {
		m := menu.(*pb.Menu)
		if m.Category == nil {
			continue
		}
		glog.Infof("Menu: %s", m.DiningHallName)
		for _, cat := range m.Category {
			if cat == nil {
				glog.Infof("Dining Hall: %s", m.DiningHallName)
				glog.Fatalf("Nil Cat: %v", m)
				continue
			}
			glog.Infof("Cat: %s", cat.Name)
			for _, menuItem := range cat.MenuItem {
				if menuItem == nil {
					glog.Fatalf("Cat: %v", cat)
				}
				food, exists := foods[menuItem.Name]
				if !exists {
					foods[menuItem.Name] = &pb.Food{Key: menuItem.Name + m.Date, Date: m.Date, Name: menuItem.Name, Category: cat.Name, MenuItem: menuItem, DiningHallMatch: map[string]*pb.FoodDiningHallMatch{}}
					food, _ = foods[menuItem.Name]
				}
				f := food.(*pb.Food)
				match, e := f.DiningHallMatch[m.DiningHallName]
				if !e {
					match = &pb.FoodDiningHallMatch{Name: m.DiningHallName, MealTime: map[string]*pb.MealTime{}}
					f.DiningHallMatch[m.DiningHallName] = match
				}
				mealTime, e2 := match.MealTime[m.Date]
				if !e2 {
					mealTime = &pb.MealTime{Date: m.Date, FormattedDate: m.FormattedDate, MealNames: []string{}}
					match.MealTime[m.Date] = mealTime
				}
				mealTime.MealNames = append(mealTime.MealNames, m.Meal)
			}
		}
	}
	foodsSlice := make([]proto.Message, 0, len(foods))
	for _, v := range foods {
		foodsSlice = append(foodsSlice, v)
	}
	wg.Add(1)
	go func() {
		dynamoclient.PutProtoBatch(&dc.FoodTableName, foodsSlice)
		wg.Done()
	}()
	wg.Wait()
}
