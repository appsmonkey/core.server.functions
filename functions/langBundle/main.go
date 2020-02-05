package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/appsmonkey/core.server.functions/integration/cognito"
	h "github.com/appsmonkey/core.server.functions/tools/helper"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	cog *cognito.Cognito
)

// Handler will handle our request comming from the API gateway
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := new(vm.LangBundleRequest)
	response := request.Validate(req.QueryStringParameters)

	res, err := h.GetLang(request.Language)
	fmt.Println(request.Language, "RES ::: ", res.Body)

	if err != nil {
		fmt.Println("Error fatching lang file")
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	if err != nil {
		fmt.Printf("%s", err)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	toReturn := new(map[string]string)
	err = json.Unmarshal(contents, &toReturn)
	if err != nil {
		fmt.Println("Unmarshaling er")
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	response.Data = toReturn
	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

func main() {
	if os.Getenv("ENV") == "local" {
		return
	}

	lambda.Start(Handler)
}
