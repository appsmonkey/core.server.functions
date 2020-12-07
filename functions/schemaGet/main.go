package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	s "github.com/appsmonkey/core.server.functions/models/schema"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := vm.SchemaInitResponse()
	version := "1"
	avHdr := req.Headers["Accept-Version"]
	if len(avHdr) > 0 {
		version = avHdr
	}

	var schemaTable = "schema"
	if value, ok := os.LookupEnv("dynamodb_table_schema"); ok {
		schemaTable = value
	}

	res, err := dal.Get(schemaTable, map[string]*dal.AttributeValue{
		"version": {
			S: aws.String(version),
		},
	})
	if err != nil {
		errData := es.ErrSchemaNotFound
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	type versionData struct {
		Version string   `json:"version"`
		Data    s.Schema `json:"data"`
	}

	model := versionData{}
	err = res.Unmarshal(&model)
	if err != nil {
		errData := es.ErrSchemaNotFound
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	response.Data = model.Data

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

func main() {
	lambda.Start(Handler)
}
