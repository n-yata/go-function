package model

import (
	"encoding/json"
	"log"
)

// レスポンスボディ
type ResponseBody struct {
	Message string `json:"message"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
}

func (rb *ResponseBody) ToJSON() string {
	jsonData, err := json.Marshal(rb)
	if err != nil {
		log.Printf("Error marshalling response body: %v", err)
		return "{}" // Return an empty JSON object in case of error
	}
	return string(jsonData)
}
