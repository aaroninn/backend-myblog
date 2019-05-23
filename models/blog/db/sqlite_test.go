package db

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
)

const blogIDPrefix = "blogid_"
const blogTitlePrefix = "blogtitle_"
const blogContentPrefix = "blogcontent_"
const userIDPrefix = "userid_"
const userNamePrefix = "username_"
const commentIDPrefix = "commentid_"
const commentContentPrefix = "commentcontent_"

func Test_sqlite(t *testing.T) {
	conn, err := sqlx.Open("postgres", "user=postgres password=postgres")
	conn.Exec("create database test")
	conn.Close()
	db, err := sqlx.Open("postgres", "user=postgres password=postgres dbname=test")
	if err != nil {
		fmt.Printf("%v", err)
	}

	if err != nil {
		t.Errorf("connect to sqlite3 failed: %s", err.Error())
	}

	defer db.Close()

	Convey("test INSERT", t, func() {
		engine := NewBlogPostgre(db)
		for idx := 0; idx < 3; idx++ {
			c := &blog.Comment{
				ID:       commentIDPrefix + strconv.Itoa(idx),
				BlogID:   blogIDPrefix + strconv.Itoa(idx),
				Content:  blogContentPrefix + strconv.Itoa(idx),
				UserID:   userIDPrefix + strconv.Itoa(idx+3),
				UserName: userNamePrefix + strconv.Itoa(idx+3),
			}
			b := &blog.Blog{
				ID:       blogIDPrefix + strconv.Itoa(idx),
				Title:    blogTitlePrefix + strconv.Itoa(idx),
				Content:  blogTitlePrefix + strconv.Itoa(idx),
				UserID:   userIDPrefix + strconv.Itoa(idx),
				UserName: userNamePrefix + strconv.Itoa(idx),
			}
			_, err := engine.CreateBlog(b)
			So(err, ShouldBeNil)
			_, err = engine.CreateComment(c)
			So(err, ShouldBeNil)
		}
		Convey("test FindByBlogID", func() {
			for idx := 0; idx < 3; idx++ {
				_, err := engine.FindBlogByID(blogIDPrefix + strconv.Itoa(idx))
				So(err, ShouldBeNil)
			}
		})

	})
}
