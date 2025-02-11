package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"example.com/m/config"
	"example.com/m/domain/model"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/go-playground/validator"
)

var validate *validator.Validate

// コールドスタート時のみ実行
func init() {
	log.Println("Init function executed (cold start).")
	validate = validator.New()
}

func fetchPostalCodeInfo(postalCode string) (map[string]interface{}, error) {
	// 接続先のURLを環境変数から取得
	baseUrl := config.ZIPCLOUD_API_URL

	url := fmt.Sprintf("%s?zipcode=%s", baseUrl, postalCode)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch postal code info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	return result, nil
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

	log.Printf("Lambda Request ID: %s", lc.AwsRequestID)
	log.Printf("環境識別子： %s", config.ENV)

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
			Body:       `{"error": "Validation error"}`,
		}, nil
	}

	// 郵便番号検索APIの呼び出し
	postalInfo, err := fetchPostalCodeInfo(requestBody.PostalCode)
	if err != nil {
		log.Printf("Error fetching postal info: %v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		}, nil
	}

	// レスポンス作成
	responseBody, _ := json.Marshal(postalInfo)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseBody),
	}, nil
}

func main() {
	lambda.Start(Handler)
}
