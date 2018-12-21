package forms

type CreateBlog struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	UserID   string
	UserName string
}

type CreateComment struct {
	BlogID   string `json:'blogid"`
	Content  string `json:"content"`
	UserID   string
	UserName string
}

type UpdateBlog struct {
	BlogID  string
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  string
}

type UpdateComment struct {
	CommentID string
	BlogID    string
	Content   string `json:"content"`
	UserID    string
}
