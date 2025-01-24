package main

type (
	GetTodoResponse struct {
		Message string `json:"message"`
		Data    []Todo `json:"data"`
	}
)
