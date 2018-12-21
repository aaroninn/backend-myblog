package user

import (
	"errors"
	"github.com/jmoiron/sqlx"
	mUser "hypermedlab/myblog/models/user"
	userDB "hypermedlab/myblog/models/user/db"
	"hypermedlab/myblog/pkgs/forms"
	"hypermedlab/myblog/pkgs/jwt"
	"hypermedlab/myblog/pkgs/password"
	"hypermedlab/myblog/pkgs/uuid"
	"log"
)

type user struct {
	db mUser.DB
}

//NewService user
func NewService(conn *sqlx.DB) Service {
	return &user{
		userDB.NewPostgre(conn),
	}
}

func (u *user) RegisterUser(form *forms.CreateUser) (*mUser.User, error) {
	hashedpwd, err := password.HashePassword(form.Password)
	if err != nil {
		return nil, err
	}
	usr := &mUser.User{
		ID:             uuid.NewV1(),
		Account:        form.Account,
		Email:          form.EMail,
		Name:           form.Name,
		HashedPassword: hashedpwd,
		Description:    form.Description,
	}

	return u.db.CreateUser(usr)
}

func (u *user) Login(form *forms.LoginForm, secret string) (*mUser.User, error) {
	usr, err := u.db.FindUserByAccount(form.Account)
	if err != nil {
		return nil, errEmptyAccount
	}

	if password.ComparePassword(form.Password, usr.HashedPassword) != nil {
		return nil, errWrongPassword
	}

	claims := jwt.NewCustomClaims(usr.ID, usr.Name, usr.Account, usr.Email)
	token, err := jwt.NewToken(claims, secret)
	if err != nil {
		return nil, err
	}

	usr.Token = token
	return usr, nil
}

func (u *user) UpdatePassword(form *forms.UpdatePassword) error {
	user, err := u.db.FindUserByID(form.UserID)
	if err != nil {
		log.Println(form.UserID)
		log.Println(err)
		return errors.New("user not exist")
	}

	err = password.ComparePassword(form.PrePassword, user.HashedPassword)
	if err != nil {
		return errors.New("password not correct")
	}

	hashedPw, err := password.HashePassword(form.Password)
	if err != nil {
		return err
	}

	usr := &mUser.User{
		ID:             form.UserID,
		HashedPassword: hashedPw,
	}

	err = u.db.UpdateUserPassword(usr)
	if err != nil {
		return errors.New("update pw failed")
	}

	return nil
}
