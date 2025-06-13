package entity

type ValidationErrorField struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value any    `json:"value"`
}
