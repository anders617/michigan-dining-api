package dynamoclient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/golang/glog"
)

func (d *DynamoClient) CreateTablesIfNotExists() {
	glog.Info("Checking for existence of dynamodb tables...")
	for _, table := range TableNames {
		if !d.tableExists(table) {
			glog.Infof("Table %s does not exist. Creating now...", table)
			d.createTable(table)
		} else {
			glog.Infof("Table %s exists.", table)
		}
	}
}

func (d *DynamoClient) CreateTables() error {
	glog.Info("Creating dynamodb tables...")
	for _, table := range TableNames {
		glog.Infof("Creating table %s...", table)
		d.createTable(table)
		glog.Infof("Created table %s.", table)
	}
	return nil
}

func (d *DynamoClient) tableExists(table string) bool {
	describeReq := d.client.DescribeTableRequest(&dynamodb.DescribeTableInput{
		TableName: aws.String(table)})
	_, err := describeReq.Send(context.Background())
	return err == nil
}

func (d *DynamoClient) createTable(table string) {
	read, write := int64(5), int64(5)
	keys, _ := TableKeys[table]
	attrs, _ := TableAttributes[table]
	streamSpec, _ := TableStreamSpecs[table]
	createReq := d.client.CreateTableRequest(&dynamodb.CreateTableInput{
		TableName:             &table,
		KeySchema:             keys,
		AttributeDefinitions:  attrs,
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{ReadCapacityUnits: &read, WriteCapacityUnits: &write},
		StreamSpecification:   &streamSpec})
	_, err := createReq.Send(context.Background())
	if err != nil {
		glog.Fatalf("Failed to create table %s %v", table, err)
	}
	glog.Infof("Created table %s.", table)
}
