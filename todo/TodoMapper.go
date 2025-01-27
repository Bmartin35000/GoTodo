package todo

func MapTodoListModelToDto(models []TodoModel) []TodoDto {
	dtoList := []TodoDto{}
	for _, dto := range models {
		dtoList = append(dtoList, MapTodoModelToDto(dto))
	}
	return dtoList
}

func MapTodoModelToDto(model TodoModel) TodoDto {
	return TodoDto{model.ID, model.Title, model.Completed, model.CreatedAt}
}
