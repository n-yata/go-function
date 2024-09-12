package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
)

type FuncRequest struct {
	URL     string            `json:"url,omitempty"`
	Method  string            `json:"method,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

type FuncResponse struct {
	Status string `json:"status"`
	Body   string `json:"text"`
}

func HandleLambdaEvent(event *FuncRequest) (*FuncResponse, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}

	var reqBody io.Reader
	if event.Body != "" {
		reqBody = strings.NewReader(event.Body)
	}

	req, _ := http.NewRequest(event.Method, event.URL, reqBody)
	for key, value := range event.Headers {
		req.Header.Set(key, value)
	}

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer res.Body.Close()

	byteArray, _ := io.ReadAll(res.Body)
	return &FuncResponse{Status: res.Status, Body: string(byteArray)}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
