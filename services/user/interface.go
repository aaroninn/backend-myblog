package user

import (
	mUser "hypermedlab/backend-myblog/models/user"
	"hypermedlab/backend-myblog/pkgs/forms"
)

//Service user
type Service interface {
	RegisterUser(form *forms.CreateUser) (*mUser.User, error)
	Login(form *forms.LoginForm, secret string) (*mUser.User, error)
	UpdatePassword(form *forms.UpdatePassword) error
	FindAllUsers() ([]*mUser.User, error)
	UpdateUserStatus(string, bool) error
}
