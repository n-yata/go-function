package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"example.com/m/domain/model"
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
	mockRequest := events.APIGatewayProxyRequest{
		Body:    `{"postalCode": "1000001"}`,
		Headers: map[string]string{"Content-Type": "application/json"},
	}

	lambdaCtx := lambdacontext.LambdaContext{AwsRequestID: "mock-request-id"}
	ctx := lambdacontext.NewContext(context.Background(), &lambdaCtx)

	response, err := Handler(ctx, mockRequest)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody model.ResponseBody
	err = json.Unmarshal([]byte(response.Body), &responseBody)
	log.Println(responseBody)
	assert.NoError(t, err)

	// "results" ではなく、直接フィールドを確認する
	assert.NotEmpty(t, responseBody.Address1)
	assert.NotEmpty(t, responseBody.Address2)
	assert.NotEmpty(t, responseBody.Address3)
}

func TestHandler_InvalidRequestBody(t *testing.T) {
	mockRequest := events.APIGatewayProxyRequest{
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
	mockRequest := events.APIGatewayProxyRequest{
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
	mockRequest := events.APIGatewayProxyRequest{
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
