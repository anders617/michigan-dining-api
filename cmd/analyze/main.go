package main

import (
	"flag"

	pb "github.com/MichiganDiningAPI/api/proto"
	"github.com/MichiganDiningAPI/db/dynamoclient"
	containers "github.com/MichiganDiningAPI/internal/util/containers"
	"github.com/MichiganDiningAPI/internal/util/date"
	"github.com/golang/glog"
	"github.com/montanaflynn/stats"
)

func logStats(container interface{}) {
	data := stats.LoadRawData(*containers.Values(container))
	mean, _ := data.Mean()
	median, _ := data.Median()
	mode, _ := data.Mode()
	glog.Infof("Stats: mean %f median %f mode %f", mean, median, mode)
}

func countTimesServed(food *pb.Food) int64 {
	count := int64(0)
	for _, dh := range food.DiningHallMatch {
		for _, mealTime := range dh.MealTime {
			count += int64(len(mealTime.MealNames))
		}
	}
	return count
}

func main() {
	flag.Parse()

	dc := dynamoclient.New()

	// counts := make(map[string]int)
	// foodDhCounts := make(map[string]map[string]int)
	// dhFoodCounts := make(map[string]map[string]int)
	// categoryCounts := make(map[string]int)
	// allergenCounts := make(map[string]int)
	// attributeCounts := make(map[string]int)
	// foodWeekdayCounts := make(map[string]map[time.Weekday]int)
	// weekdayFoodCounts := make(map[time.Weekday]map[string]int)

	foodStats := &pb.FoodStat{
		Date:                  date.Format(date.Now()),
		TimesServed:           map[string]int64{},
		FoodDiningHallCounts:  map[string]*pb.StringToInt{},
		DiningHallFoodCounts:  map[string]*pb.StringToInt{},
		CategoryCounts:        map[string]int64{},
		AllergenCounts:        map[string]int64{},
		AttributeCounts:       map[string]int64{},
		WeekdayFoodCounts:     map[int64]*pb.StringToInt{},
		FoodWeekdayCounts:     map[string]*pb.IntToInt{},
		NumUniqueFoods:        0,
		TotalFoodMealsServed:  0,
		DiningHallMealsServed: map[string]int64{},
	}
	glog.Infof("%v", foodStats)

	dc.ForEachFood(nil, nil, func(food *pb.Food) {
		timesServed := countTimesServed(food)
		foodStats.TotalFoodMealsServed += timesServed
		foodStats.TimesServed[food.Key] += timesServed
		for _, cat := range food.Category {
			foodStats.CategoryCounts[cat] += timesServed
		}
		_, e := foodStats.FoodWeekdayCounts[food.Key]
		if !e {
			foodStats.FoodWeekdayCounts[food.Key] = &pb.IntToInt{Data: make(map[int64]int64)}
		}
		d, _ := date.Parse(&food.Date)
		foodStats.FoodWeekdayCounts[food.Key].Data[int64(d.Weekday())] += timesServed
		_, e = foodStats.WeekdayFoodCounts[int64(d.Weekday())]
		if !e {
			foodStats.WeekdayFoodCounts[int64(d.Weekday())] = &pb.StringToInt{Data: make(map[string]int64)}
		}
		foodStats.WeekdayFoodCounts[int64(d.Weekday())].Data[food.Key] += timesServed
		_, e = foodStats.FoodDiningHallCounts[food.Key]
		if !e {
			foodStats.FoodDiningHallCounts[food.Key] = &pb.StringToInt{Data: make(map[string]int64)}
		}
		for dhName, dh := range food.DiningHallMatch {
			_, e = foodStats.DiningHallFoodCounts[dhName]
			if !e {
				foodStats.DiningHallFoodCounts[dhName] = &pb.StringToInt{Data: make(map[string]int64)}
			}
			for range dh.MealTime {
				foodStats.FoodDiningHallCounts[food.Key].Data[dhName]++
				foodStats.DiningHallFoodCounts[dhName].Data[food.Key]++
				foodStats.DiningHallMealsServed[dhName]++
			}
		}
		if len(food.MenuItem.Allergens) == 0 {
			foodStats.AllergenCounts["none"] += timesServed
		}
		for _, allergen := range food.MenuItem.Allergens {
			foodStats.AllergenCounts[allergen] += timesServed
		}
		if len(food.MenuItem.Attribute) == 0 {
			foodStats.AttributeCounts["none"] += timesServed
		}
		for _, attribute := range food.MenuItem.Attribute {
			foodStats.AttributeCounts[attribute] += timesServed
		}
	})

	foodStats.NumUniqueFoods = int64(len(foodStats.TimesServed))

	err := dc.PutProto(&dynamoclient.FoodStatsTableName, foodStats)
	if err != nil {
		glog.Fatalf("Error putting proto: %s", err)
	}
}
