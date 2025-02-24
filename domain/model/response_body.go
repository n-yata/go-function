package model

import (
	"encoding/json"
	"log"
)

// レスポンスボディ
type ResponseBody struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	Address3 string `json:"address3"`
}

func (rb *ResponseBody) ToJSON() string {
	jsonData, err := json.Marshal(rb)
	if err != nil {
		log.Printf("Error marshalling response body: %v", err)
		return "{}"
	}
	return string(jsonData)
}
