package dto

type LessonRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Position int    `json:"position" binding:"required,min=1"`
}
