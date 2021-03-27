package db

import (
	"fmt"
	"neighbourlink-api/types"
)

func (db *DB) AddComment (comment types.Comment) (types.Comment,error) {
	// adds comment to database
	tx := db.MustBegin()
	// starts a transaction so data collisions don't occur
	_ ,err  := tx.Exec(`INSERT INTO post_comment(post_id,user_id,comment_message) VALUES ($1,$2,$3) `,comment.Post.PostID,comment.UserID,comment.CommentMessage)
	if err != nil {
		return comment, err
	}

	// commits the transaction
	err = tx.Commit()
	return comment, err
}


func (db *DB) GetCommentAll () ([]types.Comment, error) {
	// gets all comments for all posts
	var comments []types.Comment

	tx := db.MustBegin() // starts 2 transactions for both processes that are carried out
	tx1 := db.MustBegin()
	defer tx.Commit()  // defers the commit (defer means that the command executes after the function returns and is finished)
	defer tx1.Commit()


	rows, err := tx.Queryx(` 
	SELECT
		comment_id,
	    user_id,
	    comment_message
	FROM
		post_comment
	;`)
	// this query gets all the comments
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
		// this query attaches the post_id to the comment as it's stored in another child object

		if err != nil{
			fmt.Println(3)
			return nil,err
		}

		comments = append(comments,c )
	}

	return comments,nil
}

func (db *DB) GetCommentsByPost (post types.Post) ([]types.Comment, error) {
	// gets all comments specific to a post
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
	// gets all comments that have the post_id desired

	if err != nil{
		return nil,err
	}

	for i,_ := range comments{
		comments[i].Post.PostID = post.PostID
	}
	// attaches post_id to the comments

	return comments,nil
}




func (db *DB) GetComment (c types.Comment) (types.Comment, error) {


	tx := db.MustBegin()
	defer tx.Commit()

	err := tx.QueryRowx(`
	SELECT
		comment_id,
	    user_id,
	    comment_message
	FROM
		post_comment
	WHERE
		comment_id = $1
	;`,c.CommentID).StructScan(&c)
	if err != nil{
		//fmt.Println(1)
		return c,err
	}

	err = tx.QueryRowx(`
		SELECT
			post_id
		FROM
		    post_comment
		WHERE
			post_comment.comment_id=$1
		;`, c.CommentID).StructScan(&c.Post)
	if err != nil{
		//fmt.Println(3)
		return c,err
	}

	return c,nil

}




func (db *DB) UpdateComment (comment types.Comment) (types.Comment,error) {

	tx := db.MustBegin()
	_ ,err  := tx.Exec(
		`UPDATE post_comment 
				SET 
				comment_message = $1
				WHERE
				comment_id = $2`,
				comment.CommentMessage,
				comment.CommentID,
	)
	if err != nil {
		return comment,err
	}

	err = tx.Commit()
	return comment,err
}

