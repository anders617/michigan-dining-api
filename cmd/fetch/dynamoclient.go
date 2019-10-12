package dynamoclient

import (
	"context"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

var (
	DiningHallsTableName = "DiningHalls"
	ItemsTableName       = "Items"
)

var (
	TableNames = []string{
		DiningHallsTableName,
		ItemsTableName}
	TableKeys = []string{
		"name",
		"name"}
)

type DynamoClient struct {
	client *dynamodb.Client
}

func New() *DynamoClient {
	dc := new(DynamoClient)
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		glog.Fatalf("Unable to load SDK config, %v" + err.Error())
	}
	// TODO: Make this configurable
	cfg.Region = endpoints.UsEast1RegionID
	dc.client = dynamodb.New(cfg)
	return dc
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

func (d *DynamoClient) CreateTables() error {
	glog.Info("Checking for existence of dynamodb tables...")
	for i, table := range TableNames {
		describeReq := d.client.DescribeTableRequest(&dynamodb.DescribeTableInput{TableName: aws.String(table)})
		_, err := describeReq.Send(context.Background())
		if err != nil {
			glog.Infof("Table %s does not exist. Creating now...", table)
			d.createTable(table, TableKeys[i])
			glog.Infof("Created table %s.", table)
		} else {
			glog.Infof("Table %s exists.", table)
		}
	}
	return nil
}

func (d *DynamoClient) PutProto(table *string, p proto.Message) error {
	// Convert from proto to dynamodb friendly structure
	av, err := dynamodbattribute.MarshalMap(&p)
	if err != nil {
		return err
	}
	// Create and send put request
	req := d.client.PutItemRequest(&dynamodb.PutItemInput{
		TableName: table,
		Item:      av})
	_, err = req.Send(context.Background())
	if err != nil {
		glog.Errorf("Error putting item %s", err)
		return err
	}
	glog.Infof("Successfully Put %s", reflect.TypeOf(p))
	return nil
}
