package blog

import (
	"time"
)

type Blog struct {
	Title    string     `db:"title" json:"title"`
	ID       string     `db:"id" json:"id"`
	UserName string     `db:"username" json:"userName"`
	UserID   string     `db:"userid" json:"userID"`
	Content  string     `db:"content" json:"content"`
	Comment  []*Comment `json:"comment,omitempty"`
	CreateAt time.Time  `db:"create_at" json:"createAt"`
	UpdateAt time.Time  `db:"update_at" json:"updateAt"`
}

type Comment struct {
	ID       string    `db:"id" json:"id,omitempty"`
	BlogID   string    `db:"blogid" json:"blogID,omitempty"`
	UserName string    `db:"username" json:"userName,omitempty"`
	UserID   string    `db:"userid" json:"userID,omitempty"`
	Content  string    `db:"content" json:"content,omitempty"`
	CreateAt time.Time `db:"create_at" json:"createAt,omitempty"`
	UpdateAt time.Time `db:"update_at" json:"updateAt,omitempty"`
}
