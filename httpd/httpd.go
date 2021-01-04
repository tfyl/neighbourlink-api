package httpd

import (
	"github.com/go-chi/chi"
	"log"
	"neighbourlink-api/db"
	"net/http"
)

type Server struct {
	db *db.DB
	secretKey []byte
}

func New(db *db.DB) Server {
	return Server{
		db:db,
		secretKey:[]byte("askldjoifbaiusndoia"),
	}
}


func (s Server) Start () {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/post", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				// Public routes
				r.Get("/", func(w http.ResponseWriter, r *http.Request) { RetrieveAllPost(w, r, s.db) })
				r.Get("/{PostID}", func(w http.ResponseWriter, r *http.Request) { RetrievePost(w, r, s.db) })
				// Private routes
				r.Group(func(r chi.Router) {
					// Handle valid / invalid tokens
					r.Use(func(handler http.Handler) http.Handler { return JWTAuthMiddleware(handler, s.secretKey) })
					r.Post("/", func(w http.ResponseWriter, r *http.Request) { CreatePost(w, r, s.db) })          // POST  - create a new post
					r.Patch("/{PostID}", func(w http.ResponseWriter, r *http.Request) { UpdatePost(w, r, s.db) }) // POST  - Update post
				})

			})
		})
		r.Route("/comment", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				// Public routes
				r.Get("/", func(w http.ResponseWriter, r *http.Request) { RetrieveAllComments(w, r, s.db) })
				r.Get("/{CommentID}", func(w http.ResponseWriter, r *http.Request) { RetrieveComment(w, r, s.db) })
				// Private routes
				r.Group(func(r chi.Router) {
					// Handle valid / invalid tokens
					r.Use(func(handler http.Handler) http.Handler { return JWTAuthMiddleware(handler, s.secretKey) })
					r.Post("/", func(w http.ResponseWriter, r *http.Request) { CreateComment(w, r, s.db) })          // POST  - create a new comment
					r.Patch("/{CommentID}", func(w http.ResponseWriter, r *http.Request) { UpdateComment(w, r, s.db) }) // POST  - Update post
				})

			})
		})
		r.Route("/user", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request){ CreateUser(w,r,s.db) })
			r.Group(func(r chi.Router) {
				// Public routes
				r.Use(func(handler http.Handler) http.Handler { return JWTAuthMiddleware(handler, s.secretKey) })
				r.Get("/", func(w http.ResponseWriter, r *http.Request){ RetrieveUser(w,r,s.db) })
				r.Get("/{Username}", func(w http.ResponseWriter, r *http.Request){ RetrieveUser(w,r,s.db) })

			})
		})
		r.Route("/session", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request){ LoginUser(w,r,s.db,s.secretKey) })
		})

	})

	log.Fatal(http.ListenAndServe(":80", r))

//	r.Route("/api/v1" ,func(r chi.Router) {
//		r.Route("/post", func(r chi.Router) {
//			r.Get("/", func(w http.ResponseWriter, r *http.Request) { post.RequestPost(w, r, s.db) }) // change func to get all posts
//			r.Post("/", func(w http.ResponseWriter, r *http.Request) { post.AddPost(w, r, db, secretKey) })
//			r.Get("/{PostID}", func(w http.ResponseWriter, r *http.Request) { post.RequestPost(w, r, db) })
//			r.Put("/{PostID}", func(w http.ResponseWriter, r *http.Request) { post.RequestPost(w, r, db) }) // change func to update
//			r.Get("/{PostID}", func(w http.ResponseWriter, r *http.Request) { post.RequestPost(w, r, db) }) // change func to get a specific post
//		})
//		r.Route("/account", func(r chi.Router) {
//			r.Get("/", func(w http.ResponseWriter, r *http.Request) { auth.AccountInfo(w, r, db,secretKey) })
//			r.Post("/signup", func(w http.ResponseWriter, r *http.Request) { auth.Signup(w, r, db) })
//			r.Post("/login", func(w http.ResponseWriter, r *http.Request) { auth.Login(w, r, db, secretKey) })
//			r.Post("/testJWT", func(w http.ResponseWriter, r *http.Request) { auth.ValidateJWT(w, r, secretKey) })
//		})
//		r.Route("/area", func(r chi.Router) {
//			r.Post("/calcPath", func(w http.ResponseWriter, r *http.Request) { calcpath.ShortestPath(w, r) })
//		})
//	})
}