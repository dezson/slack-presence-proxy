package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

var (
	ErrMissingEnv     = errors.New("Environment variable cannot be empty")
	ErrNon200Response = errors.New("Non 200 response found")
)

type slackPresenceResponse struct {
	Ok       string `json:"ok"`
	Presence string `json:"presence,omitempty"`
	Error    string `json:"error,omitempty"`
}

func GetUserPresence(userSecret string, token string) (string, error) {
	fmt.Println("main.GetUserPresence")
	if len(userSecret) == 0 || len(token) == 0 {
		return "", ErrMissingEnv
	}

	data := url.Values{}
	data.Add("token", token)
	data.Add("user", userSecret)
	fmt.Printf("Request Arguments: %s\n", data.Encode())

	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://slack.com/api/users.getPresence", bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	resp, _ := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		return "", ErrNon200Response
	}

	var apiResp slackPresenceResponse
	buf, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(buf, &apiResp)

	if len(apiResp.Error) != 0 {
		return "", errors.New(apiResp.Error)
	}

	fmt.Printf("Response body: %s\n", resp.Body)
	return apiResp.Presence, nil
}

func Handler(ctx context.Context) (Response, error) {
	fmt.Println("main.Handler")
	slackUserSecret := os.Getenv("SLACK_USER_SECRET")
	sleckToken := os.Getenv("SLACK_AUTH_TOKEN")
	fmt.Printf("Configuration loaded: \n\t SLACK_USER_SECRET=%s \n\t SLACK_AUTH_TOKEN=%s\n", slackUserSecret, sleckToken)

	userPresence, err := GetUserPresence(slackUserSecret, sleckToken)

	var message string
	if err != nil {
		message = err.Error()
	} else {
		message = fmt.Sprintf("Your friend is: %s", userPresence)
	}

	var buf bytes.Buffer
	body, err := json.Marshal(map[string]interface{}{
		"message": message,
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
