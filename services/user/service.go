package user

import (
	mUser "hypermedlab/backend-myblog/models/user"
	userDB "hypermedlab/backend-myblog/models/user/db"
	"hypermedlab/backend-myblog/pkgs/forms"
	"hypermedlab/backend-myblog/pkgs/jwt"
	"hypermedlab/backend-myblog/pkgs/password"
	"hypermedlab/backend-myblog/pkgs/session"
	"hypermedlab/backend-myblog/pkgs/uuid"

	"errors"
)

type Service struct {
	db       *userDB.Sqlite3
	sessions *session.SessionsStorageInMemory
}

//NewService user
func NewService(sql *userDB.Sqlite3, sessions *session.SessionsStorageInMemory) *Service {
	return &Service{
		sql,
		sessions,
	}
}

func (s *Service) RegisterUser(form *forms.CreateUser) (*mUser.User, error) {
	hashedpwd, err := password.HashedPassword(form.Password)
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

	return s.db.CreateUser(usr)
}

func (s *Service) Login(form *forms.LoginForm, secret string) (*mUser.User, error) {
	usr, err := s.db.FindUserByAccount(form.Account)
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
	sess, ok := s.sessions.Get(usr.ID)
	if !ok {
		newSess := session.NewSession(usr.ID)
		newSess.SetData(usr.Token)
		s.sessions.Add(newSess)
	} else {
		sess.SetData(usr.Token)
		s.sessions.RefeshSession(usr.ID)
	}

	return usr, nil
}

func (s *Service) LogOut(id string) {
	s.sessions.Delete(id)
}

func (s *Service) UpdatePassword(form *forms.UpdatePassword) error {
	user, err := s.db.FindUserByID(form.UserID)
	if err != nil {
		return errors.New("user not exist")
	}

	err = password.ComparePassword(form.PrePassword, user.HashedPassword)
	if err != nil {
		return errors.New("password not correct")
	}

	hashedPw, err := password.HashedPassword(form.Password)
	if err != nil {
		return err
	}

	usr := &mUser.User{
		ID:             form.UserID,
		HashedPassword: hashedPw,
	}

	err = s.db.UpdateUserPassword(usr)
	if err != nil {
		return errors.New("update pw failed")
	}

	s.sessions.Delete(usr.ID)

	return nil
}

func (s *Service) FindAllUsers() ([]*mUser.User, error) {
	return s.FindAllUsers()
}

func (s *Service) UpdateUserStatus(userid string, status bool) error {
	return s.UpdateUserStatus(userid, status)
}

func (s *Service) GetSession(userid string) (*session.Session, bool) {
	return s.sessions.Get(userid)
}
