package main

import (
	"fmt"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

type resultData struct {
	Date  float64 `json:"date"`
	Value float64 `json:"value"`
}

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := new(vm.ChartLiveDeviceRequest)
	response := request.Validate(req.QueryStringParameters)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	res, err := dal.QueryMultiple("live",
		dal.Condition{
			"token": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dal.AttributeValue{
					{
						S: aws.String(request.Token),
					},
				},
			},
			"timestamp": {
				ComparisonOperator: aws.String("GT"),
				AttributeValueList: []*dal.AttributeValue{
					{
						N: aws.String(request.From),
					},
				},
			},
		},
		dal.Projection(dal.Name("timestamp"), dal.Name(request.Sensor)),
		true)

	if err != nil {
		response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	var dbData []map[string]float64
	err = res.Unmarshal(&dbData)
	if err != nil {
		response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
	}

	result := make([]*resultData, 0)
	for _, v := range dbData {
		result = append(result, &resultData{
			Date:  v["timestamp"],
			Value: v[request.Sensor],
		})
	}

	response.Data = result

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

func main() {
	lambda.Start(Handler)
}
