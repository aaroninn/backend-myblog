package blog

import (
	"github.com/jmoiron/sqlx"
	mBlog "hypermedlab/myblog/models/blog"
	blogDB "hypermedlab/myblog/models/blog/db"
	"hypermedlab/myblog/pkgs/forms"
	"hypermedlab/myblog/pkgs/uuid"
	"time"
)

type service struct {
	db mBlog.DB
}

func NewBlogService(conn *sqlx.DB) Service {
	return &service{
		db: blogDB.NewBlogPostgre(conn),
	}
}

func (s *service) CreateBlog(form *forms.CreateBlog) (*mBlog.Blog, error) {
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

func (s *service) CreateComment(form *forms.CreateComment) (*mBlog.Blog, error) {
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

func (s *service) FindBlogByID(id string) (*mBlog.Blog, error) {
	b, err := s.db.FindBlogByID(id)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *service) FindBlogsByUserID(userid string) ([]*mBlog.Blog, error) {
	blogs, err := s.db.FindBlogsByUserID(userid)
	if err != nil {
		return nil, err
	}

	return blogs, nil
}

func (s *service) FindBlogsByTitle(title string) ([]*mBlog.Blog, error) {
	blogs, err := s.db.FindBlogsByTitle(title)
	if err != nil {
		return nil, err
	}

	return blogs, nil
}

func (s *service) FindCommentByID(id string) (*mBlog.Comment, error) {
	comment, err := s.db.FindCommentByID(id)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *service) FindCommentsByUserID(id string) ([]*mBlog.Comment, error) {
	comments, err := s.db.FindCommentsByUserID(id)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *service) UpdateBlog(form *forms.UpdateBlog) (*mBlog.Blog, error) {
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

func (s *service) UpdateComment(form *forms.UpdateComment) (*mBlog.Blog, error) {
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

func (s *service) DeleteBlogByID(id string) error {
	err := s.db.DeleteBlogByID(id)
	return err
}

func (s *service) DeleteBlogByUserID(id string) error {
	err := s.db.DeleteBlogByUserID(id)
	return err
}

func (s *service) DeleteCommentByID(id string) error {
	err := s.db.DeleteCommentByID(id)
	return err
}

func (s *service) DeleteCommentByUserID(id string) error {
	err := s.db.DeleteCommentByUserID(id)
	return err
}
