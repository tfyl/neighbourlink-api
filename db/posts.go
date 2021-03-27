package db

import (
	"neighbourlink-api/types"
	"time"
)

func (db *DB) AddPost (post types.Post) (types.Post,error) {

	tx := db.MustBegin()
	_ ,err  := tx.Exec(`INSERT INTO post(user_id,post_time,post_title,post_description,post_urgency) VALUES ($1,$2,$3,$4,$5) `,post.UserID,time.Now(),post.Title,post.Description,post.Urgency)
	if err != nil {
		return post,err
	}
	err = tx.Commit()
	return post,err
}


func (db *DB) UpdatePost (post types.Post) (types.Post,error) {

	tx := db.MustBegin()
	_ ,err  := tx.Exec(
		`UPDATE post 
				SET 
				user_id = $1,
				post_title = $2,
				post_description = $3,
				post_urgency = $4
				WHERE post_id = $5;`,

				post.UserID,
				post.Title,
				post.Description,
				post.Urgency,
				post.PostID,
				)
	if err != nil {
		return post,err
	}

	err = tx.Commit()
	return post,err
}


func (db *DB) GetPostAll () ([]types.Post, error) {

	var posts []types.Post
	// creates variable posts which is an array of the type types.Post

	err := db.Select(&posts, `
	SELECT
	       post.post_id,
	       post.user_id,
	       user_detail.username,
	       post.post_time,
	       post.post_title,
	       post.post_description,
	       post.post_urgency,
	       user_attribute.local_area,
		   user_attribute.reputation
	FROM
	     post
	
	INNER JOIN user_detail ON post.user_id=user_detail.user_id
	INNER JOIN user_attribute ON post.user_id=user_attribute.user_id

	ORDER BY post.post_time DESC;
	`, )// Returns posts ordered by time descending
	// Performs inner join to get the corresponding record from user_detail and user_attribute

	if err != nil {
		return posts,err
	}



	return posts,err

}


func (db *DB) GetPostByArea (post types.Post) ([]types.Post, error) {


	var posts []types.Post
	// creates variable posts which is an array of the type types.Post

	err := db.Select(&posts, `
	SELECT
	       post.post_id,
	       post.user_id,
	       user_detail.username,
	       post.post_time,
	       post.post_title,
	       post.post_description,
	       post.post_urgency,
	       user_attribute.local_area,
		   user_attribute.reputation
	FROM
	     post
	
	INNER JOIN user_detail ON post.user_id=user_detail.user_id
	INNER JOIN user_attribute ON post.user_id=user_attribute.user_id
	
	WHERE
		user_attribute.local_area=$1

	ORDER BY post.post_time DESC;
	`, post.LocalArea) // Returns posts ordered by time descending
	// Performs inner join to get the corresponding record from user_detail and user_attribute
	if err != nil {
		return posts,err
	}

	return posts,err

}


func (db *DB) GetPost (post types.Post) (types.Post, error) {
	// One singular post using post_id
	tx := db.MustBegin()
	err := tx.Get(&post, `
	SELECT
	       post.post_id,
	       post.user_id,
	       user_detail.username,
	       post.post_time,
	       post.post_title,
	       post.post_description,
	       post.post_urgency,
	       user_attribute.local_area,
	       user_attribute.reputation
	FROM
	     post
	
	INNER JOIN user_detail ON post.user_id=user_detail.user_id
	INNER JOIN user_attribute ON post.user_id=user_attribute.user_id
	
	WHERE
		post.post_id=$1;
	`, post.PostID) // Gets post by checking post_id
	// Performs inner join to get the corresponding record from user_detail and user_attribute
	if err != nil {
		return post,err
	}
	err = tx.Commit()


	return post,err
}

