package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	RequestMethod  string `json:"RequestMethod"`
	RequestBody    string `json:"RequestBody"`
	PathParameter  string `json:"PathParameter"`
	QueryParameter string `json:"QueryParameter"`
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	method := req.HTTPMethod
	body := req.Body
	pathParam := req.PathParameters["pathparam"]
	queryParam := req.QueryStringParameters["queryparam"]

	res := Response{
		RequestMethod:  method,
		RequestBody:    body,
		PathParameter:  pathParam,
		QueryParameter: queryParam,
	}
	jsonBytes, _ := json.Marshal(res)

	return events.APIGatewayProxyResponse{
		Body:       string(jsonBytes),
		StatusCode: 200,
	}, nil
}
func main() {
	lambda.Start(handler)
}
