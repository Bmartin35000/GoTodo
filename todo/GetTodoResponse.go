package todo

type GetTodoResponse struct {
	Message string    `json:"message"`
	Data    []TodoDto `json:"data"`
}
