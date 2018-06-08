package controllers

type Response struct {
	Data    map[string]interface{} `json:"data"`
	Message string                 `json:"message"`
}
