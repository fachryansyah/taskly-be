package task

type CreateTaskRequest struct {
	Title string `json:"title" validate:"required"`
	Desc  string `json:"desc"`
	Label string `json:"label" validate:"required"`
}

type EditTaskRequest struct {
	Title string `json:"title" validate:"required"`
	Desc  string `json:"desc"`
	Label string `json:"label" validate:"required"`
}
