package main

import (
	"flag"

	"github.com/MichiganDiningAPI/internal/util/io"
	containers "github.com/MichiganDiningAPI/internal/util/containers"
	pb "github.com/MichiganDiningAPI/api/proto"
	"github.com/MichiganDiningAPI/db/dynamoclient"
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

func main() {
	flag.Parse()

	dc := dynamoclient.New()

	counts := make(map[string]int)
	foodDhCounts := make(map[string]map[string]int)
	dhFoodCounts := make(map[string]map[string]int)
	categoryCounts := make(map[string]int)
	allergenCounts := make(map[string]int)
	attributeCounts := make(map[string]int)
	
	dc.ForEachFood(func(food *pb.Food) {
		counts[food.Key]++
		categoryCounts[food.Category]++
		_, e := foodDhCounts[food.Key]
		if !e {
			foodDhCounts[food.Key] = make(map[string]int)
		}
		for dhName := range food.DiningHallMatch {
			_, e = dhFoodCounts[dhName]
			if !e {
				dhFoodCounts[dhName] = make(map[string]int)
			}
			foodDhCounts[food.Key][dhName]++
			dhFoodCounts[dhName][food.Key]++
		}
		if len(food.MenuItem.Allergens) == 0 {
			allergenCounts["none"]++
		}
		for _, allergen := range food.MenuItem.Allergens {
			allergenCounts[allergen]++
		}
		if len(food.MenuItem.Attribute) == 0 {
			attributeCounts["none"]++
		}
		for _, attribute := range food.MenuItem.Attribute {
			attributeCounts[attribute]++
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
}
