package dynamoclient

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	DiningHallsTableName = "DiningHalls"
	ItemsTableName       = "Items"
	MenuTableName        = "Menus"
	FoodTableName        = "Foods"
)

var (
	NameKey               = "name"
	DiningHallDateMealKey = "key"
	DateKey               = "date"
	NameDateKey           = "key"
	FoodTableNameKey      = "key"
)

var (
	TableNames = []string{
		DiningHallsTableName,
		ItemsTableName,
		MenuTableName,
		FoodTableName}
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
				AttributeName: &DiningHallDateMealKey,
				KeyType:       "HASH"},
			dynamodb.KeySchemaElement{
				AttributeName: &DateKey,
				KeyType:       "RANGE"}},
		FoodTableName: []dynamodb.KeySchemaElement{
			dynamodb.KeySchemaElement{
				AttributeName: &FoodTableNameKey,
				KeyType:       "HASH"},
			dynamodb.KeySchemaElement{
				AttributeName: &DateKey,
				KeyType:       "RANGE"}}}
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
				AttributeName: &DiningHallDateMealKey,
				AttributeType: dynamodb.ScalarAttributeTypeS},
			dynamodb.AttributeDefinition{
				AttributeName: &DateKey,
				AttributeType: dynamodb.ScalarAttributeTypeS}},
		FoodTableName: []dynamodb.AttributeDefinition{
			dynamodb.AttributeDefinition{
				AttributeName: &FoodTableNameKey,
				AttributeType: dynamodb.ScalarAttributeTypeS},
			dynamodb.AttributeDefinition{
				AttributeName: &DateKey,
				AttributeType: dynamodb.ScalarAttributeTypeS}}}
)
