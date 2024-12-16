package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/assert"
)

func TestHandler_ValidRequest(t *testing.T) {
	// モックリクエストの作成
	mockRequest := events.APIGatewayV2HTTPRequest{
		Body: `{"name": "hello", "age": 5}`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// モックLambdaコンテキストを作成
	lambdaCtx := lambdacontext.LambdaContext{
		AwsRequestID: "mock-request-id",
	}
	ctx := lambdacontext.NewContext(context.Background(), &lambdaCtx)

	// ハンドラーの呼び出し
	response, err := Handler(ctx, mockRequest)

	// エラーチェック
	assert.NoError(t, err, "Handler should not return an error")

	// ステータスコードの確認
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected status code to be 200")

	// レスポンスボディの確認
	var responseBody map[string]interface{}
	err = json.Unmarshal([]byte(response.Body), &responseBody)
	log.Printf("Response body: %+v", responseBody)
	assert.NoError(t, err, "Response body should be valid JSON")
	assert.Equal(t, "Validation passed", responseBody["message"], "Response message should be 'Validation passed'")
	assert.Equal(t, "hello", responseBody["name"], "Response 'name' should match request")
}

func TestHandler_InvalidRequestBody(t *testing.T) {
	// 無効なリクエストボディ
	mockRequest := events.APIGatewayV2HTTPRequest{
		Body: `invalid-json`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// モックLambdaコンテキストを作成
	lambdaCtx := lambdacontext.LambdaContext{
		AwsRequestID: "mock-request-id",
	}
	ctx := lambdacontext.NewContext(context.Background(), &lambdaCtx)

	// ハンドラーの呼び出し
	response, err := Handler(ctx, mockRequest)

	// エラーチェック
	assert.NoError(t, err, "Handler should not return an error")

	// ステータスコードの確認
	assert.Equal(t, http.StatusBadRequest, response.StatusCode, "Expected status code to be 400")

	// エラーメッセージの確認
	var responseBody map[string]interface{}
	err = json.Unmarshal([]byte(response.Body), &responseBody)
	assert.NoError(t, err, "Response body should be valid JSON")
	assert.Equal(t, "Invalid request body", responseBody["error"], "Error message should indicate invalid request body")
}

func TestHandler_ValidationError(t *testing.T) {
	// バリデーションエラーを含むリクエストボディ
	mockRequest := events.APIGatewayV2HTTPRequest{
		Body: `{"name": "", "age": -1}`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// モックLambdaコンテキストを作成
	lambdaCtx := lambdacontext.LambdaContext{
		AwsRequestID: "mock-request-id",
	}
	ctx := lambdacontext.NewContext(context.Background(), &lambdaCtx)

	// ハンドラーの呼び出し
	response, err := Handler(ctx, mockRequest)

	// エラーチェック
	assert.NoError(t, err, "Handler should not return an error")

	// ステータスコードの確認
	assert.Equal(t, http.StatusBadRequest, response.StatusCode, "Expected status code to be 400")

	// エラーメッセージの確認
	var responseBody map[string]interface{}
	err = json.Unmarshal([]byte(response.Body), &responseBody)
	assert.NoError(t, err, "Response body should be valid JSON")
	assert.Contains(t, responseBody["error"], "Validation error", "Error message should indicate validation error")
}

func TestHandler_MissingField(t *testing.T) {
	// 必須フィールドが欠けているリクエスト
	mockRequest := events.APIGatewayV2HTTPRequest{
		Body: `{"age": 10}`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// モックLambdaコンテキストを作成
	lambdaCtx := lambdacontext.LambdaContext{
		AwsRequestID: "mock-request-id",
	}
	ctx := lambdacontext.NewContext(context.Background(), &lambdaCtx)

	// ハンドラーの呼び出し
	response, err := Handler(ctx, mockRequest)

	// エラーチェック
	assert.NoError(t, err, "Handler should not return an error")

	// ステータスコードの確認
	assert.Equal(t, http.StatusBadRequest, response.StatusCode, "Expected status code to be 400")

	// エラーメッセージの確認
	var responseBody map[string]interface{}
	err = json.Unmarshal([]byte(response.Body), &responseBody)
	assert.NoError(t, err, "Response body should be valid JSON")
	assert.Contains(t, responseBody["error"], "Validation error", "Error message should indicate missing field")
}
