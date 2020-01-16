package main

import (
	"fmt"
	"net/http"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	m "github.com/appsmonkey/core.server.functions/models"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/avct/uasurfer"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("VERIFICATION REQUEST: ", req.QueryStringParameters)
	request := new(vm.VerifyRedirectRequest)
	response := request.Validate(req.QueryStringParameters)
	if response.Code != 0 {
		fmt.Printf("errors on request: %v, requestID: %v", response.Errors, response.RequestID)

		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	ua := uasurfer.Parse(req.Headers["User-Agent"])
	fmt.Println("REQUEST_TYPE ::: ", request.Type)

	// create verification URL
	verificationURL := "https://cityos.auth.us-east-1.amazoncognito.com/confirmUser?client_id=" + request.ClientID + "&user_name=" + request.UserName + "&response_type=code" + "&confirmation_code=" + request.ConfirmationCode
	fmt.Println("VERIFICATION URL:", verificationURL)

	verificationResponse, err := http.Get(verificationURL)
	fmt.Println("Verification response ::: ", verificationResponse.StatusCode)

	if err != nil {
		fmt.Println("Verification error: ", err)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	res, err := dal.Get("users", map[string]*dal.AttributeValue{
		"cognito_id": {
			S: aws.String(request.CognitoID),
		},
		"email": {
			S: aws.String(request.UserName),
		},
	})

	if err != nil {
		fmt.Println("User missing error: ", err)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	user := new(m.User)
	res.Unmarshal(&user)

	fmt.Println("USER :::", user)
	if user.Attributes["cognito:user_status"] != "CONFIRMED" {
		fmt.Println("User not confirmed, verification failed.")
		errData := es.VerificationFailed
		errData.Data = err.Error()
		response.Errors = append(response.Errors, errData)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 403, Headers: response.Headers()}, nil
	}

	response.Data = user

	headers := response.Headers()

	if ua.OS.Name.String() == "OSAndroid" {
		return events.APIGatewayProxyResponse{Body: `[{
			"relation": ["delegate_permission/common.handle_all_urls"],
			"target": {
			  "namespace": "android_app",
			  "package_name": "io.cityos.cityosair",
			  "sha256_cert_fingerprints":
			  ["66:6D:4F:2F:AA:94:E4:77:C1:57:EB:95:8E:58:DF:42:60:9D:92:34:3E:F8:B0:D9:7F:D6:25:2F:2A:95:9B:EC"]
			}
		  }]`, StatusCode: 200, Headers: headers}, nil
	} else if ua.OS.Name.String() == "OSiOS" {
		// headers["Location"] = "http://links.cityos.io/.well-known/apple-app-site-association"
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: headers}, nil
	} else {
		// headers["Location"] = "https://dev.cityos.io"
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: headers}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
