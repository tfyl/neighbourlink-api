package db

import (
	"neighbourlink-api/types"
)

func (db *DB) AddPost (post types.Post) (types.Post,error) {

	tx := db.MustBegin()
	_ ,err  := tx.Exec(`INSERT INTO post(user_id,post_title,post_description,post_urgency) VALUES ($1,$2,$3,$4) `,post.UserID,post.Title,post.Description,post.Urgency)
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
				WHERE post_id = $5`,

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
	err := db.Select(&posts, `
	SELECT
	       post.post_id,
	       post.user_id,
	       post.post_title,
	       post.post_description,
	       post.post_urgency,
	       user_attribute.local_area,
		   user_attribute.reputation
	FROM
	     post
	
	INNER JOIN user_attribute ON post.user_id=user_attribute.user_id
	`, )
	if err != nil {
		return posts,err
	}



	return posts,err

}


func (db *DB) GetPostByArea (post types.Post) ([]types.Post, error) {


	var posts []types.Post

	err := db.Select(&posts, `
	SELECT
	       post.post_id,
	       post.user_id,
	       post.post_title,
	       post.post_description,
	       post.post_urgency,
	       user_attribute.local_area,
		   user_attribute.reputation
	FROM
	     post
	
	INNER JOIN user_attribute ON post.user_id=user_attribute.user_id
	
	WHERE
		user_attribute.local_area=$1;
	`, post.LocalArea)
	if err != nil {
		return posts,err
	}

	return posts,err

}


func (db *DB) GetPost (post types.Post) (types.Post, error) {

	tx := db.MustBegin()
	err := tx.Get(&post, `
	SELECT
	       post.post_id,
	       post.user_id,
	       post.post_title,
	       post.post_description,
	       post.post_urgency,
	       user_attribute.local_area,
	       user_attribute.reputation
	FROM
	     post
	
	INNER JOIN user_attribute ON post.user_id=user_attribute.user_id
	
	WHERE
		post.post_id=$1;
	`, post.PostID)
	if err != nil {
		return post,err
	}
	err = tx.Commit()


	return post,err
}

