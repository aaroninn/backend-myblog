package user

import "hypermedlab/myblog/pkgs/forms"

//DB User DB interface
type DB interface {
	CreateUser(*User) (*User, error)
	InsertFriend(forms.InsertFriend) error

	FindUserByID(string) (*User, error)
	FindUserByAccount(string) (*User, error)

	UpdateUserPassword(*User) error
	UpdateUser(*User) error //update description and name
	UpdateUserStatus(string, bool) error

	DeleteUserByID(string) error
}
