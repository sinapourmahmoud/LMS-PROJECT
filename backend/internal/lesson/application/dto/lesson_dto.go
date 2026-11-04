package dto

type LessonRequest struct {
	Title    string `json:"title" validate:"required"`
	Content  string `json:"content" validate:"required"`
	Position int    `json:"position" validate:"required,min=1"`
}
