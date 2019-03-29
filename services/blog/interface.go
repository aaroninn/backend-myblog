package blog

import (
	"hypermedlab/backend-myblog/models/blog"
	"hypermedlab/backend-myblog/pkgs/forms"
)

type Service interface {
	CreateBlog(*forms.CreateBlog) (*blog.Blog, error)
	CreateComment(*forms.CreateComment) (*blog.Blog, error)

	FindBlogByID(string) (*blog.Blog, error)
	FindBlogsByUserID(string) ([]*blog.Blog, error)
	FindBlogsByTitle(string) ([]*blog.Blog, error)
	FindBlogsByUserName(string) ([]*blog.Blog, error)
	FindCommentByID(string) (*blog.Comment, error)
	FindCommentsByUserID(string) ([]*blog.Comment, error)

	UpdateBlog(*forms.UpdateBlog) (*blog.Blog, error)
	UpdateComment(*forms.UpdateComment) (*blog.Blog, error)

	DeleteBlogByID(string) error
	DeleteBlogByUserID(string) error
	DeleteCommentByID(string) error
	DeleteCommentByUserID(string) error
}
