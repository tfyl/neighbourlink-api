package httpd

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"neighbourlink-api/db"
	"neighbourlink-api/httpd/middleware"
	"neighbourlink-api/types"
	"net/http"
	"strconv"
)

func RetrieveAllComments(w http.ResponseWriter, r *http.Request, db *db.DB) {
	// returns all comments in json form

	search, err := db.GetCommentAll() // gets list of comments
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest) // if there is an error return 400 error
		return
	}

	// if all is successful return all comments
	_ = json.NewEncoder(w).Encode(search)
}


func RetrieveComment(w http.ResponseWriter, r *http.Request, db *db.DB) {

	var c types.Comment
	err := json.NewDecoder(r.Body).Decode(&c)
	CommentIDstr := chi.URLParam(r, "CommentID")
	CommentID,err := strconv.Atoi(CommentIDstr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.CommentID = CommentID

	c, err = db.GetComment(c)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(c)
}


func CreateComment(w http.ResponseWriter, r *http.Request, db *db.DB) {
	var c types.Comment
	JWTID := middleware.JWTUserID(r)

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.UserID = JWTID
	fmt.Println(c)
	c, err = db.AddComment(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	_ = json.NewEncoder(w).Encode(c)
}


func UpdateComment(w http.ResponseWriter, r *http.Request, db *db.DB) {
	JWTID := middleware.JWTUserID(r)

	CommentIDStr := chi.URLParam(r, "CommentID")
	CommentID,err := strconv.Atoi(CommentIDStr)

	var c types.Comment

	err = json.NewDecoder(r.Body).Decode(&c)

	if c.UserID != JWTID{
		http.Error(w,fmt.Sprintf("UserID does not match logged in user Post:%d JWT:%d",c.UserID,JWTID),http.StatusUnauthorized)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.UserID = JWTID
	c.CommentID = CommentID


	_, err = db.UpdateComment(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	_ = json.NewEncoder(w).Encode(c)

}
