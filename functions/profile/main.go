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
	request := new(vm.ProfileRequest)

	response := request.Validate(req.Body)
	if response.Code != 0 {
		return response, nil
	}

	data, err := cog.Profile(request.Email)
	if err != nil {
		errData := es.ErrProfileMissingEmail
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)
		return response, nil
	}

	usrGroups, err := cog.ListGroupsForUser(request.Email)
	fmt.Println("USER GROUPS::", usrGroups)
	if err != nil {
		errData := es.ErrProfileMissingEmail
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)
		return response, nil
	}

	// TODO: check what data do we need from user's profile.
	response.Data = data
	response.Groups = usrGroups
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
	if len(os.Args) != 2 {
		fmt.Println(`
			ERROR: missing or having extra params.\n
			HINT: required params are: email
		`)
		return
	}

	data, _ := json.Marshal(vm.ProfileRequest{
		Email: os.Args[1],
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
