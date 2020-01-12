package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

	// create verification URL
	verificationURL := "https://cityos.auth.us-east-1.amazoncognito.com/confirmUser?client_id=" + request.ClientID + "&user_name=" + request.UserName + "&response_type=code" + "&confirmation_code=" + request.ConfirmationCode
	fmt.Println("VERIFICATION URL:", verificationURL)

	_, err := http.Get(verificationURL)
	if err != nil {
		fmt.Println("Verification error: ", err)
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}

	// asset links for android
	var assetLinks map[string]interface{}
	json.Unmarshal([]byte(`[{
		"relation": ["delegate_permission/common.handle_all_urls"],
		"target": {
		  "namespace": "android_app",
		  "package_name": "io.cityos.cityosair",
		  "sha256_cert_fingerprints":
		  ["66:6D:4F:2F:AA:94:E4:77:C1:57:EB:95:8E:58:DF:42:60:9D:92:34:3E:F8:B0:D9:7F:D6:25:2F:2A:95:9B:EC"]
		}
	  }]`), &assetLinks)

	headers := response.Headers()

	json, err := json.Marshal(assetLinks)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	}
	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: 200, Headers: headers}, nil

	// // redired URL - default
	// switch request.ClientID {
	// // - Android
	// case "km0afsc8ua4f0bc56brcn7t90":
	// 	json, err := json.Marshal(assetLinks)
	// 	if err != nil {
	// 		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 400, Headers: response.Headers()}, nil
	// 	}
	// 	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: 200, Headers: headers}, nil
	// // - IOS
	// case "70mq6uphtmmorkjt74ei0rj5fr":
	// 	// - TODO: change URL for IOS
	// 	headers["Location"] = "https://dev.cityos.io"
	// 	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 302, Headers: headers}, nil
	// default:
	// 	headers["Location"] = "https://dev.cityos.io"
	// 	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 302, Headers: headers}, nil
	// }

}

func main() {
	lambda.Start(Handler)
}
