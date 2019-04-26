package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/appsmonkey/core.server.functions/dal"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler will handle our request comming from the API gateway
func Handler(ctx context.Context, req interface{}) error {
	input, ok := req.(map[string]interface{})
	if !ok {
		err := errors.New("incorrect data received. input has incorrect format")
		fmt.Println(err)
		return err
	}

	dal.Insert(input["table"].(string), input["data"])

	return nil
}

func main() {
	lambda.Start(Handler)
}
