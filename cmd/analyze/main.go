package main

import (
	"flag"
	"time"

	pb "github.com/MichiganDiningAPI/api/proto"
	"github.com/MichiganDiningAPI/db/dynamoclient"
	"github.com/MichiganDiningAPI/internal/util/io"
	"github.com/MichiganDiningAPI/internal/util/date"
	containers "github.com/MichiganDiningAPI/internal/util/containers"
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

func countTimesServed(food *pb.Food) int {
	count := 0
	for _, dh := range food.DiningHallMatch {
		for _, mealTime := range dh.MealTime {
			count += len(mealTime.MealNames)
		}
	}
	return count
}

func main() {
	flag.Parse()

	dc := dynamoclient.New()

	counts := make(map[string]int)
	foodDhCounts := make(map[string]map[string]int)
	dhFoodCounts := make(map[string]map[string]int)
	categoryCounts := make(map[string]int)
	allergenCounts := make(map[string]int)
	attributeCounts := make(map[string]int)
	foodWeekdayCounts := make(map[string]map[time.Weekday]int)
	weekdayFoodCounts := make(map[time.Weekday]map[string]int)

	dc.ForEachFood(func(food *pb.Food) {
		timesServed := countTimesServed(food)
		counts[food.Key] += timesServed
		categoryCounts[food.Category] += timesServed
		_, e := foodWeekdayCounts[food.Key]
		if !e {
			foodWeekdayCounts[food.Key] = make(map[time.Weekday]int)
		}
		d, _:= date.Parse(&food.Date)
		foodWeekdayCounts[food.Key][d.Weekday()] += timesServed
		_, e = weekdayFoodCounts[d.Weekday()]
		if !e {
			weekdayFoodCounts[d.Weekday()] = make(map[string]int)
		}
		weekdayFoodCounts[d.Weekday()][food.Key] += timesServed
		_, e = foodDhCounts[food.Key]
		if !e {
			foodDhCounts[food.Key] = make(map[string]int)
		}
		for dhName, dh := range food.DiningHallMatch {
			_, e = dhFoodCounts[dhName]
			if !e {
				dhFoodCounts[dhName] = make(map[string]int)
			}
			for range dh.MealTime {
				foodDhCounts[food.Key][dhName]++
				dhFoodCounts[dhName][food.Key]++
			}
		}
		if len(food.MenuItem.Allergens) == 0 {
			allergenCounts["none"] += timesServed
		}
		for _, allergen := range food.MenuItem.Allergens {
			allergenCounts[allergen] += timesServed
		}
		if len(food.MenuItem.Attribute) == 0 {
			attributeCounts["none"] += timesServed
		}
		for _, attribute := range food.MenuItem.Attribute {
			attributeCounts[attribute] += timesServed
		}
	})

	glog.Infof("Counts:\n%s", util.MapToString(counts))
	glog.Infof("Counts:\n%s", util.MapToString(foodDhCounts))
	glog.Infof("Counts:\n%s", util.MapToString(dhFoodCounts))
	glog.Infof("Counts:\n%s", util.MapToString(categoryCounts))
	logStats(categoryCounts)
	glog.Infof("Counts:\n%s", util.MapToString(allergenCounts))
	logStats(allergenCounts)
	glog.Infof("Counts:\n%s", util.MapToString(attributeCounts))
	logStats(attributeCounts)
	glog.Infof("Counts:\n%s", util.MapToString(foodWeekdayCounts))
	glog.Infof("Counts:\n%s", util.MapToString(weekdayFoodCounts))
}
