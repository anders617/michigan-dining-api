package mdiningprocessing

import (
	"strings"

	util "github.com/MichiganDiningAPI/internal/util/containers"
	pb "github.com/anders617/mdining-proto/proto/mdining"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

func FoodStatsToSummaryStats(foodStats *[]*pb.FoodStat) *pb.SummaryStats {
	summaryStats := &pb.SummaryStats{
		Dates:                   []string{},
		NumUniqueFoods:          []int64{},
		TotalFoodMealsServed:    []int64{},
		AllergenCounts:          map[string]*pb.CountArray{},
		AttributeCounts:         map[string]*pb.CountArray{},
		LatestWeekdayFoodCounts: map[string]*pb.StringToInt{},
	}
	allergens := map[string]bool{}
	attributes := map[string]bool{}
	days := map[string]bool{}
	foods := map[string]bool{}
	for _, stat := range *foodStats {
		for allergen := range stat.AllergenCounts {
			allergens[allergen] = true
		}
		for attribute := range stat.AttributeCounts {
			attributes[attribute] = true
		}
		for day, counts := range stat.WeekdayFoodCounts {
			days[day] = true
			for food := range counts.Data {
				foods[food] = true
			}
		}
	}
	for allergen := range allergens {
		summaryStats.AllergenCounts[allergen] = &pb.CountArray{}
	}
	for attribute := range attributes {
		summaryStats.AttributeCounts[attribute] = &pb.CountArray{}
	}
	summaryStats.LatestWeekdayFoodCounts = map[string]*pb.StringToInt{}
	for day := range days {
		summaryStats.LatestWeekdayFoodCounts[day] = &pb.StringToInt{Data: map[string]int64{}}
		for food := range foods {
			summaryStats.LatestWeekdayFoodCounts[day].Data[food] = 0
		}
	}
	for _, stat := range *foodStats {
		summaryStats.Dates = append(summaryStats.Dates, stat.Date)
		summaryStats.NumUniqueFoods = append(summaryStats.NumUniqueFoods, stat.NumUniqueFoods)
		summaryStats.TotalFoodMealsServed = append(summaryStats.TotalFoodMealsServed, stat.TotalFoodMealsServed)

		for allergen := range allergens {
			count := int64(-1)
			if val, exists := stat.AllergenCounts[allergen]; exists {
				count = val
			}
			summaryStats.AllergenCounts[allergen].Counts = append(summaryStats.AllergenCounts[allergen].Counts, count)
		}
		for attribute := range attributes {
			count := int64(-1)
			if val, exists := stat.AttributeCounts[attribute]; exists {
				count = val
			}
			summaryStats.AttributeCounts[attribute].Counts = append(summaryStats.AttributeCounts[attribute].Counts, count)
		}
		for day, counts := range stat.WeekdayFoodCounts {
			for food, count := range counts.Data {
				summaryStats.LatestWeekdayFoodCounts[day].Data[food] += count
			}
		}
	}
	return summaryStats
}

func ItemsToFilterableEntries(items *pb.Items) *pb.FilterableEntries {
	filterableEntries := pb.FilterableEntries{}
	filterableEntries.FilterableEntries = make([]*pb.FilterableEntry, 0, len(items.Items))
	for _, item := range items.Items {
		for _, match := range item.DiningHallMatches {
			for _, time := range match.MealTimes {
				entry := pb.FilterableEntry{
					ItemName:       item.Name,
					Date:           time.Date,
					DiningHallName: match.Name,
					MealNames:      time.MealNames,
					Attributes:     item.Attributes,
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
		item, exists := items.Items[food.Key]
		if !exists {
			item = &pb.Item{
				Name:       food.Name,
				Attributes: food.MenuItem.Attribute,
			}
			item.DiningHallMatches = make(map[string]*pb.Item_DiningHallMatch)
			item.DiningHallMatchesArray = make([]*pb.Item_DiningHallMatch, 0, len(food.DiningHallMatch))
			items.Items[food.Key] = item
		}
		for _, match := range food.DiningHallMatch {
			itemMatch, exists := item.DiningHallMatches[match.Name]
			if match.Campus != "" && match.Campus != "DINING HALLS" {
				// For legacy purposes, Item should only contain dining hall foods
				continue
			}
			if !exists {
				itemMatch = &pb.Item_DiningHallMatch{
					Name:           match.Name,
					MealTimes:      map[string]*pb.MealTime{},
					MealTimesArray: []*pb.MealTime{},
				}
				item.DiningHallMatches[match.Name] = itemMatch
				item.DiningHallMatchesArray = append(item.DiningHallMatchesArray, itemMatch)
			}
			for key, value := range match.MealTime {
				_, exists := itemMatch.MealTimes[key]
				if !exists {
					itemMatch.MealTimes[key] = value
					itemMatch.MealTimesArray = append(itemMatch.MealTimesArray, value)
				}
			}
		}
		if len(item.DiningHallMatches) == 0 {
			// If we didn't add any dining hall matches, having an item is pointless
			delete(items.Items, food.Key)
		}
	}
	return &items
}

func FoodDiningHallMatchToDiningHallMatch(f *pb.FoodDiningHallMatch) *pb.Item_DiningHallMatch {
	diningHallMatch := pb.Item_DiningHallMatch{
		Name:           f.Name,
		MealTimes:      f.MealTime,
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
				food, exists := foods[menuItem.Name+m.Date]
				if !exists {
					foods[menuItem.Name+m.Date] = &pb.Food{
						Key:             strings.ToLower(menuItem.Name),
						Date:            m.Date,
						Name:            menuItem.Name,
						Category:        []string{},
						MenuItem:        menuItem,
						DiningHallMatch: map[string]*pb.FoodDiningHallMatch{}}
					food, _ = foods[menuItem.Name+m.Date]
				}
				f := food.(*pb.Food)
				containsCategory := false
				for _, c := range f.Category {
					if c == cat.Name {
						containsCategory = true
					}
				}
				if !containsCategory {
					f.Category = append(f.Category, cat.Name)
				}
				var match *pb.FoodDiningHallMatch
				match, exists = f.DiningHallMatch[m.DiningHallName]
				if !exists {
					match = &pb.FoodDiningHallMatch{Name: m.DiningHallName, MealTime: map[string]*pb.MealTime{}, Campus: m.DiningHallCampus}
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
