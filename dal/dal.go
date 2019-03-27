package dal

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// QueryResult holds the result of a single row result
type QueryResult struct {
	items *dynamodb.QueryOutput
}

// Unmarshal the QUERY result into your type
// `Make sure that the `*out*` parameter is a ptr to a slice
func (r *QueryResult) Unmarshal(out interface{}) error {
	err := dynamodbattribute.UnmarshalListOfMaps(r.items.Items, &out)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

// GetResult holds the result of a single row result
type GetResult struct {
	item  *dynamodb.GetItemOutput
	key   map[string]*dynamodb.AttributeValue
	table string
	svc   *dynamodb.DynamoDB
}

// Unmarshal the GET result into your type
func (r *GetResult) Unmarshal(out interface{}) error {
	// if no data returned we will not try to do anything, `out` will remain unchanged
	if r.item == nil || len(r.item.Item) == 0 {
		return nil
	}

	err := dynamodbattribute.UnmarshalMap(r.item.Item, &out)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

// Update an item
func (r *GetResult) Update(updateExpression string, data Query) error {
	// Update Item in table and return
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: data,
		TableName:                 aws.String(r.table),
		Key:                       r.key,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String(updateExpression),
	}

	_, err := r.svc.UpdateItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

// ListResult holds the result of a teh List request
type ListResult struct {
	items *dynamodb.ScanOutput
}

// Unmarshal the GET result into your type
// `Make sure that the `*out*` parameter is a ptr to a slice
func (r *ListResult) Unmarshal(out interface{}) error {
	err := dynamodbattribute.UnmarshalListOfMaps(r.items.Items, &out)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

// AttributeValue for query parameters
type AttributeValue = dynamodb.AttributeValue

// Query definition
type Query = map[string]*AttributeValue

// Condition definition
type Condition = map[string]*dynamodb.Condition

// NameBuilder for a list of desired named parameters
type NameBuilder = expression.NameBuilder

// ValueBuilder for a list of desired value parameters
type ValueBuilder = expression.ValueBuilder

// ConditionBuilder is the filter for our scan
type ConditionBuilder = expression.ConditionBuilder

// ProjectionBuilder is the desired result we expect from the scan
type ProjectionBuilder = expression.ProjectionBuilder

// Projection will formulate a list of names to be returned
func Projection(name NameBuilder, names ...NameBuilder) ProjectionBuilder {
	return name.NamesList(names...)
}

// Name will return a new name builder to formulate a filter expression
func Name(name string) NameBuilder {
	return expression.Name(name)
}

// Value will return a new name builder to formulate a filter expression
func Value(name interface{}) ValueBuilder {
	return expression.Value(name)
}

// Insert the data into the DynamoDB table (Single Item)
func Insert(table string, data interface{}) error {
	// Create the dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Marshall the Item into a Map DynamoDB can deal with
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		fmt.Println("got error marshalling map:")
		fmt.Println(err.Error())
		return err
	}

	// Create Item in table and return
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}
	pOut, err := svc.PutItem(input)
	fmt.Println("saving item output:", pOut)
	if err != nil {
		fmt.Println("got error saving item:", err)
		return err
	}

	return err
}

// Get data from the table (Single Item)
func Get(table string, query Query) (*GetResult, error) {
	// Create the dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Perform the query
	fmt.Println("Trying to read from table: ", table)
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key:       query,
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &GetResult{item: result, table: table, key: query, svc: svc}, err
}

// GetFromIndex data from the table (Single Item)
func GetFromIndex(table, index string, condition Condition) (*QueryResult, error) {
	// Create the dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Perform the query
	fmt.Println("Trying to read from table: ", table)
	var queryInput = &dynamodb.QueryInput{
		TableName:     aws.String(table),
		IndexName:     aws.String(index),
		KeyConditions: condition,
	}

	result, err := svc.Query(queryInput)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &QueryResult{items: result}, err
}

// List data (returns possible multiple values)
func List(table string, filter ConditionBuilder, projection ProjectionBuilder) (*ListResult, error) {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		fmt.Println("got error building expression:")
		fmt.Println(err.Error())

		return nil, err
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(table),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	fmt.Println("Result", result)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		return nil, err
	}

	return &ListResult{items: result}, nil
}

// ListNoFilter data (returns possible multiple values)
func ListNoFilter(table string, projection ProjectionBuilder) (*ListResult, error) {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	expr, err := expression.NewBuilder().WithProjection(projection).Build()
	if err != nil {
		fmt.Println("Got error building expression:")
		fmt.Println(err.Error())

		return nil, err
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(table),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	fmt.Println("Result", result)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		return nil, err
	}

	return &ListResult{items: result}, nil
}

// Update an item
func Update(table, updateExpression string, key, data Query) error {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Update Item in table and return
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: data,
		TableName:                 aws.String(table),
		Key:                       key,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String(updateExpression),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
