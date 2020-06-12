package main

import (
	"flag"

	"github.com/MichiganDiningAPI/db/dynamoclient"
	containers "github.com/MichiganDiningAPI/internal/util/containers"
	"github.com/MichiganDiningAPI/internal/util/date"
	pb "github.com/anders617/mdining-proto/proto/mdining"
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
		if dh.Campus != "" && dh.Campus != "DINING HALLS" {
			// For now, only analyze actual Dining Hall foods
			continue
		}
		for _, mealTime := range dh.MealTime {
			count += int64(len(mealTime.MealNames))
		}
	}
	return count
}

func updateStats(foodStats *pb.FoodStat, food *pb.Food) {
	timesServed := countTimesServed(food)
	foodStats.TotalFoodMealsServed += timesServed
	foodStats.TimesServed[food.Key] += timesServed
	for _, cat := range food.Category {
		foodStats.CategoryCounts[cat] += timesServed
	}
	_, e := foodStats.FoodWeekdayCounts[food.Key]
	if !e {
		foodStats.FoodWeekdayCounts[food.Key] = &pb.StringToInt{Data: make(map[string]int64)}
	}
	d, _ := date.ParseNoTime(&food.Date)
	foodStats.FoodWeekdayCounts[food.Key].Data[d.Weekday().String()] += timesServed
	_, e = foodStats.WeekdayFoodCounts[d.Weekday().String()]
	if !e {
		foodStats.WeekdayFoodCounts[d.Weekday().String()] = &pb.StringToInt{Data: make(map[string]int64)}
	}
	foodStats.WeekdayFoodCounts[d.Weekday().String()].Data[food.Key] += timesServed
	_, e = foodStats.FoodDiningHallCounts[food.Key]
	if !e {
		foodStats.FoodDiningHallCounts[food.Key] = &pb.StringToInt{Data: make(map[string]int64)}
	}
	for dhName, dh := range food.DiningHallMatch {
		if dh.Campus != "" && dh.Campus != "DINING HALLS" {
			// For now, only analyze actual Dining Hall foods
			continue
		}
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
}

func NewFoodStat(date string) *pb.FoodStat {
	return &pb.FoodStat{
		Date:                  date,
		TimesServed:           map[string]int64{},
		FoodDiningHallCounts:  map[string]*pb.StringToInt{},
		DiningHallFoodCounts:  map[string]*pb.StringToInt{},
		CategoryCounts:        map[string]int64{},
		AllergenCounts:        map[string]int64{},
		AttributeCounts:       map[string]int64{},
		WeekdayFoodCounts:     map[string]*pb.StringToInt{},
		FoodWeekdayCounts:     map[string]*pb.StringToInt{},
		NumUniqueFoods:        0,
		TotalFoodMealsServed:  0,
		DiningHallMealsServed: map[string]int64{},
	}
}

func main() {
	flag.Parse()

	dc := dynamoclient.New()

	stats := map[string]*pb.FoodStat{}

	// Find all foods and calculate
	//startDate := date.FormatNoTime(date.Now())
	dc.ForEachFood(nil, nil, func(food *pb.Food) {
		stat, exists := stats[food.Date]
		if !exists {
			stat = NewFoodStat(food.Date)
			stats[food.Date] = stat
		}
		updateStats(stat, food)
	})

	for _, stat := range stats {
		stat.NumUniqueFoods = int64(len(stat.TimesServed))
	}

	// Push results to dynamodb
	for date, stat := range stats {
		glog.Infof("Putting stats for date %s", date)
		err := dc.PutProto(&dynamoclient.FoodStatsTableName, stat)
		if err != nil {
			glog.Fatalf("Error putting proto: %s", err)
		}
		glog.Infof("Sucessfully put stats for date %s", date)
	}

}
