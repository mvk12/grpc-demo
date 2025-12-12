package models

type Student struct {
	ID    int32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Test struct {
	ID          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}