package model

// リクエストボディ
type RequestBody struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"min=0"`
}
