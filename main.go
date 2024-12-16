package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"example.com/m/domain/model"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
)

var validate *validator.Validate

// コールドスタート時のみ実行
func init() {
	log.Println("Init function executed (cold start).")
	validate = validator.New()
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Println("Handler function executed.")

	// LambdaのリクエストIDを取得
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		log.Println("Failed to retrieve Lambda context")
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Unable to retrieve Lambda context"}`,
		}, nil
	}

	requestID := lc.AwsRequestID
	log.Printf("Lambda Request ID: %s", requestID)

	// 環境ファイルの読み込み
	err := godotenv.Load()
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "fatal loading env file"}`,
		}, nil
	}
	host := os.Getenv("USER")
	pass := os.Getenv("PASSWORD")
	dbName := os.Getenv("DATABASE")
	log.Printf("host: %s, pass: %s, dbName: %s", host, pass, dbName)

	// リクエストボディのパース
	var requestBody model.RequestBody
	if err := json.Unmarshal([]byte(req.Body), &requestBody); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid request body"}`,
		}, nil
	}

	// バリデーション
	if err := validate.Struct(requestBody); err != nil {
		log.Printf("Validation failed: %v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Validation error: ` + err.Error() + `"}`,
		}, nil
	}

	// レスポンスの作成
	responseBody := model.ResponseBody{
		Message: "Validation passed",
		Name:    requestBody.Name,
		Age:     requestBody.Age,
	}

	response := events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"unique-key":   "unique-value",
		},
		Body: responseBody.ToJSON(),
	}

	return response, nil
}

func main() {
	lambda.Start(Handler)
}
