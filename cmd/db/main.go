package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	dc "github.com/MichiganDiningAPI/db/dynamoclient"
	"github.com/golang/glog"
)

func toInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func main() {
	create := flag.Bool("create", false, "Specify this flag to create necessary tables on dynamodb")
	delete := flag.Bool("delete", false, "Specify this flag to delete necessary table on dynamo db")
	query := flag.Bool("query", false, "Specify this flag to query tables")
	stream := flag.Bool("stream", false, "Specify this flag to stream from the hearts table")
	flag.Parse()

	if toInt(*create)+toInt(*delete)+toInt(*query)+toInt(*stream) > 1 {
		glog.Fatal("You must specify either create or delete, not both")
	}

	dynamoclient := dc.New()
	if *create {
		dynamoclient.CreateTablesIfNotExists()
	}
	if *delete {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("WARNING! This is an irreversible action that will result in data loss.\n")
		fmt.Printf("Are you sure you want to delete tables? [y/n]")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, " \n\t")
		if text == "y" {
			dynamoclient.DeleteTables()
		} else {
			fmt.Printf("Not Deleting!\n")
		}
	}
	if *stream {
		records, done := dynamoclient.StreamHearts()
		time.AfterFunc(time.Second*10, func() { done <- struct{}{} })
		for record := range records {
			glog.Infof("Record: %v", record)
		}
	}
}
