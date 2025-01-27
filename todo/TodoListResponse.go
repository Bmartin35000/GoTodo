package todo

type TodoListResponse struct {
	Message string    `json:"message"`
	Data    []TodoDto `json:"data"`
}
