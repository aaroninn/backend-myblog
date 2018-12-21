package user

import (
	"time"
)

//User struct
type User struct {
	ID             string `db:"id" json:"id,omitempty"`
	Account        string `db:"account" json:"account,omitempty"`
	Email          string `db:"email" json:"email,omitempty"`
	Name           string `db:"name" json:"name,omitempty"`
	Description    string `db:"description" json:"description,omitempty"`
	HashedPassword string `db:"hashedpassword" json:"-"`
	Token          string `json:"token"`
	Baned          bool   `db:"baned" json:"baned,omitempty"`
	Root           bool   `db:"root" json:"root,omitempty"`
	Friends        []*User
	CreateAt       time.Time `db:"create_at" json:"createAt,omitempty"`
	LastLogin      time.Time `db:"last_login" json:"lastLogin,omitempty"`
}
