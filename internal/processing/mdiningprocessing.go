package mdiningprocessing

import (
	"strings"

	pb "github.com/MichiganDiningAPI/api/proto"
	"github.com/MichiganDiningAPI/internal/util/containers"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

func ItemsToFilterableEntries(items *pb.Items) *pb.FilterableEntries {
	filterableEntries := pb.FilterableEntries{}
	filterableEntries.FilterableEntries = make([]*pb.FilterableEntry, 0, len(items.Items))
	for _, item := range items.Items {
		for _, match := range item.DiningHallMatches {
			for _, time := range match.MealTimes {
				entry := pb.FilterableEntry{
					ItemName: item.Name,
					Date: time.Date,
					DiningHallName: match.Name,
					MealNames: time.MealNames,
					Attributes: item.Attributes,
				}
				filterableEntries.FilterableEntries = append(filterableEntries.FilterableEntries, &entry)
			}
		}
	}
	return &filterableEntries
}

func FoodsToItems(foods *[]*pb.Food) *pb.Items {
	items := pb.Items{Items: map[string]*pb.Item{}}
	for _, food := range *foods {
		i := pb.Item{
			Name: food.Name,
			Attributes: food.MenuItem.Attribute,
		}
		i.DiningHallMatches = make(map[string]*pb.Item_DiningHallMatch)
		i.DiningHallMatchesArray = make([]*pb.Item_DiningHallMatch, 0, len(food.DiningHallMatch))
		for _, match := range food.DiningHallMatch {
			itemMatch := FoodDiningHallMatchToDiningHallMatch(match)
			i.DiningHallMatches[match.Name] = itemMatch
			i.DiningHallMatchesArray = append(i.DiningHallMatchesArray, itemMatch)
		}
		items.Items[food.Name] = &i
	}
	return &items
}

func FoodDiningHallMatchToDiningHallMatch(f *pb.FoodDiningHallMatch) *pb.Item_DiningHallMatch {
	diningHallMatch := pb.Item_DiningHallMatch {
		Name: f.Name,
		MealTimes: f.MealTime,
		MealTimesArray: util.AsSliceType(util.Values(f.MealTime), []*pb.MealTime{}).([]*pb.MealTime),
	}
	return &diningHallMatch
}

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
				glog.Warningf("Category nil for menu %s", m.DiningHallMeal+m.Date)
				continue
			}
			for _, menuItem := range cat.MenuItem {
				if menuItem == nil {
					glog.Warningf("MenuItem is nil for category %s in menu %s", cat.Name, m.DiningHallMeal+m.Date)
					continue
				}
				food, exists := foods[menuItem.Name]
				if !exists {
					foods[menuItem.Name] = &pb.Food{
						Key:             strings.ToLower(menuItem.Name),
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
