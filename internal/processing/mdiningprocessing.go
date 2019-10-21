package mdiningprocessing

import (
	pb "github.com/MichiganDiningAPI/api/proto"
	"github.com/MichiganDiningAPI/util/containers"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

func MenusToFoods(menus *[]proto.Message) ([]proto.Message, error) {
	foods := make(map[string]proto.Message)
	for _, menu := range *menus {
		m := menu.(*pb.Menu)
		if m.Category == nil {
			continue
		}
		glog.Infof("Processing Dining Hall Menu: %s", m.DiningHallName)
		for _, cat := range m.Category {
			if cat == nil {
				glog.Warningf("Category nil for menu %s", m.Key)
				continue
			}
			for _, menuItem := range cat.MenuItem {
				if menuItem == nil {
					glog.Warningf("MenuItem is nil for category %s in menu %s", cat.Name, m.Key)
					continue
				}
				food, exists := foods[menuItem.Name]
				if !exists {
					foods[menuItem.Name] = &pb.Food{
						Key:             menuItem.Name + m.Date,
						Date:            m.Date,
						Name:            menuItem.Name,
						Category:        cat.Name,
						MenuItem:        menuItem,
						DiningHallMatch: map[string]*pb.FoodDiningHallMatch{}}
					food, _ = foods[menuItem.Name]
				}
				f := food.(*pb.Food)
				var match *pb.FoodDiningHallMatch
				match, exists = f.DiningHallMatch[m.DiningHallName]
				if !exists {
					match = &pb.FoodDiningHallMatch{Name: m.DiningHallName, MealTime: map[string]*pb.MealTime{}}
					f.DiningHallMatch[m.DiningHallName] = match
				}
				var mealTime *pb.MealTime
				mealTime, exists = match.MealTime[m.Date]
				if !exists {
					mealTime = &pb.MealTime{Date: m.Date, FormattedDate: m.FormattedDate, MealNames: []string{}}
					match.MealTime[m.Date] = mealTime
				}
				mealTime.MealNames = append(mealTime.MealNames, m.Meal)
			}
		}
	}
	return util.AsSliceType(util.Values(foods), []proto.Message{}).([]proto.Message), nil
}
