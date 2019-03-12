package forms

type InsertFriend struct {
	ID       string
	FriendID string
}

type CreateUser struct {
	Account     string `json:"account"`
	Name        string `json:"name"`
	EMail       string `json:"email"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

type LoginForm struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type UpdatePassword struct {
	UserID      string
	PrePassword string `json:"prePassword"`
	Password    string `json:"password"`
}

type UpdateUserStatus struct {
	UserID string
	Status bool
}
