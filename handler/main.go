package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

var (
	ErrMissingEnv = errors.New("SLACK_USER_SECRET and SLACK_AUTH_TOKEN environment variables cannot be empty")
)

func Handler(ctx context.Context) (Response, error) {
	fmt.Println("main.Handler")

	slackUserSecret := os.Getenv("SLACK_USER_SECRET")
	sleckToken := os.Getenv("SLACK_AUTH_TOKEN")
	fmt.Printf("Configuration loaded: \n\t SLACK_USER_SECRET=%s \n\t SLACK_AUTH_TOKEN=%s", slackUserSecret, sleckToken)

	if len(slackUserSecret) == 0 || len(sleckToken) == 0 {
		return Response{}, ErrMissingEnv
	}

	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"message": "Go Serverless v1.0! Your function executed successfully!",
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
