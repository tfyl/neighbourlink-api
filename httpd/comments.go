package httpd

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"neighbourlink-api/db"
	"neighbourlink-api/types"
	"net/http"
	"strconv"
)

func RetrieveAllComments(w http.ResponseWriter, r *http.Request, db *db.DB) {

	search, err := db.GetCommentAll()
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(search)
}


func RetrieveComment(w http.ResponseWriter, r *http.Request, db *db.DB) {

	var c types.Comment
	err := json.NewDecoder(r.Body).Decode(&c)


	search, err := db.GetComment(c)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	_ = json.NewEncoder(w).Encode(search)
}




func CreateComment(w http.ResponseWriter, r *http.Request, db *db.DB) {
	var c types.Comment
	JWTID := JWTUserID(r)

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.UserID = JWTID

	_, err = db.AddComment(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

}


func UpdateComment(w http.ResponseWriter, r *http.Request, db *db.DB) {
	JWTID := JWTUserID(r)

	CommentIDStr := chi.URLParam(r, "CommentID")
	CommentID,err := strconv.Atoi(CommentIDStr)

	var c types.Comment

	//if c.UserID != JWTID{
	//	http.Error(w,fmt.Sprintf("UserID does not match logged in user Post:%d JWT:%d",p.UserID,JWTID),http.StatusUnauthorized)
	//	return
	//}

	err = json.NewDecoder(r.Body).Decode(&c)
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

}








// Not part of a RESTful API
//
//func RetrieveCommentsByPost(w http.ResponseWriter, r *http.Request, db *db.DB) {

//	var p types.Post
//	err := json.NewDecoder(r.Body).Decode(&p)
//
//
//	search, err := db.GetCommentsByPost(p)
//	if err != nil {
//		fmt.Println(err)
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	_ = json.NewEncoder(w).Encode(search)
//}
