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

type Question struct {
	ID       int32  `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	TestID   int32  `json:"test_id"`
}