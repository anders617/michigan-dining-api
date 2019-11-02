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

func (d *DynamoClient) ForEachFood(startDate *string, endDate *string, fn func(*pb.Food)) error {
	var filter expression.ConditionBuilder
	if startDate != nil && endDate != nil {
		filter = expression.Name("date").Between(expression.Value(*startDate), expression.Value(*endDate))
	} else if startDate != nil {
		filter = expression.Name("date").GreaterThanEqual(expression.Value(*startDate))
	} else if endDate != nil {
		filter = expression.Name("date").LessThanEqual(expression.Value(*endDate))
	}
	params := &dynamodb.ScanInput{
		TableName: aws.String(FoodTableName),
	}
	if startDate != nil || endDate != nil {
		glog.Info(*startDate)
		expr, _ := expression.NewBuilder().WithFilter(filter).Build()
		params.FilterExpression = expr.Filter()
		params.ExpressionAttributeNames = expr.Names()
		params.ExpressionAttributeValues = expr.Values()
	}
	req := d.client.ScanRequest(params)
	p := dynamodb.NewScanPaginator(req)

	for p.Next(context.Background()) {
		page := p.CurrentPage()
		for _, item := range page.Items {
			food := pb.Food{}
			dynamodbattribute.UnmarshalMap(item, &food)
			fn(&food)
		}
	}

	if err := p.Err(); err != nil {
		return err
	}
	return nil
}

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

func (d *DynamoClient) QueryFoodsDateRange(name *string, startDate *string, endDate *string) (*[]*pb.Food, error) {
	glog.Infof("QueryFoodsDateRange %v %v %v", name, startDate, endDate)
	if name != nil {
		keyCond := expression.Key(FoodTableNameKey).Equal(expression.Value(*name))
		if startDate != nil && endDate != nil {
			keyCond = keyCond.And(expression.Key(DateKey).Between(expression.Value(*startDate), expression.Value(*endDate)))
		} else if startDate != nil {
			keyCond = keyCond.And(expression.Key(DateKey).GreaterThanEqual(expression.Value(*startDate)))
		} else if endDate != nil {
			keyCond = keyCond.And(expression.Key(DateKey).LessThanEqual(expression.Value(*endDate)))
		}
		expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()
		glog.Infof("Expr: %v", keyCond)
		params := &dynamodb.QueryInput{
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			TableName:                 aws.String(FoodTableName),
		}
		return d.queryFoods(params)
	}
	foods := make([]*pb.Food, 0)
	err := d.ForEachFood(startDate, endDate, func(food *pb.Food) {
		foods = append(foods, food)
	})
	glog.Info(len(foods))
	if err != nil {
		return nil, err
	}
	return &foods, nil
}

func (d *DynamoClient) QueryFoods(name *string, date *string) (*[]*pb.Food, error) {
	glog.Infof("QueryFoods %v %v", name, date)
	if name != nil && date != nil {
		food := pb.Food{}
		err := d.GetProto(FoodTableName, map[string]string{FoodTableNameKey: *name, DateKey: *date}, &food)
		if err != nil {
			return nil, err
		}
		return &[]*pb.Food{&food}, nil
	}
	if name != nil {
		keyCond := expression.Key(FoodTableNameKey).Equal(expression.Value(*name))
		expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()
		params := &dynamodb.QueryInput{
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			TableName:                 aws.String(FoodTableName),
		}
		return d.queryFoods(params)
	}
	return nil, errors.New("Unimplemented Foods Query")
}

func (d *DynamoClient) queryFoods(params *dynamodb.QueryInput) (*[]*pb.Food, error) {
	req := d.client.QueryRequest(params)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	foods := make([]*pb.Food, len(result.Items))
	for idx, item := range result.Items {
		food := pb.Food{}
		dynamodbattribute.UnmarshalMap(item, &food)
		foods[idx] = &food
	}
	return &foods, nil
}

func (d *DynamoClient) QueryMenus(diningHallName *string, date *string, meal *string) (*[]*pb.Menu, error) {
	glog.Infof("QueryMenus %v, %v, %v", diningHallName, date, meal)
	// If we have all three then we can just do a Get since the key is fully specified
	if diningHallName != nil && date != nil && meal != nil {
		menu := pb.Menu{}
		err := d.GetProto(MenuTableName, map[string]string{DateKey: *date, MenuTableDiningHallMealKey: *diningHallName + *meal}, &menu)
		if err != nil {
			return nil, err
		}
		menus := []*pb.Menu{&menu}
		return &menus, nil
	}
	// If we are missing the meal, then do a PartitionKey lookup with BeginsWith condition on DiningHallMealKey for diningHallName
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
		return d.queryMenus(params)
	}
	// If we are missing diningHallName do a PartitionKey lookup with filter expression for meal
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
		return d.queryMenus(params)
	}
	// If we are missing date, we have to do a scan with filter
	return nil, errors.New("Unimplemented Menu Query")
}

// Execute a query with the given parameters and marshal the output into a slice of *pb.Menu
func (d *DynamoClient) queryMenus(params *dynamodb.QueryInput) (*[]*pb.Menu, error) {
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
