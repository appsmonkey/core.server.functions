package main

import (
	"fmt"
	"strings"

	"github.com/appsmonkey/core.server.functions/dal"
	es "github.com/appsmonkey/core.server.functions/errorStatuses"
	"github.com/appsmonkey/core.server.functions/integration/cognito"
	m "github.com/appsmonkey/core.server.functions/models"
	"github.com/appsmonkey/core.server.functions/tools/defaultDevice"
	h "github.com/appsmonkey/core.server.functions/tools/helper"
	vm "github.com/appsmonkey/core.server.functions/viewmodels"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var (
	cog *cognito.Cognito
)

// Handler will handle our request comming from the API gateway
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := new(vm.DeviceListResponse)
	response.Init()
	cognitoID := h.CognitoIDZeroValue
	userName := ""
	authHdr := header("AccessToken", req.Headers)
	if len(authHdr) > 0 {
		c, u, isExpired, err := cog.ValidateToken(authHdr)
		if err != nil {
			fmt.Println(err)
			if isExpired {
				return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 401, Headers: response.Headers()}, nil
			}
		} else {
			cognitoID = c
			userName = u
		}
	}

	City := req.QueryStringParameters["city"]
	if len(City) < 1 {
		// Sarajevo is the default city
		City = "Sarajevo"
	}

	if cognitoID == h.CognitoIDZeroValue {
		rd := make([]*vm.DeviceGetData, 0)

		// Add the default device on the top
		dd := defaultDevice.Get(City)
		rd = append(rd, &dd)

		response.Data = rd
		return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
	}
	var userGroupsRes *cognitoidentityprovider.AdminListGroupsForUserOutput
	if len(userName) > 0 {
		groupRes, err := cog.ListGroupsForUser(userName)

		if err != nil {
			fmt.Println("User groups error ::: ", err)
			return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
		}
		userGroupsRes = groupRes
	}

	isAdmin := false
	for _, g := range userGroupsRes.Groups {
		if g.GroupName != nil && (*g.GroupName == "AdminGroup" || *g.GroupName == "SuperAdminGroup") {
			isAdmin = true
		}
	}
	fmt.Println("IS ADMIN ::: ", isAdmin)

	dbData := make([]m.Device, 0)

	if isAdmin {
		res, err := dal.ListNoFilter("devices", dal.Projection(
			dal.Name("token"),
			dal.Name("meta"),
			dal.Name("cognito_id"),
			dal.Name("active"),
			dal.Name("model"),
			dal.Name("indoor"),
			dal.Name("default_device"),
			dal.Name("map_meta"),
			dal.Name("latest"),
			dal.Name("timestamp"),
			dal.Name("zone_id"),
		))

		if err != nil {
			fmt.Println(err)
			response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
			return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
		}

		err = res.Unmarshal(&dbData)
		if err != nil {
			fmt.Println(err)
			response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
			return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
		}

	} else {
		res, err := dal.GetFromIndex("devices", "CognitoID-index", dal.Condition{
			"cognito_id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dal.AttributeValue{
					{
						S: aws.String(cognitoID),
					},
				},
			},
		})
		if err != nil {
			fmt.Println(err)
			response.AddError(&es.Error{Message: err.Error(), Data: "could not fetch data from the DB"})
			return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
		}

		err = res.Unmarshal(&dbData)
		if err != nil {
			fmt.Println(err)
			response.AddError(&es.Error{Message: err.Error(), Data: "could not unmarshal data from the DB"})
			return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 500, Headers: response.Headers()}, nil
		}
	}

	rd := make([]*vm.DeviceGetData, 0)

	// Add the default device on the top
	dd := defaultDevice.Get(City)
	rd = append(rd, &dd)

	for _, d := range dbData {
		owner := ""
		if isAdmin {
			res, err := dal.GetFromIndex("users", "CognitoID-index", dal.Condition{
				"cognito_id": {
					ComparisonOperator: aws.String("EQ"),
					AttributeValueList: []*dal.AttributeValue{
						{
							S: aws.String(d.CognitoID),
						},
					},
				},
			})

			if err != nil {
				fmt.Println("Failed to fetch user: ", err, d.CognitoID)
				owner = ""
			}

			owners := make([]m.User, 0)
			err = res.Unmarshal(&owners)

			if err != nil {
				fmt.Println("Failed to unmarshal user: ", err, d.CognitoID)
				owner = ""
			}

			if len(owners) > 0 {
				owner = owners[0].Email
			}
		} else {
			owner = userName
		}

		data := vm.DeviceGetData{
			DeviceID:  d.Token,
			Name:      d.Meta.Name,
			Active:    d.Active,
			Mine:      d.CognitoID == cognitoID,
			Model:     d.Meta.Model,
			Indoor:    d.Meta.Indoor,
			Location:  d.Meta.Coordinates,
			MapMeta:   d.MapMeta,
			Timestamp: d.Timestamp,
			Owner:     owner,
		}
		rd = append(rd, &data)
	}

	response.Data = rd

	return events.APIGatewayProxyResponse{Body: response.Marshal(), StatusCode: 200, Headers: response.Headers()}, nil
}

// CognitoData for user
func CognitoData(in map[string]interface{}) string {
	data := in["claims"].(map[string]interface{})

	return data["sub"].(string)
}

func header(hdr string, in map[string]string) string {
	result, ok := in[hdr]
	if !ok {
		lwr := strings.ToLower(hdr)
		result = in[lwr]
	}

	return result
}

func init() {
	cog = cognito.NewCognito()
}

func main() {
	lambda.Start(Handler)
}
