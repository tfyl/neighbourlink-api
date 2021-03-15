package httpd

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"neighbourlink-api/alg"
	"neighbourlink-api/db"
	"neighbourlink-api/httpd/middleware"
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

	cQueue := types.NewQueue()

	// attach comments
	for _,post := range search{
		comments, err := db.GetCommentsByPost(post)
		if err != nil{
			fmt.Println("httpd: RetrieveAllPost",err)
			continue
		}

		cQueue.Push(comments)

	}

	queueLen := cQueue.Len()

	for i:=0;i<queueLen;i++{
		search[i].Comments = cQueue.Pop().([]types.Comment)
	}

	var nodes []*alg.Hnode
	for _,p := range search{
		nTime := alg.NormaliseTime(search[0].Time,p.Time)
		nUrgency := alg.NormalisePriority(p.Urgency)
		nodes = append(nodes, &alg.Hnode{Value: nUrgency*nTime, Data: p})
	}

	// creates new heap
	heap := alg.NewHeap(nodes)
	heap.Sort()  // sorts the object using the public method
	nodes = heap.ReturnArray() // returns array of nodes

	var OrderedP []types.Post
	for _,n:= range nodes{
		OrderedP = append(OrderedP, n.Data.(types.Post) )
	}

	_ = json.NewEncoder(w).Encode(OrderedP)
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

	JWTID := middleware.JWTUserID(r)


	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.UserID = JWTID

	p, err = db.AddPost(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	_ = json.NewEncoder(w).Encode(p)

}

func UpdatePost(w http.ResponseWriter, r *http.Request, db *db.DB) {
	JWTID := middleware.JWTUserID(r)

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

	p, err = db.UpdatePost(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	_ = json.NewEncoder(w).Encode(p)
}
