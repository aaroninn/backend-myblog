package db

import (
	"hypermedlab/backend-myblog/models/user"
	"hypermedlab/backend-myblog/pkgs/forms"
	"hypermedlab/backend-myblog/pkgs/password"
	"log"

	"github.com/jmoiron/sqlx"
)

const createTableUser = `
CREATE TABLE IF NOT EXISTS bloguser  (
	id CHAR(40) NOT NULL PRIMARY KEY,
	account TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL,
	name TEXT NOT NULL UNIQUE,
	description TEXT,
	hashedpassword TEXT NOT NULL,
	baned BOOL DEFAULT FALSE,
	root BOOL DEFAULT FALSE,
	create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS friend (
	id CHAR(40) NOT NULL,
	friend_id CHAR(40) NOT NULL
)
`

type Sqlite3 struct {
	db *sqlx.DB
}

//NewSqlite3 sql connection
func NewSqlite3(conn *sqlx.DB) *Sqlite3 {
	s := &Sqlite3{
		db: conn,
	}

	if s.createNewTable() != nil {
		panic("create table usr failed")
	}

	if s.createAdmin() != nil {
		panic("create admin failed")
	}

	return s
}
func (s *Sqlite3) createNewTable() error {
	_, err := s.db.Exec(createTableUser)
	return err
}

func (s *Sqlite3) createAdmin() error {
	u, _ := s.FindUserByID("000000")
	if u != nil {
		return nil
	}
	hasehdpwd, err := password.HashedPassword("123456")
	if err != nil {
		return err
	}
	_, err = s.db.Exec("INSERT INTO bloguser (id, account, name, email, hashedpassword, root, description) VALUES ($1, $2, $3, $4, $5, $6, $7)", "000000", "rootadmin", "admin", "448338094@qq.com", hasehdpwd, true, "")

	return err
}

func (s *Sqlite3) CreateUser(u *user.User) (*user.User, error) {
	_, err := s.db.Exec("INSERT INTO bloguser (id, account, name, email, description, hashedpassword) VALUES ($1, $2, $3, $4, $5, $6)", u.ID, u.Account, u.Name, u.Email, u.Description, u.HashedPassword)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Sqlite3) InsertFriend(form forms.InsertFriend) error {
	_, err := s.db.Exec("INSERT INTO friend (id, friend_id) VALUES ($1, $2)", form.ID, form.FriendID)
	return err
}

func (s *Sqlite3) FindUserByID(id string) (*user.User, error) {
	users := make([]*user.User, 0)
	err := s.db.Select(&users, "SELECT id, account, email, name, description, hashedpassword, baned, root, create_at, last_login FROM bloguser WHERE id IN (SELECT friend_id  FROM friend WHERE id=$1) OR id=$1", id)
	if err != nil {
		return nil, err
	}

	u := &user.User{}
	var i int
	for index, usr := range users {
		if usr.ID == id {
			u = usr
			i = index
			break
		}
	}

	if len(users) > 1 {
		users = append(users[:i], users[i+1:]...)
	}

	u.Friends = users

	return u, nil
}

func (s *Sqlite3) FindUserByAccount(acocunt string) (*user.User, error) {
	var usr user.User
	err := s.db.Get(&usr, "SELECT id, account, email, name, description, hashedpassword, baned, root, create_at, last_login FROM bloguser WHERE account=$1", acocunt)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &usr, nil
}

func (s *Sqlite3) UpdateUserPassword(u *user.User) error {
	_, err := s.db.Exec("UPDATE bloguser SET hashedpassword=$1 WHERE id=$2", u.HashedPassword, u.ID)
	return err
}

//update description and name
func (s *Sqlite3) UpdateUser(u *user.User) error {
	_, err := s.db.Exec("UPDATE bloguser SET description=$1, name=$2 WHERE id=$3", u.Description, u.Name, u.ID)
	return err
}

func (s *Sqlite3) UpdateUserStatus(id string, status bool) error {
	_, err := s.db.Exec("UPDATE bloguser SET baned=$1 WHERE id=$2", status, id)
	return err
}

func (s *Sqlite3) DeleteUserByID(id string) error {
	_, err := s.db.Exec("DELETE * FROM bloguser WHERE id=$1", id)
	return err
}

func (s *Sqlite3) FindAllUsers() ([]*user.User, error) {
	users := make([]*user.User, 0)
	err := s.db.Select(&users, "SELECT * FROM bloguser")
	if err != nil {
		return nil, err
	}

	return users, nil
}
