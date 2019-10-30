package dynamoclient

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	DiningHallsTableName = "DiningHalls"
	ItemsTableName       = "Items"
	MenuTableName        = "Menus"
	FoodTableName        = "Foods"
	FoodStatsTableName   = "FoodStats"
)

var (
	NameKey                    = "name"
	DiningHallDateMealKey      = "key"
	DateKey                    = "date"
	NameDateKey                = "key"
	FoodTableNameKey           = "key"
	MenuTableDiningHallMealKey = "diningHallMeal"
	FoodStatsDateKey           = "date"
)

var (
	TableNames = []string{
		DiningHallsTableName,
		ItemsTableName,
		MenuTableName,
		FoodTableName,
		FoodStatsTableName}
	TableKeys = map[string][]dynamodb.KeySchemaElement{
		DiningHallsTableName: []dynamodb.KeySchemaElement{
			dynamodb.KeySchemaElement{
				AttributeName: &NameKey,
				KeyType:       "HASH"}},
		ItemsTableName: []dynamodb.KeySchemaElement{
			dynamodb.KeySchemaElement{
				AttributeName: &NameKey,
				KeyType:       "HASH"}},
		MenuTableName: []dynamodb.KeySchemaElement{
			dynamodb.KeySchemaElement{
				AttributeName: &DateKey,
				KeyType:       "HASH"},
			dynamodb.KeySchemaElement{
				AttributeName: &MenuTableDiningHallMealKey,
				KeyType:       "RANGE"}},
		FoodTableName: []dynamodb.KeySchemaElement{
			dynamodb.KeySchemaElement{
				AttributeName: &FoodTableNameKey,
				KeyType:       "HASH"},
			dynamodb.KeySchemaElement{
				AttributeName: &DateKey,
				KeyType:       "RANGE"}},
		FoodStatsTableName: []dynamodb.KeySchemaElement{
			dynamodb.KeySchemaElement{
				AttributeName: &FoodStatsDateKey,
				KeyType:       "HASH",
			},
		}}
	TableAttributes = map[string][]dynamodb.AttributeDefinition{
		DiningHallsTableName: []dynamodb.AttributeDefinition{
			dynamodb.AttributeDefinition{
				AttributeName: &NameKey,
				AttributeType: dynamodb.ScalarAttributeTypeS}},
		ItemsTableName: []dynamodb.AttributeDefinition{
			dynamodb.AttributeDefinition{
				AttributeName: &NameKey,
				AttributeType: dynamodb.ScalarAttributeTypeS}},
		MenuTableName: []dynamodb.AttributeDefinition{
			dynamodb.AttributeDefinition{
				AttributeName: &DateKey,
				AttributeType: dynamodb.ScalarAttributeTypeS},
			dynamodb.AttributeDefinition{
				AttributeName: &MenuTableDiningHallMealKey,
				AttributeType: dynamodb.ScalarAttributeTypeS}},
		FoodTableName: []dynamodb.AttributeDefinition{
			dynamodb.AttributeDefinition{
				AttributeName: &FoodTableNameKey,
				AttributeType: dynamodb.ScalarAttributeTypeS},
			dynamodb.AttributeDefinition{
				AttributeName: &DateKey,
				AttributeType: dynamodb.ScalarAttributeTypeS}},
		FoodStatsTableName: []dynamodb.AttributeDefinition{
			dynamodb.AttributeDefinition{
				AttributeName: &FoodStatsDateKey,
				AttributeType: dynamodb.ScalarAttributeTypeS},
		}}
)
