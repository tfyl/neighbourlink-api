package types

type Post struct {
	PostID int    `json:"post_id" db:"post_id"`
	User
	Title       string `json:"post_title,omitempty" db:"post_title"`
	Description string `json:"post_description,omitempty" db:"post_description"`
	Urgency     int    `json:"post_urgency,omitempty" db:"post_urgency"`
//	Timestamp   time.Time `json:"timestamp,omitempty" db:"timestamp"`
	Comments []Comment `json:"comments,omitempty"`
}