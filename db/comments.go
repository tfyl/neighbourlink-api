package db

import (
	"fmt"
	"neighbourlink-api/types"
)

func (db *DB) AddComment (comment types.Comment) (types.Comment,error) {

	tx := db.MustBegin()
	_ ,err  := tx.Exec(`INSERT INTO post_comment(post_id,user_id,comment_message) VALUES ($1,$2,$3) `,comment.Post.PostID,comment.UserID,comment.CommentMessage)
	if err != nil {
		return comment, err
	}

	err = tx.Commit()
	return comment, err
}


func (db *DB) GetCommentAll () ([]types.Comment, error) {

	var comments []types.Comment

	tx := db.MustBegin()
	tx1 := db.MustBegin()
	defer tx.Commit()
	defer tx1.Commit()

	rows, err := tx.Queryx(`
	SELECT
		comment_id,
	    user_id,
	    comment_message
	FROM
		post_comment
	;`)
	if err != nil{
		fmt.Println(1)
		return nil,err
	}

	defer rows.Close()
	for rows.Next() {
		var c types.Comment
		err = rows.StructScan(&c)
		if err != nil{
			fmt.Println(2)
			return nil,err
		}
		err = tx1.QueryRowx(`
		SELECT
			post_id
		FROM
		    post_comment
		WHERE
			post_comment.comment_id=$1
		;`, c.CommentID).StructScan(&c.Post)
		if err != nil{
			fmt.Println(3)
			return nil,err
		}

		comments = append(comments,c )
	}

	return comments,nil
}

func (db *DB) GetCommentsByPost (post types.Post) ([]types.Comment, error) {

	var comments []types.Comment


	err := db.Select(&comments,`
	SELECT
		comment_id,
	    user_id,
	    comment_message
	FROM
		post_comment
	WHERE
		post_comment.post_id=$1
	;`,post.PostID)

	if err != nil{
		return nil,err
	}

	for i,_ := range comments{
		comments[i].Post.PostID = post.PostID
	}

	fmt.Println(comments)
	return comments,nil
}




func (db *DB) GetComment (comment types.Comment) (types.Comment, error) {

	err := db.Get(&comment, `
	SELECT
		*  
	FROM
	     post_comment
	WHERE
		post_comment.comment_id=$1
	;`, comment.CommentID)


	return comment,err

}




func (db *DB) UpdateComment (comment types.Comment) (types.Comment,error) {

	tx := db.MustBegin()
	_ ,err  := tx.Exec(
		`UPDATE post_comment 
				SET 
				comment_message = $1
				WHERE
				comment_id = $2`,
				comment.CommentID,
	)
	if err != nil {
		return comment,err
	}

	err = tx.Commit()
	return comment,err
}

