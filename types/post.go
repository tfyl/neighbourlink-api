package types

import "time"

type Post struct {
	PostID int    `json:"post_id,omitempty" db:"post_id"`
	User
	Title       string `json:"post_title,omitempty" db:"post_title"`
	Description string `json:"post_description,omitempty" db:"post_description"`
	Urgency     int    `json:"post_urgency,omitempty" db:"post_urgency"`
	Time   *time.Time `json:"post_time,omitempty" db:"post_time"`
	Comments []Comment `json:"comments,omitempty"`
}