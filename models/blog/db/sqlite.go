package db

import (
	"hypermedlab/backend-myblog/models/blog"
	"hypermedlab/backend-myblog/pkgs/sort"
	"hypermedlab/backend-myblog/pkgs/uuid"

	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type Sqlite3 struct {
	db *sqlx.DB
}

func NewSqlite3(conn *sqlx.DB) *Sqlite3 {
	p := &Sqlite3{
		db: conn,
	}

	if p.createNewTable() != nil {
		panic("create table blog, comment failed")
	}

	return p
}

func (s *Sqlite3) createNewTable() error {
	_, err := s.db.Exec(createBlogTable)
	return err
}

func (s *Sqlite3) CreateBlog(b *blog.Blog) (*blog.Blog, error) {
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

func (s *Sqlite3) CreateComment(c *blog.Comment) (*blog.Comment, error) {
	_, err := s.db.Exec("INSERT INTO comment (id, content, userid, username, blogid) VALUES ($1, $2, $3, $4, $5)", c.ID, c.Content, c.UserID, c.UserName, c.BlogID)
	if err != nil {
		return nil, err
	}

	return c, nil
}

type tmpBlog struct {
	Title           string         `db:"title"`
	Content         string         `db:"content"`
	ID              string         `db:"id"`
	UserID          string         `db:"userid"`
	UserName        string         `db:"username"`
	CreateAt        time.Time      `db:"create_at"`
	UpdateAt        time.Time      `db:"update_at"`
	CommentID       string         `db:"commentid"`
	CommentUserID   string         `db:"commentuserid"`
	CommentUserName string         `db:"commentusername"`
	CommentContent  string         `db:"commentcontent"`
	CommentBlogID   string         `db:"commentblogid"`
	CommentCreateAt time.Time      `db:"commentcreate_at"`
	CommentUpdateAt time.Time      `db:"commentupdate_at"`
	TagID           sql.NullString `db:"tagid"`
	TagName         sql.NullString `db:"tagname"`
}

func (s *Sqlite3) FindBlogsByContent(content string) ([]*blog.Blog, error) {
	blogs := make([]*blog.Blog, 0)
	err := s.db.Select(&blogs, searchBlog, "%"+content+"%", "%"+content+"%")
	if err != nil {
		return nil, err
	}

	return blogs, nil
}

func (s *Sqlite3) FindBlogByID(id string) (*blog.Blog, error) {
	tmpblgs := make([]*tmpBlog, 0)
	err := s.db.Select(&tmpblgs, findBlogByID, id)
	if err != nil {
		return nil, err
	}

	if len(tmpblgs) == 0 {
		return nil, errors.New("err, blog not exist")
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

	commentsMap := make(map[string]bool)
	tagsMap := make(map[string]bool)
	comments := make([]*blog.Comment, 0)
	tags := make([]*blog.Tag, 0)
	for _, tmpblg := range tmpblgs {
		if _, ok := commentsMap[tmpblg.CommentID]; !ok {
			comment := &blog.Comment{
				ID:       tmpblg.CommentID,
				Content:  tmpblg.CommentContent,
				UserID:   tmpblg.CommentUserID,
				UserName: tmpblg.CommentUserName,
				CreateAt: tmpblg.CommentCreateAt,
				UpdateAt: tmpblg.CommentUpdateAt,
			}
			comments = append(comments, comment)
			commentsMap[tmpblg.CommentID] = true
		}

		if _, ok := tagsMap[tmpblg.TagID.String]; !ok {
			tag := &blog.Tag{
				ID:   tmpblg.TagID.String,
				Name: tmpblg.TagName.String,
			}
			tags = append(tags, tag)
			tagsMap[tmpblg.TagID.String] = true
		}
	}
	if len(comments) > 1 {
		blg.Comment = comments[1:]
	}
	blg.Tags = tags
	return blg, nil
}

func (s *Sqlite3) FindBlogsByTitle(title string) ([]*blog.Blog, error) {
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

	sort.Sort(blgs)

	return blgs, nil
}

func (s *Sqlite3) FindBlogsByUserID(userid string) ([]*blog.Blog, error) {
	tmpblogs := make([]*tmpBlog, 0)
	err := s.db.Select(&tmpblogs, findBlogsByUserID, userid)
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

	sort.Sort(blgs)

	return blgs, nil
}

func (s *Sqlite3) FindBlogsByUserName(username string) ([]*blog.Blog, error) {
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

	sort.Sort(blgs)

	return blgs, nil
}

func (s *Sqlite3) FindBlogsByTagID(id string) ([]*blog.Blog, error) {
	tmpblogs := make([]*tmpBlog, 0)
	err := s.db.Select(&tmpblogs, findBlogsByTagID, id)
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

	sort.Sort(blgs)

	return blgs, nil
}

func (s *Sqlite3) FindCommentsByBlogID(blogid string) ([]*blog.Comment, error) {
	comments := make([]*blog.Comment, 0)
	err := s.db.Select(&comments, "SELECT id, content. blogid, userid, username, create_at, update_at FROM comment WHERE blogid=$1", blogid)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *Sqlite3) FindCommentsByUserID(userid string) ([]*blog.Comment, error) {
	comments := make([]*blog.Comment, 0)
	err := s.db.Select(&comments, "SELECT id, content. blogid, userid, username, create_at, update_at FROM comment WHERE userid=$1", userid)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *Sqlite3) FindCommentByID(id string) (*blog.Comment, error) {
	var c *blog.Comment
	err := s.db.Get(c, "SELECT id, content. blogid, userid, username, create_at, update_at FROM comment WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *Sqlite3) UpdateBlog(b *blog.Blog) (*blog.Blog, error) {
	_, err := s.db.Exec("UPDATE blog SET content=$1, title=$2, update_at=$3 WHERE id=$4", b.Content, b.Title, time.Now(), b.ID)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Sqlite3) UpdateComment(c *blog.Comment) (*blog.Comment, error) {
	_, err := s.db.Exec("UPDATE comment SET content=$1, update_at=$2 WHERE id=$3", c.Content, time.Now(), c.ID)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *Sqlite3) DeleteBlogByID(id string) error {
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

func (s *Sqlite3) DeleteCommentByID(id string) error {
	_, err := s.db.Exec("DELETE FROM comment WHERE id=$1", id)
	return err
}

func (s *Sqlite3) DeleteBlogByUserID(id string) error {
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

func (s *Sqlite3) DeleteCommentByUserID(id string) error {
	_, err := s.db.Exec("DELETE FROM comment WHERE userid=$1", id)
	return err
}

func (s *Sqlite3) CreateNewTagForBlog(tag *blog.Tag) (*blog.Tag, error) {
	_, err := s.db.Exec("INSERT INTO tag (id, name) VALUES ($1, $2)", tag.ID, tag.Name)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *Sqlite3) AllocateTagForBlog(tagid, blogid string) error {
	_, err := s.db.Exec("INSERT INTO tags (tag_id, blog_id) VALUES ($1, $2)", tagid, blogid)
	if err != nil {
		return err
	}

	return nil
}

func (s *Sqlite3) FindTagByName(name string) (*blog.Tag, error) {
	tag := new(blog.Tag)
	err := s.db.Get(tag, "SELECT * FROM tag WHERE name=$1", name)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (s *Sqlite3) DeleteTag(tagid, blogid string) error {
	_, err := s.db.Exec("DELETE * FROM tags WHERE tag_id=$1 AND blog_id=$2", tagid, blogid)
	if err != nil {
		return err
	}

	return nil
}
