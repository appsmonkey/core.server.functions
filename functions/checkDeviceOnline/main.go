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

	// Fetch live data for defined period
	from := time.Now().Add(-time.Hour * 2).Unix()
	fmt.Println(from)
	// for _, d := range activeDevices {
	// 	if int64(d.Timestamp) < from {
	// 		d.Active = false
	// 	}
	// }

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
