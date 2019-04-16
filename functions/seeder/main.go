package main

import (
	s "github.com/appsmonkey/core.server.functions/dal/seed"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler will handle our request comming from the API gateway
func Handler(req interface{}) error {
	seeder := req.(map[string]interface{})

	seed := seeder["seed"].(string)
	s.Run(seed)

	return nil
}

func main() {
	lambda.Start(Handler)
}
