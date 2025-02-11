package main

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestHandler_ValidRequest(t *testing.T) {
	mockRequest := events.APIGatewayV2HTTPRequest{
		Body:    `{"postalCode": "1000001"}`,
		Headers: map[string]string{"Content-Type": "application/json"},
	}

	lambdaCtx := lambdacontext.LambdaContext{AwsRequestID: "mock-request-id"}
	ctx := lambdacontext.NewContext(context.Background(), &lambdaCtx)

	response, err := Handler(ctx, mockRequest)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody map[string]interface{}
	err = json.Unmarshal([]byte(response.Body), &responseBody)
	assert.NoError(t, err)
	assert.Contains(t, responseBody, "results")
}

func TestHandler_InvalidRequestBody(t *testing.T) {
	mockRequest := events.APIGatewayV2HTTPRequest{
		Body:    "invalid-json",
		Headers: map[string]string{"Content-Type": "application/json"},
	}

	lambdaCtx := lambdacontext.LambdaContext{AwsRequestID: "mock-request-id"}
	ctx := lambdacontext.NewContext(context.Background(), &lambdaCtx)

	response, err := Handler(ctx, mockRequest)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	var responseBody map[string]interface{}
	err = json.Unmarshal([]byte(response.Body), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request body", responseBody["error"])
}

func TestHandler_ValidationError(t *testing.T) {
	mockRequest := events.APIGatewayV2HTTPRequest{
		Body:    `{"postalCode": ""}`,
		Headers: map[string]string{"Content-Type": "application/json"},
	}

	lambdaCtx := lambdacontext.LambdaContext{AwsRequestID: "mock-request-id"}
	ctx := lambdacontext.NewContext(context.Background(), &lambdaCtx)

	response, err := Handler(ctx, mockRequest)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	var responseBody map[string]interface{}
	err = json.Unmarshal([]byte(response.Body), &responseBody)
	assert.NoError(t, err)
	assert.Contains(t, responseBody["error"], "Validation error")
}

func TestHandler_APIFailure(t *testing.T) {
	mockRequest := events.APIGatewayV2HTTPRequest{
		Body:    `{"postalCode": "9999999"}`,
		Headers: map[string]string{"Content-Type": "application/json"},
	}

	lambdaCtx := lambdacontext.LambdaContext{AwsRequestID: "mock-request-id"}
	ctx := lambdacontext.NewContext(context.Background(), &lambdaCtx)

	// Simulate API failure
	response, err := Handler(ctx, mockRequest)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)

	var responseBody map[string]interface{}
	err = json.Unmarshal([]byte(response.Body), &responseBody)
	assert.NoError(t, err)
	assert.Contains(t, responseBody["error"], "failed to fetch postal code info")
}
