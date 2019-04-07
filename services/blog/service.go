package blog

import (
	mBlog "hypermedlab/backend-myblog/models/blog"
	blogDB "hypermedlab/backend-myblog/models/blog/db"
	"hypermedlab/backend-myblog/pkgs/forms"
	"hypermedlab/backend-myblog/pkgs/uuid"

	"time"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	db *blogDB.Sqlite3
}

func NewBlogService(conn *sqlx.DB) *Service {
	return &Service{
		db: blogDB.NewSqlite3(conn),
	}
}

func (s *Service) CreateBlog(form *forms.CreateBlog) (*mBlog.Blog, error) {
	b := &mBlog.Blog{
		ID:       uuid.NewV1(),
		Title:    form.Title,
		Content:  form.Content,
		UserID:   form.UserID,
		UserName: form.UserName,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	_, err := s.db.CreateBlog(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) CreateComment(form *forms.CreateComment) (*mBlog.Blog, error) {
	b, err := s.FindBlogByID(form.BlogID)
	if err != nil {
		return nil, errBlogNotExist
	}

	c := &mBlog.Comment{
		ID:       uuid.NewV1(),
		Content:  form.Content,
		BlogID:   form.BlogID,
		UserID:   form.UserID,
		UserName: form.UserName,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	_, err = s.db.CreateComment(c)
	if err != nil {
		return nil, err
	}

	b.Comment = append(b.Comment, c)

	return b, nil
}

func (s *Service) FindBlogByID(id string) (*mBlog.Blog, error) {
	return s.db.FindBlogByID(id)
}

func (s *Service) FindBlogsByUserID(userid string) ([]*mBlog.Blog, error) {
	return s.db.FindBlogsByUserID(userid)
}

func (s *Service) FindBlogsByTitle(title string) ([]*mBlog.Blog, error) {
	return s.db.FindBlogsByTitle(title)
}

func (s *Service) FindBlogsByUserName(username string) ([]*mBlog.Blog, error) {
	return s.db.FindBlogsByUserName(username)
}

func (s *Service) FindCommentByID(id string) (*mBlog.Comment, error) {
	return s.db.FindCommentByID(id)
}

func (s *Service) FindCommentsByUserID(id string) ([]*mBlog.Comment, error) {
	return s.db.FindCommentsByUserID(id)
}

func (s *Service) UpdateBlog(form *forms.UpdateBlog) (*mBlog.Blog, error) {
	b, err := s.FindBlogByID(form.BlogID)
	if err != nil {
		return nil, errBlogNotExist
	}

	b.Content = form.Content
	b.Title = form.Title
	b.UpdateAt = time.Now()

	_, err = s.db.UpdateBlog(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) UpdateComment(form *forms.UpdateComment) (*mBlog.Blog, error) {
	comment := &mBlog.Comment{
		ID:       form.CommentID,
		Content:  form.Content,
		UpdateAt: time.Now(),
	}

	_, err := s.db.UpdateComment(comment)
	if err != nil {
		return nil, err
	}

	b, err := s.FindBlogByID(form.BlogID)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) DeleteBlogByID(id string) error {
	return s.db.DeleteBlogByID(id)
}

func (s *Service) DeleteBlogByUserID(id string) error {
	return s.db.DeleteBlogByUserID(id)
}

func (s *Service) DeleteCommentByID(id string) error {
	return s.db.DeleteCommentByID(id)
}

func (s *Service) DeleteCommentByUserID(id string) error {
	return s.db.DeleteCommentByUserID(id)
}
