package types

// defines the struct/class user
type User struct {
	UserID   int    `json:"user_id,omitempty" db:"user_id"`
	Username string `json:"username,omitempty" db:"username"`
	Email    string `json:"email,omitempty" db:"email"`
	// inherits UserAuth Struct
	UserAuth
	// inherits UserAttributes Struct
	UserAttributes
}

type UserAuth struct {
	Password    string `json:"password,omitempty" db:"password"`
	Permissions string `json:"permissions,omitempty" db:"permissions"`
}

type UserAttributes struct {
	LocalArea  string `json:"local_area,omitempty" db:"local_area"`
	Reputation int    `json:"reputation,omitempty" db:"reputation"`
}

// method of User
// returns sanitised object without password
func (u *User) Data () User{
	removeP := *u
	removeP.Password = ""
	return removeP
}