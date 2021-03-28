package types

// class / struct that defines Comment that inherits attributes from "User"
type Comment struct {
	Post           Post `json:"post"`
	User           // inherits from "User" stuct / class
	CommentID      int    `json:"comment_id" db:"comment_id"`
	CommentMessage string `json:"comment_message" db:"comment_message"`
}

