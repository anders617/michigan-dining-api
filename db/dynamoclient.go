package dynamoclient

import (
	"context"
	"math"
	"reflect"
	"time"

	pb "github.com/anders617/mdining-proto/proto/mdining"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodbstreams"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type DynamoClient struct {
	client       *dynamodb.Client
	streamClient *dynamodbstreams.Client
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
	dc.streamClient = dynamodbstreams.New(cfg)
	return dc
}

func (d *DynamoClient) AddHeart(key string) (*pb.HeartCount, error) {
	updateExpression := expression.Add(expression.Name("count"), expression.Value(1))
	expr, _ := expression.NewBuilder().WithUpdate(updateExpression).Build()
	dynamoKey, err := dynamodbattribute.Marshal(&key)
	if err != nil {
		return nil, err
	}
	params := dynamodb.UpdateItemInput{
		TableName:                 &HeartsTableName,
		Key:                       map[string]dynamodb.AttributeValue{HeartsTableKey: *dynamoKey},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              dynamodb.ReturnValueAllNew,
	}
	req := d.client.UpdateItemRequest(&params)

	resp, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	heartCount := pb.HeartCount{}
	err = dynamodbattribute.UnmarshalMap(resp.Attributes, &heartCount)
	if err != nil {
		return nil, err
	}
	return &heartCount, nil
}

func (d *DynamoClient) GetProto(table string, keys map[string]string, p proto.Message) error {
	dynamoKeys := make(map[string]dynamodb.AttributeValue)
	var keyErr error
	var k *dynamodb.AttributeValue
	for keyName, key := range keys {
		k, keyErr = dynamodbattribute.Marshal(&key)
		dynamoKeys[keyName] = *k
	}
	if keyErr != nil {
		glog.Errorf("Error marshalling key to attribute: %s", keyErr)
		return keyErr
	}
	req := d.client.GetItemRequest(&dynamodb.GetItemInput{
		TableName: &table,
		Key:       dynamoKeys})
	res, err := req.Send(context.Background())
	if err != nil {
		glog.Errorf("Error sending get request for %s %s", reflect.TypeOf(p), err)
		return err
	}
	err = dynamodbattribute.UnmarshalMap(res.Item, p)
	if err != nil {
		glog.Errorf("Error unmarshalling response into %s %s", reflect.TypeOf(p), err)
		return err
	}
	glog.Infof("Succesfully Got %s", reflect.TypeOf(p))
	return nil
}

func (d *DynamoClient) PutProtoBatch(table *string, protos []proto.Message) error {
	reqs := make([]dynamodb.WriteRequest, 0)
	for _, p := range protos {
		av, err := dynamodbattribute.MarshalMap(&p)
		if err != nil {
			return err
		}
		reqs = append(reqs, dynamodb.WriteRequest{PutRequest: &dynamodb.PutRequest{Item: av}})
	}
	numBatches := int(math.Ceil(float64(len(reqs)) / 25.0))
	currentBatch := 0
	for len(reqs) > 0 {
		// Take last 25 reqs (or all if <25 left)
		// Dynamo db restricts batch calls to 25 or fewer items
		startIdx := len(reqs) - 25
		if startIdx < 0 {
			startIdx = 0
		}
		req := d.client.BatchWriteItemRequest(&dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]dynamodb.WriteRequest{
				*table: reqs[startIdx:]}})
		_, err := req.Send(context.Background())
		if err != nil {
			glog.Errorf("Error batch putting %s %s", reflect.TypeOf(protos), err)
			glog.Errorf("Retrying...")
			time.Sleep(time.Second)
			continue
		}
		reqs = reqs[:startIdx]
		glog.Infof("Batch Put %s (%d/%d): %d Items Remaining", *table, currentBatch, numBatches, len(reqs))
		currentBatch += 1
	}
	glog.Infof("Successful Batch Put %s", reflect.TypeOf(protos))
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
