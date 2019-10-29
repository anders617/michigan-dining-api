package dynamoclient

import (
	"context"
	"errors"

	pb "github.com/MichiganDiningAPI/api/proto"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/expression"
	"github.com/golang/glog"
)

func (d *DynamoClient) QueryDiningHalls() (*pb.DiningHalls, error) {
	params := &dynamodb.ScanInput{
		TableName: aws.String(DiningHallsTableName),
	}
	// Make the DynamoDB Query API call
	req := d.client.ScanRequest(params)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	diningHalls := pb.DiningHalls{}
	for _, i := range result.Items {
		dh := pb.DiningHall{}
		dynamodbattribute.UnmarshalMap(i, &dh)
		diningHalls.DiningHalls = append(diningHalls.DiningHalls, &dh)
	}
	return &diningHalls, nil
}

func (d *DynamoClient) QueryMenus(diningHallName *string, date *string, meal *string) (*[]*pb.Menu, error) {
	glog.Infof("QueryMenus %v, %v, %v", diningHallName, date, meal)
	if diningHallName != nil && date != nil && meal != nil {
		menu := pb.Menu{}
		err := d.GetProto(MenuTableName, map[string]string{DateKey: *date, MenuTableDiningHallMealKey: *diningHallName + *meal}, &menu)
		if err != nil {
			return nil, err
		}
		menus := []*pb.Menu{&menu}
		return &menus, nil
	}
	if date != nil && meal == nil {
		keyCond := expression.Key(DateKey).Equal(expression.Value(*date))
		if diningHallName != nil {
			keyCond = keyCond.And(expression.Key(MenuTableDiningHallMealKey).BeginsWith(*diningHallName))
		}
		expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()
		params := &dynamodb.QueryInput{
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			TableName:                 aws.String(MenuTableName),
		}
		req := d.client.QueryRequest(params)
		result, err := req.Send(context.Background())
		if err != nil {
			return nil, err
		}
		menus := make([]*pb.Menu, len(result.Items))
		for idx, item := range result.Items {
			menu := pb.Menu{}
			dynamodbattribute.UnmarshalMap(item, &menu)
			menus[idx] = &menu
		}
		return &menus, nil
	}
	if date != nil && meal != nil && diningHallName == nil {
		keyCond := expression.Key(DateKey).Equal(expression.Value(*date))
		filter := expression.Name("meal").Equal(expression.Value(*meal))
		expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).WithFilter(filter).Build()
		params := &dynamodb.QueryInput{
			KeyConditionExpression:    expr.KeyCondition(),
			FilterExpression:          expr.Filter(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			TableName:                 aws.String(MenuTableName),
		}
		req := d.client.QueryRequest(params)
		result, err := req.Send(context.Background())
		if err != nil {
			return nil, err
		}
		menus := make([]*pb.Menu, len(result.Items))
		for idx, item := range result.Items {
			menu := pb.Menu{}
			dynamodbattribute.UnmarshalMap(item, &menu)
			menus[idx] = &menu
		}
		return &menus, nil
	}
	return nil, errors.New("Unimplemented Menu Query")
}
