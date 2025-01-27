package todo

func MapTodoListModelToDto(todoModels []TodoModel) []TodoDto {
	todoList := []TodoDto{}
	for _, td := range todoModels {
		todoList = append(todoList, TodoDto{
			Title:     td.Title,
			ID:        td.ID,
			CreatedAt: td.CreatedAt,
			Completed: td.Completed,
		})
	}
	return todoList
}
