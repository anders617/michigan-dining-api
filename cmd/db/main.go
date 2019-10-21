package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	dc "github.com/MichiganDiningAPI/db/dynamoclient"
	"github.com/golang/glog"
)

func main() {
	create := flag.Bool("create", false, "Specify this flag to create necessary tables on dynamodb")
	delete := flag.Bool("delete", false, "Specify this glag to delete necessary table on dynamo db")
	flag.Parse()

	if *create && *delete {
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
}
