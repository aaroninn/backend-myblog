package blog

type DB interface {
	CreateBlog(*Blog) (*Blog, error)
	CreateComment(*Comment) (*Comment, error)

	FindBlogsByTitle(string) ([]*Blog, error)
	FindBlogByID(string) (*Blog, error)
	FindBlogsByUserID(string) ([]*Blog, error)
	FindBlogsByUserName(string) ([]*Blog, error)
	FindCommentsByBlogID(string) ([]*Comment, error)
	FindCommentByID(string) (*Comment, error)
	FindCommentsByUserID(string) ([]*Comment, error)

	UpdateBlog(*Blog) (*Blog, error)
	UpdateComment(*Comment) (*Comment, error)

	DeleteBlogByID(string) error
	DeleteBlogByUserID(string) error
	DeleteCommentByID(string) error
	DeleteCommentByUserID(string) error
}
