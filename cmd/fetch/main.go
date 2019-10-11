package main

import (
	"context"
	"flag"

	pb "github.com/MichiganDiningAPI/api/proto"
	"github.com/MichiganDiningAPI/util/io"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/golang/glog"
)

var (
	DiningHallsTableName = "DiningHalls"
	ItemsTableName       = "Items"
)

var TableNames = []string{
	DiningHallsTableName,
	ItemsTableName}
var TableKeys = []string{
	"Name",
	"Name"}

type DynamoClient struct {
	client *dynamodb.Client
}

func diningHallToItem(diningHall *pb.DiningHall) *map[string]dynamodb.AttributeValue {
	item := make(map[string]dynamodb.AttributeValue)
	item["Name"] = dynamodb.AttributeValue{S: &diningHall.Name}
	item["Campus"] = dynamodb.AttributeValue{S: &diningHall.Campus}
	// TODO: Rest of attributes
	return &item
}

func itemToItem(item *pb.Item) *map[string]dynamodb.AttributeValue {
	i := make(map[string]dynamodb.AttributeValue)
	i["Name"] = dynamodb.AttributeValue{S: &item.Name}
	i["Attributes"] = dynamodb.AttributeValue{SS: item.Attributes}
	return &i
}

func (d *DynamoClient) createTable(table string, key string) {
	read, write := int64(5), int64(5)
	createReq := d.client.CreateTableRequest(&dynamodb.CreateTableInput{
		TableName:             &table,
		KeySchema:             []dynamodb.KeySchemaElement{dynamodb.KeySchemaElement{AttributeName: &key, KeyType: "HASH"}},
		AttributeDefinitions:  []dynamodb.AttributeDefinition{dynamodb.AttributeDefinition{AttributeName: &key, AttributeType: dynamodb.ScalarAttributeTypeS}},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{ReadCapacityUnits: &read, WriteCapacityUnits: &write}})
	_, err := createReq.Send(context.Background())
	if err != nil {
		glog.Fatalf("Failed to create table %s %v", table, err)
	}
	glog.Infof("Created table %s.", table)
}

func (d *DynamoClient) createTables() error {
	glog.Info("Checking for existence of dynamodb tables...")
	for i, table := range TableNames {
		describeReq := d.client.DescribeTableRequest(&dynamodb.DescribeTableInput{TableName: aws.String(table)})
		_, err := describeReq.Send(context.Background())
		if err != nil {
			glog.Infof("Table %s does not exist. Creating now...", table)
			d.createTable(table, TableNames[i])
			glog.Infof("Created table %s.", table)
		} else {
			glog.Infof("Table %s exists.", table)
		}
	}
	return nil
}

func (d *DynamoClient) addDiningHall(diningHall *pb.DiningHall) error {
	item := diningHallToItem(diningHall)
	req := d.client.PutItemRequest(&dynamodb.PutItemInput{
		TableName: &DiningHallsTableName,
		Item:      *item})
	_, err := req.Send(context.Background())
	if err != nil {
		glog.Errorf("Error adding dining hall %s %s", diningHall.Name, err)
		return err
	}
	glog.Infof("Added dining hall %s", diningHall.Name)
	return nil
}

func (d *DynamoClient) addItem(item *pb.Item) error {
	i := itemToItem(item)
	req := d.client.PutItemRequest(&dynamodb.PutItemInput{
		TableName: &ItemsTableName,
		Item:      *i})
	_, err := req.Send(context.Background())
	if err != nil {
		glog.Errorf("Error adding item %s %s", item.Name, err)
		return err
	}
	glog.Infof("Added item %s", item.Name)
	return nil
}

func main() {
	flag.Parse()
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		glog.Fatalf("Unable to load SDK config, %v" + err.Error())
	}

	// Set the AWS Region that the service clients should use
	cfg.Region = endpoints.UsEast1RegionID

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.New(cfg)

	d := new(DynamoClient)
	d.client = svc
	d.createTables()
	var mockDiningHalls pb.DiningHalls
	if util.ReadProtoFromFile("api/proto/sample/dininghalls.proto.txt", &mockDiningHalls) != nil {
		glog.Fatalf("Failed to read dining hall proto")
	}
	av, err := dynamodbattribute.MarshalMap(&mockDiningHalls.DiningHalls[0])
	if err != nil {
		glog.Fatal("Bad marshall")
	}
	glog.Infof("Marshalled: %v", av)
	d.addDiningHall(&pb.DiningHall{Name: "Test", Campus: "Testing2"})
	d.addItem(&pb.Item{Name: "Cheese", Attributes: []string{"Attr1", "Attr2"}})
}
