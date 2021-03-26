package httpd

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"log"
	"neighbourlink-api/db"
	"neighbourlink-api/httpd/admin"
	"neighbourlink-api/httpd/middleware"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS","PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request){
		http.Redirect(w, r, "/web/posts.html", http.StatusMovedPermanently)
	})
	r.Get(`/{^[a-zA-Z]*\.html}`, func(w http.ResponseWriter, r *http.Request){
		http.Redirect(w, r, "/web/posts.html", http.StatusMovedPermanently)
	})
	r.Get(`/web/`, func(w http.ResponseWriter, r *http.Request){
		http.Redirect(w, r, "/web/posts.html", http.StatusMovedPermanently)
	})


	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/post", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				// Public routes
				r.Get("/", func(w http.ResponseWriter, r *http.Request) { RetrieveAllPost(w, r, s.db) })
				r.Get("/{PostID}", func(w http.ResponseWriter, r *http.Request) { RetrievePost(w, r, s.db) })
				// Private routes
				r.Group(func(r chi.Router) {
					// Handle valid / invalid tokens
					r.Use(func(handler http.Handler) http.Handler { return middleware.JWTAuthMiddleware(handler, s.secretKey) })
					r.Post("/", func(w http.ResponseWriter, r *http.Request) { CreatePost(w, r, s.db) })          // POST  - create a new post
					r.Patch("/{PostID}", func(w http.ResponseWriter, r *http.Request) { UpdatePost(w, r, s.db) }) // POST  - Update post
				})

			})
		})
		r.Route("/graph", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				// Public routes
				r.Get("/", func(w http.ResponseWriter, r *http.Request) { RetrieveShortestPath(w, r, s.db) })
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
					r.Use(func(handler http.Handler) http.Handler { return middleware.JWTAuthMiddleware(handler, s.secretKey) })
					r.Post("/", func(w http.ResponseWriter, r *http.Request) { CreateComment(w, r, s.db) })          // POST  - create a new comment
					r.Patch("/{CommentID}", func(w http.ResponseWriter, r *http.Request) { UpdateComment(w, r, s.db) }) // PATCH  - Update post
				})

			})
		})
		r.Route("/user", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request){ CreateUser(w,r,s.db) })
			r.Get("/{UserID}", func(w http.ResponseWriter, r *http.Request){ RetrieveUser(w,r,s.db) })  // doesn't use JWT Middleware
			r.Group(func(r chi.Router) {
				r.Use(func(handler http.Handler) http.Handler { return middleware.JWTAuthMiddleware(handler, s.secretKey) })

				r.Get("/", func(w http.ResponseWriter, r *http.Request){ RetrieveUser(w,r,s.db) }) // uses JWT middleware to determine authorised user and return data
				r.Patch("/{UserID}", func(w http.ResponseWriter, r *http.Request){ UpdateUser(w,r,s.db) }) // updates user and thus must be protected by validating JWT Cookies
				r.Delete("/{UserID}", func(w http.ResponseWriter, r *http.Request){ admin.DeleteUser(w,r,s.db) }) // delete user (must have auth level of `Admin`)
			})
		})
		r.Route("/session", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request){ LoginUser(w,r,s.db,s.secretKey) })
		})
		r.Route("/admin", func(r chi.Router) {
			r.Use(func(handler http.Handler) http.Handler { return middleware.JWTAuthMiddleware(handler, s.secretKey) })
			r.Get("/user", func(w http.ResponseWriter, r *http.Request){ admin.RetrieveAllUsers(w,r,s.db) })
		})
		r.Route("/ws", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(func(handler http.Handler) http.Handler { return middleware.JWTAuthMiddleware(handler, s.secretKey) })
				r.Get("/", func(w http.ResponseWriter, r *http.Request) { serveWs(w, r, s.db) })
			})
		})

	})

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "./static"))
	fmt.Println(filesDir)
	FileServer(r, "/web/", filesDir)
	log.Fatal(http.ListenAndServe(":81", r))

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

// FileServer sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}