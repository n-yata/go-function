package model

// リクエストボディ
type RequestBody struct {
	PostalCode string `json:"postalCode" validate:"required"`
}
