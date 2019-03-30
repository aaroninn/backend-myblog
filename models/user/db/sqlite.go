package db

import (
	"github.com/jmoiron/sqlx"
	"hypermedlab/backend-myblog/models/user"
	"hypermedlab/backend-myblog/pkgs/forms"
	"hypermedlab/backend-myblog/pkgs/password"
	"log"
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

type postgre struct {
	db *sqlx.DB
}

//NewPostgre sql connection
func NewPostgre(conn *sqlx.DB) user.DB {
	p := &postgre{
		db: conn,
	}

	if p.createNewTable() != nil {
		panic("create table usr failed")
	}

	if p.createAdmin() != nil {
		panic("create admin failed")
	}

	return p
}
func (p *postgre) createNewTable() error {
	_, err := p.db.Exec(createTableUser)
	return err
}

func (p *postgre) createAdmin() error {
	u, _ := p.FindUserByID("000000")
	if u != nil {
		return nil
	}
	hasehdpwd, err := password.HashPassword("123456")
	if err != nil {
		return err
	}
	_, err = p.db.Exec("INSERT INTO bloguser (id, account, name, email, hashedpassword, root, description) VALUES ($1, $2, $3, $4, $5, $6, $7)", "000000", "rootadmin", "admin", "448338094@qq.com", hasehdpwd, true, "")

	return err
}

func (p *postgre) CreateUser(u *user.User) (*user.User, error) {
	_, err := p.db.Exec("INSERT INTO bloguser (id, account, name, email, description, hashedpassword) VALUES ($1, $2, $3, $4, $5, $6)", u.ID, u.Account, u.Name, u.Email, u.Description, u.HashedPassword)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (p *postgre) InsertFriend(form forms.InsertFriend) error {
	_, err := p.db.Exec("INSERT INTO friend (id, friend_id) VALUES ($1, $2)", form.ID, form.FriendID)
	return err
}

func (p *postgre) FindUserByID(id string) (*user.User, error) {
	users := make([]*user.User, 0)
	err := p.db.Select(&users, "SELECT id, account, email, name, description, hashedpassword, baned, root, create_at, last_login FROM bloguser WHERE id IN (SELECT friend_id  FROM friend WHERE id=$1) OR id=$1", id)
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

func (p *postgre) FindUserByAccount(acocunt string) (*user.User, error) {
	var usr user.User
	err := p.db.Get(&usr, "SELECT id, account, email, name, description, hashedpassword, baned, root, create_at, last_login FROM bloguser WHERE account=$1", acocunt)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &usr, nil
}

func (p *postgre) UpdateUserPassword(u *user.User) error {
	_, err := p.db.Exec("UPDATE bloguser SET hashedpassword=$1 WHERE id=$2", u.HashedPassword, u.ID)
	return err
}

//update description and name
func (p *postgre) UpdateUser(u *user.User) error {
	_, err := p.db.Exec("UPDATE bloguser SET description=$1, name=$2 WHERE id=$3", u.Description, u.Name, u.ID)
	return err
}

func (p *postgre) UpdateUserStatus(id string, status bool) error {
	_, err := p.db.Exec("UPDATE bloguser SET baned=$1 WHERE id=$2", status, id)
	return err
}

func (p *postgre) DeleteUserByID(id string) error {
	_, err := p.db.Exec("DELETE * FROM bloguser WHERE id=$1", id)
	return err
}

func (p *postgre) FindAllUsers() ([]*user.User, error) {
	users := make([]*user.User, 0)
	err := p.db.Select(&users, "SELECT * FROM bloguser")
	if err != nil {
		return nil, err
	}

	return users, nil
}
