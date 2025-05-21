package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

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

// APIリクエスト(POST)用のデータ構造体
type PostalCodeRequest struct {
	PostalCode string `json:"postalCode"`
}

func fetchPostalCodeInfo(postalCode string) (*model.ResponseBody, error) {
	// 接続先のURLを環境変数から取得
	baseUrl := config.ZIPCLOUD_API_URL
	url := fmt.Sprintf("%s?zipcode=%s", baseUrl, postalCode)

	// // リクエストボディをJSONに変換(POSTリクエスト)
	// requestBody := PostalCodeRequest{PostalCode: postalCode}
	// jsonBody, err := json.Marshal(requestBody)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to marshal request body: %w", err)
	// }

	// // HTTPリクエストの作成（POSTリクエスト）
	// req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create request: %w", err)
	// }

	// HTTPリクエストの作成
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := callApi(req)
	if err != nil {
		return nil, fmt.Errorf("failed callApi: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	// レスポンスボディを読み込む
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %w", err)
	}

	// APIのレスポンスをパース
	var apiResponse map[string]interface{}
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	// レスポンスデータをマッピング
	results, ok := apiResponse["results"].([]interface{})
	if !ok || len(results) == 0 {
		return nil, fmt.Errorf("no address found for postal code: %s", postalCode)
	}

	firstResult, ok := results[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid API response format")
	}

	response := &model.ResponseBody{
		Address1: firstResult["address1"].(string),
		Address2: firstResult["address2"].(string),
		Address3: firstResult["address3"].(string),
	}

	return response, nil
}

func callApi(req *http.Request) (*http.Response, error) {
	const (
		maxRetries     = 5
		initialBackoff = 200 * time.Millisecond
	)

	var resp *http.Response
	var err error
	backoff := initialBackoff
	client := &http.Client{}

	// リクエストヘッダーの追加
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "GoLambdaFunction/1.0")

	for i := 0; i <= maxRetries; i++ {
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		if resp != nil {
			resp.Body.Close()
		}

		if i == maxRetries {
			break
		}

		time.Sleep(backoff)
		backoff *= 2
	}

	return nil, fmt.Errorf("failed to fetch postal code info: %w", err)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Handler function executed.")

	// LambdaのリクエストIDを取得
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		log.Println("Failed to retrieve Lambda context")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Unable to retrieve Lambda context"}`,
		}, nil
	}

	log.Printf("Lambda Request ID: %s", lc.AwsRequestID)
	log.Printf("環境識別子： %s", config.ENV)

	// リクエストボディのパース
	var requestBody model.RequestBody
	if err := json.Unmarshal([]byte(req.Body), &requestBody); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid request body"}`,
		}, nil
	}

	// バリデーション
	if err := validate.Struct(requestBody); err != nil {
		log.Printf("Validation failed: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Validation error"}`,
		}, nil
	}

	// 郵便番号検索APIの呼び出し
	postalInfo, err := fetchPostalCodeInfo(requestBody.PostalCode)
	if err != nil {
		log.Printf("Error fetching postal info: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		}, nil
	}

	// レスポンス作成
	responseBody, _ := json.Marshal(postalInfo)
	return events.APIGatewayProxyResponse{
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
