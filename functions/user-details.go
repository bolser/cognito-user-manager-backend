package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"net/http"
)

type Request struct {
	UserPoolId string `json:"user_pool_id"`
	Username   string `json:"username"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var r Request
	json.Unmarshal([]byte(request.Body), &r)

	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-west-1")}))

	cip := cognitoidentityprovider.New(sess)
	result, err := cip.AdminGetUser(&cognitoidentityprovider.AdminGetUserInput{Username: aws.String(r.Username), UserPoolId: aws.String(r.UserPoolId)})
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: http.StatusInternalServerError}, nil
	}
	user, _ := json.Marshal(result)
	return events.APIGatewayProxyResponse{Body: string(user), Headers: map[string]string{"Access-Control-Allow-Origin": "*"}, StatusCode: http.StatusOK}, nil
}

func main() {
	lambda.Start(Handler)
}
