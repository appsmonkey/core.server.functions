package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	"github.com/appsmonkey/core.server.functions/integration/cognito"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

var (
	cog *cognito.Cognito
)

// Handler will handle our request comming from the API gateway
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (interface{}, error) {
	request := new(vm.SigninRequest)

	response := request.Validate(req.Body)
	if response.Code != 0 {
		return response, nil
	}

	data, err := cog.SignIn(request.Email, request.Password)
	if err != nil {
		errData := es.ErrRegistrationSignInError
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)
		return response, nil
	}

	response.Data = data
	return response, nil
}

func init() {
	if os.Getenv("ENV") == "local" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("error loading .env: %v\n", err)
		}
	}

	cog = cognito.NewCognito()
}

func local() {
	if len(os.Args) != 3 {
		fmt.Println(`
			ERROR: missing or having extra params.\n
			HINT: required params are: email password
		`)
		return
	}

	data, _ := json.Marshal(vm.SigninRequest{
		Email:    os.Args[1],
		Password: os.Args[2],
	})

	resp, err := Handler(context.Background(), events.APIGatewayProxyRequest{
		Body: string(data),
	})

	if err != nil {
		fmt.Printf("unhandled error! \nError: %v\n", err)
	} else {
		j, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Println(string(j))
	}
}

func main() {
	if os.Getenv("ENV") == "local" {
		local()
		return
	}

	lambda.Start(Handler)
}
