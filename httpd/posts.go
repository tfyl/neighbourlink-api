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

func RetrieveAllPost(w http.ResponseWriter, r *http.Request, db *db.DB) {


	var search []types.Post
	var err error

	switch filter := r.URL.Query().Get("local_area"); filter{
	case "" :
		search, err = db.GetPostAll()
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}


	default:
		var p types.Post
		p.LocalArea = filter

		search, err = db.GetPostByArea(p)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}


	// attach comments
	for i,post := range search{
		comments, err := db.GetCommentsByPost(post)
		if err != nil{
			fmt.Println("httpd: RetrieveAllPost",err)
			continue
		}
		search[i].Comments = comments
	}

	_ = json.NewEncoder(w).Encode(search)
}


func RetrievePost(w http.ResponseWriter, r *http.Request, db *db.DB) {
	PostIdStr := chi.URLParam(r, "PostID")
	PostId,err := strconv.Atoi(PostIdStr)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
		_, _ = w.Write([]byte(fmt.Sprintf("404 Not Found: %s is a string", PostIdStr)))
		return
	}

	SearchPost, err := db.GetPost(types.Post{PostID: PostId})
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comments, err := db.GetCommentsByPost(SearchPost)
	if err != nil{
		fmt.Println("httpd: RetrieveAllPost",err)
	}
	SearchPost.Comments = comments

	_ = json.NewEncoder(w).Encode(SearchPost)
}



func CreatePost(w http.ResponseWriter, r *http.Request, db *db.DB) {
	var p types.Post

	JWTID := JWTUserID(r)


	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.UserID = JWTID

	_, err = db.AddPost(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

}

func UpdatePost(w http.ResponseWriter, r *http.Request, db *db.DB) {
	JWTID := JWTUserID(r)

	PostIdStr := chi.URLParam(r, "PostID")
	PostId,err := strconv.Atoi(PostIdStr)


	p, err := db.GetPost(types.Post{PostID: PostId})
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if p.UserID != JWTID{
		http.Error(w,fmt.Sprintf("UserID does not match logged in user Post:%d JWT:%d",p.UserID,JWTID),http.StatusUnauthorized)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.UserID = JWTID
	p.PostID = PostId

	_, err = db.UpdatePost(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

}
