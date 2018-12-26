package db

import (
	"github.com/jmoiron/sqlx"
	"hypermedlab/myblog/models/blog"
	"hypermedlab/myblog/pkgs/uuid"
	"time"
)

type sqlite struct {
	db *sqlx.DB
}

const createBlogTable = `
CREATE TABLE IF NOT EXISTS blog(
	id CHAR(40) NOT NULL PRIMARY KEY,
	title TEXT,
	content TEXT,
	userid CHAR(40) NOT NULL,
	username TEXT NOT NULL,
	create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS comment(
	id CHAR(40) NOT NULL PRIMARY KEY,
	content TEXT NOT NULL,
	userid CHAR(40) NOT NULL,
	username TEXT NOT NULL,
	blogid CHAR(40) NOT NULL,
	create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
`
const findBlogByID = `
SELECT 
blog.id AS id,
blog.title AS title,
blog.content AS content,
blog.userid AS userid,
blog.username AS username,
blog.create_at AS create_at,
blog.update_at AS update_at,
comment.id AS commentid,
comment.content AS commentcontent,
comment.userid AS commentuserid,
comment.username AS commentusername,
comment.blogid AS commentblogid,
comment.create_at AS commentcreate_at,
comment.update_at As commentupdate_at
FROM comment LEFT JOIN blog
ON blog.id = comment.blogid
WHERE blog.id = $1
`

const findBlogByUserID = `
SELECT 
blog.id AS id,
blog.title AS title,
blog.content AS content,
blog.userid AS userid,
blog.username AS username,
blog.create_at AS create_at,
blog.update_at AS update_at,
comment.id AS commentid,
comment.content AS commentcontent,
comment.userid AS commentuserid,
comment.username AS commentusername,
comment.blogid AS commentblogid,
comment.create_at AS commentcreate_at,
comment.update_at As commentupdate_at
FROM comment LEFT JOIN blog
ON blog.id = comment.blogid
WHERE blog.userid = $1
ORDER BY create_at DESC
`
const findBlogsByTitle = `
SELECT 
blog.id AS id,
blog.title AS title,
blog.content AS content,
blog.userid AS userid,
blog.username AS username,
blog.create_at AS create_at,
blog.update_at AS update_at,
comment.id AS commentid,
comment.content AS commentcontent,
comment.userid AS commentuserid,
comment.username AS commentusername,
comment.blogid AS commentblogid,
comment.create_at AS commentcreate_at,
comment.update_at As commentupdate_at
FROM comment LEFT JOIN blog
ON blog.id = comment.blogid
WHERE blog.title LIKE $1
`

const findBlogsByUserName = `
SELECT 
blog.id AS id,
blog.title AS title,
blog.content AS content,
blog.userid AS userid,
blog.username AS username,
blog.create_at AS create_at,
blog.update_at AS update_at,
comment.id AS commentid,
comment.content AS commentcontent,
comment.userid AS commentuserid,
comment.username AS commentusername,
comment.blogid AS commentblogid,
comment.create_at AS commentcreate_at,
comment.update_at As commentupdate_at
FROM comment LEFT JOIN blog
ON blog.id = comment.blogid
WHERE blog.username = $1
ORDER BY create_at DESC
`

func NewBlogPostgre(conn *sqlx.DB) blog.DB {
	p := &sqlite{
		db: conn,
	}

	if p.createNewTable() != nil {
		panic("create table blog, comment failed")
	}

	return p
}

func (s *sqlite) createNewTable() error {
	_, err := s.db.Exec(createBlogTable)
	return err
}

func (s *sqlite) CreateBlog(b *blog.Blog) (*blog.Blog, error) {
	_, err := s.db.Exec("INSERT INTO blog (id, title, content, userid, username) VALUES ($1, $2, $3, $4, $5)", b.ID, b.Title, b.Content, b.UserID, b.UserName)
	if err != nil {
		return nil, err
	}

	comment := &blog.Comment{
		ID:       uuid.NewV1(),
		BlogID:   b.ID,
		UserID:   "default",
		UserName: "default",
		Content:  "default",
	}
	_, err = s.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *sqlite) CreateComment(c *blog.Comment) (*blog.Comment, error) {
	_, err := s.db.Exec("INSERT INTO comment (id, content, userid, username, blogid) VALUES ($1, $2, $3, $4, $5)", c.ID, c.Content, c.UserID, c.UserName, c.BlogID)
	if err != nil {
		return nil, err
	}

	return c, nil
}

type tmpBlog struct {
	Title           string    `db:"title"`
	Content         string    `db:"content"`
	ID              string    `db:"id"`
	UserID          string    `db:"userid"`
	UserName        string    `db:"username"`
	CreateAt        time.Time `db:"create_at"`
	UpdateAt        time.Time `db:"update_at"`
	CommentID       string    `db:"commentid"`
	CommentUserID   string    `db:"commentuserid"`
	CommentUserName string    `db:"commentusername"`
	CommentContent  string    `db:"commentcontent"`
	CommentBlogID   string    `db:"commentblogid"`
	CommentCreateAt time.Time `db:"commentcreate_at"`
	CommentUpdateAt time.Time `db:"commentupdate_at"`
}

func (s *sqlite) FindBlogByID(id string) (*blog.Blog, error) {
	tmpblgs := make([]*tmpBlog, 0)
	err := s.db.Select(&tmpblgs, findBlogByID, id)
	if err != nil {
		return nil, err
	}

	blg := &blog.Blog{
		Title:    tmpblgs[0].Title,
		Content:  tmpblgs[0].Content,
		ID:       tmpblgs[0].ID,
		UserID:   tmpblgs[0].UserID,
		UserName: tmpblgs[0].UserName,
		CreateAt: tmpblgs[0].CreateAt,
		UpdateAt: tmpblgs[0].UpdateAt,
	}

	comments := make([]*blog.Comment, 0)
	for _, tmpblg := range tmpblgs {
		comment := &blog.Comment{
			ID:       tmpblg.CommentID,
			Content:  tmpblg.CommentContent,
			UserID:   tmpblg.UserID,
			UserName: tmpblg.UserName,
			CreateAt: tmpblg.CreateAt,
			UpdateAt: tmpblg.UpdateAt,
		}
		comments = append(comments, comment)
	}
	if len(comments) > 1 {
		blg.Comment = comments[1:]
	}
	return blg, nil
}

func (s *sqlite) FindBlogsByTitle(title string) ([]*blog.Blog, error) {
	tmpblogs := make([]*tmpBlog, 0)
	err := s.db.Select(&tmpblogs, findBlogsByTitle, "%"+title+"%")
	if err != nil {
		return nil, err
	}

	blogs := make(map[string]*blog.Blog, 0)
	comments := make(map[string][]*blog.Comment)
	for _, tmpblog := range tmpblogs {
		blogs[tmpblog.ID] = &blog.Blog{
			UserID:   tmpblog.UserID,
			UserName: tmpblog.UserName,
			Title:    tmpblog.Title,
			Content:  tmpblog.Content,
			ID:       tmpblog.ID,
			CreateAt: tmpblog.CreateAt,
			UpdateAt: tmpblog.UpdateAt,
		}
		cs := comments[tmpblog.ID]
		c := &blog.Comment{
			ID:       tmpblog.CommentID,
			Content:  tmpblog.CommentContent,
			UserID:   tmpblog.CommentUserID,
			UserName: tmpblog.CommentUserName,
			BlogID:   tmpblog.ID,
			CreateAt: tmpblog.CommentCreateAt,
			UpdateAt: tmpblog.CreateAt,
		}
		cs = append(cs, c)
		comments[tmpblog.ID] = cs
	}

	blgs := make([]*blog.Blog, 0)
	for _, blg := range blogs {
		if len(comments[blg.ID]) > 1 {
			blg.Comment = comments[blg.ID][1:]
		}
		blgs = append(blgs, blg)
	}

	return blgs, nil
}

func (s *sqlite) FindBlogsByUserID(userid string) ([]*blog.Blog, error) {
	tmpblogs := make([]*tmpBlog, 0)
	err := s.db.Select(&tmpblogs, findBlogByUserID, userid)
	if err != nil {
		return nil, err
	}

	blogs := make(map[string]*blog.Blog, 0)
	comments := make(map[string][]*blog.Comment, 0)
	for _, tmpblog := range tmpblogs {
		blogs[tmpblog.ID] = &blog.Blog{
			UserID:   tmpblog.UserID,
			UserName: tmpblog.UserName,
			Title:    tmpblog.Title,
			Content:  tmpblog.Content,
			ID:       tmpblog.ID,
			CreateAt: tmpblog.CreateAt,
			UpdateAt: tmpblog.UpdateAt,
		}
		cs := comments[tmpblog.ID]
		c := &blog.Comment{
			ID:       tmpblog.CommentID,
			Content:  tmpblog.CommentContent,
			UserID:   tmpblog.CommentUserID,
			UserName: tmpblog.CommentUserName,
			BlogID:   tmpblog.ID,
			CreateAt: tmpblog.CommentCreateAt,
			UpdateAt: tmpblog.CreateAt,
		}
		cs = append(cs, c)
		comments[tmpblog.ID] = cs
	}

	blgs := make([]*blog.Blog, 0)
	for _, blg := range blogs {
		if len(comments[blg.ID]) > 1 {
			blg.Comment = comments[blg.ID][1:]
		}
		blgs = append(blgs, blg)
	}

	return blgs, nil
}

func (s *sqlite) FindBlogsByUserName(username string) ([]*blog.Blog, error) {
	tmpblogs := make([]*tmpBlog, 0)
	err := s.db.Select(&tmpblogs, findBlogsByUserName, username)
	if err != nil {
		return nil, err
	}

	blogs := make(map[string]*blog.Blog, 0)
	comments := make(map[string][]*blog.Comment, 0)
	for _, tmpblog := range tmpblogs {
		blogs[tmpblog.ID] = &blog.Blog{
			UserID:   tmpblog.UserID,
			UserName: tmpblog.UserName,
			Title:    tmpblog.Title,
			Content:  tmpblog.Content,
			ID:       tmpblog.ID,
			CreateAt: tmpblog.CreateAt,
			UpdateAt: tmpblog.UpdateAt,
		}
		cs := comments[tmpblog.ID]
		c := &blog.Comment{
			ID:       tmpblog.CommentID,
			Content:  tmpblog.CommentContent,
			UserID:   tmpblog.CommentUserID,
			UserName: tmpblog.CommentUserName,
			BlogID:   tmpblog.ID,
			CreateAt: tmpblog.CommentCreateAt,
			UpdateAt: tmpblog.CreateAt,
		}
		cs = append(cs, c)
		comments[tmpblog.ID] = cs
	}

	blgs := make([]*blog.Blog, 0)
	for _, blg := range blogs {
		if len(comments[blg.ID]) > 1 {
			blg.Comment = comments[blg.ID][1:]
		}
		blgs = append(blgs, blg)
	}

	return blgs, nil
}

func (s *sqlite) FindCommentsByBlogID(blogid string) ([]*blog.Comment, error) {
	comments := make([]*blog.Comment, 0)
	err := s.db.Select(&comments, "SELECT id, content. blogid, userid, username, create_at, update_at FROM comment WHERE blogid=$1", blogid)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *sqlite) FindCommentsByUserID(userid string) ([]*blog.Comment, error) {
	comments := make([]*blog.Comment, 0)
	err := s.db.Select(&comments, "SELECT id, content. blogid, userid, username, create_at, update_at FROM comment WHERE userid=$1", userid)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *sqlite) FindCommentByID(id string) (*blog.Comment, error) {
	var c *blog.Comment
	err := s.db.Get(c, "SELECT id, content. blogid, userid, username, create_at, update_at FROM comment WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *sqlite) UpdateBlog(b *blog.Blog) (*blog.Blog, error) {
	_, err := s.db.Exec("UPDATE blog SET content=$1, title=$2, update_at=$3 WHERE id=$4", b.Content, b.Title, time.Now(), b.ID)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *sqlite) UpdateComment(c *blog.Comment) (*blog.Comment, error) {
	_, err := s.db.Exec("UPDATE comment SET content=$1, update_at=$2 WHERE id=$3", c.Content, time.Now(), c.ID)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *sqlite) DeleteBlogByID(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM blog WHERE id=$1", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM comment WHERE blogid=$1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *sqlite) DeleteCommentByID(id string) error {
	_, err := s.db.Exec("DELETE FROM comment WHERE id=$1", id)
	return err
}

func (s *sqlite) DeleteBlogByUserID(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM blog WHERE userid=$1", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM comment WHERE blogid IN (SELECT id FROM blog WHERE id=$1)", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *sqlite) DeleteCommentByUserID(id string) error {
	_, err := s.db.Exec("DELETE FROM comment WHERE userid=$1", id)
	return err
}
