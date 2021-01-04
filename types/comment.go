package types

type Comment struct {
	Post Post `json:"post"`
	User

	CommentID int    `json:"comment_id" db:"comment_id"`
	CommentMessage string `json:"comment_message" db:"comment_message"`
}

