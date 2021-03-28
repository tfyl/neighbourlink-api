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

// retrieve all posts from database (and filter by area if needed)
func RetrieveAllPost(w http.ResponseWriter, r *http.Request, db *db.DB) {

	var search []types.Post
	var err error

	// switch case by the local_area query
	switch filter := r.URL.Query().Get("local_area"); filter{
	 // if there is no filter return all the posts
	case "" :
		search, err = db.GetPostAll()
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			// if there  an error getting data from db, return 400 error
			return
		}

	// if there is a filter then get post by filtering that area, then return that
	default:
		var p types.Post
		p.LocalArea = filter

		search, err = db.GetPostByArea(p)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			// if there  an error getting data from db, return 400 error
			return
		}

	}

	// creates queue
	cQueue := types.NewQueue()
	
	// attach comments each post
	for _,post := range search{
		// gets the comment for the post from database
		comments, err := db.GetCommentsByPost(post)
		if err != nil{
			fmt.Println("httpd: RetrieveAllPost",err)
			continue
		}
		// push to queue
		cQueue.Push(comments)
	}

	queueLen := cQueue.Len()

	for i:=0;i<queueLen;i++{
		// pop from the queue to attach te posts
		search[i].Comments = cQueue.Pop().([]types.Comment)
	}

	var nodes []*alg.Hnode // creates list of nodes heap 
	for _,p := range search{
		nTime := alg.NormaliseTime(search[0].Time,p.Time) // normalises the time to a value where it's the delta between oldest post and current post
		nUrgency := alg.NormalisePriority(p.Urgency) // normalises priority to a value so the multiplier is properly weighted
		nodes = append(nodes, &alg.Hnode{Value: nUrgency*nTime, Data: p})
	}

	// creates new heap
	heap := alg.NewHeap(nodes)
	heap.Sort()  // sorts the object using the public method
	nodes = heap.ReturnArray() // returns array of nodes

	var OrderedP []types.Post // instantiates the object that holds the ordered posts
	for _,n:= range nodes{
		OrderedP = append(OrderedP, n.Data.(types.Post) ) // adds the posts to the variable 
	}
	// returns json
	_ = json.NewEncoder(w).Encode(OrderedP)
}


// retrieve one singular post
func RetrievePost(w http.ResponseWriter, r *http.Request, db *db.DB) {
	// get post id from url parameter
	PostIdStr := chi.URLParam(r, "PostID")
	// convert post id to int
	PostId,err := strconv.Atoi(PostIdStr)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
		_, _ = w.Write([]byte(fmt.Sprintf("404 Not Found: %s is a string", PostIdStr)))
		return
	}

	// search for the post using post id
	SearchPost, err := db.GetPost(types.Post{PostID: PostId})
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// gets comment associated with the post
	comments, err := db.GetCommentsByPost(SearchPost)
	if err != nil{
		fmt.Println("httpd: RetrieveAllPost",err)
	}
	SearchPost.Comments = comments

	_ = json.NewEncoder(w).Encode(SearchPost)
}


// function to create the post
func CreatePost(w http.ResponseWriter, r *http.Request, db *db.DB) {
	// instantiates the object post to a variable called p
	var p types.Post

	// get the JWT user id from the cookie (coded in middleware)
	JWTID := middleware.JWTUserID(r)

	// decode the new post details
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// sets user id to the authenticated user_id in the JWT cookie
	p.UserID = JWTID

	// add post to the db
	p, err = db.AddPost(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// runs code to send webhook to all active users
	go alertPost(p,db)

	_ = json.NewEncoder(w).Encode(p)

}

// function to update the post
func UpdatePost(w http.ResponseWriter, r *http.Request, db *db.DB) {
	// get the JWT user id from the cookie (coded in middleware)
	JWTID := middleware.JWTUserID(r)

	// gets post id from the url parameter
	PostIdStr := chi.URLParam(r, "PostID")
	// converts the string into int
	PostId,err := strconv.Atoi(PostIdStr)

	// gets existing post details
	p, err := db.GetPost(types.Post{PostID: PostId})
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// if the authenticated user does not equal the user who posted it, return an error
	if p.UserID != JWTID{
		http.Error(w,fmt.Sprintf("UserID does not match logged in user Post:%d JWT:%d",p.UserID,JWTID),http.StatusUnauthorized)
		return
	}

	// get updated post details
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// sets post user id
	p.UserID = JWTID
	// sets post id
	p.PostID = PostId

	// updates the post in the database using post id
	p, err = db.UpdatePost(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// runs code to send webhook to all active users
	go alertPost(p,db)

	_ = json.NewEncoder(w).Encode(p)
}

func alertPost(p types.Post,db *db.DB){
	// checks if the post urgency is 5
	if p.Urgency == 5{
		// sends post to all connected clients
		err := db.Websocket.SendAll(p)
		fmt.Println("Error sending websocket",err)
	}
}