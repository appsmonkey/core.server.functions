package main

import (
	"context"
	"fmt"
	"time"

	"github.com/appsmonkey/core.server.functions/dal"
	m "github.com/appsmonkey/core.server.functions/models"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	sl "github.com/aws/aws-sdk-go/service/lambda"
)

var lambdaClient *sl.Lambda

// Handler will handle our request comming from the API gateway
func Handler(ctx context.Context, req interface{}) error {

	// Fetch active devices
	projBuilder := dal.Projection(dal.Name("token"), dal.Name("active"), dal.Name("timestamp"))
	res, err := dal.List("devices", dal.Name("active").Equal(dal.Value(true)), projBuilder)

	if err != nil {
		fmt.Println("Fetching devices from device table failed", err)
		return err
	}

	activeDevices := make([]m.Device, 0)
	err = res.Unmarshal(&activeDevices)
	if err != nil {
		fmt.Println("Falied to unmarshal devices", err)
		return err
	}

	// activeState, err := dynamodbattribute.Marshal(false)

	if err != nil {
		fmt.Println("Failed to marshal active state")
		return nil
	}

	type schemaData struct {
		Version   string
		Heartbeat int
		Data      m.Schema
	}

	schemaRes, err := dal.Get("schema", map[string]*dal.AttributeValue{
		"version": {
			S: aws.String("1"),
		},
	})
	if err != nil {
		fmt.Println("Error fetching schema from db", err)
		return err
	}

	schema := new(schemaData)
	err = schemaRes.Unmarshal(schemaRes)
	if err != nil {
		fmt.Println("Error unmarshaling schema ::. ", err)
	}

	// 120 is deafult allowed timeout
	heartbeat := 120
	if schema.Heartbeat != 0 {
		fmt.Println("Setting heartbeat from schema", schema.Heartbeat)
		heartbeat = schema.Heartbeat
	}

	// Fetch live data for defined period
	from := time.Now().Add(-time.Minute * time.Duration(heartbeat)).Unix()
	for _, d := range activeDevices {
		if d.Timestamp < float64(from) {
			fmt.Println("Changing state of: ", d.Token, " - to offline")
			d.Active = false

			err = dal.Update("devices", "set active = :a",
				map[string]*dal.AttributeValue{
					"token": {
						S: aws.String(d.Token),
					},
				}, map[string]*dal.AttributeValue{
					":a": {
						BOOL: &d.Active,
					},
				})
		}
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}

func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	lambdaClient = sl.New(sess, &aws.Config{Region: aws.String("us-east-1")})
}
