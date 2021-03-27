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
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS","PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// redirects any urls to the correct frontend
	r.Get("/", func(w http.ResponseWriter, r *http.Request){
		http.Redirect(w, r, "/web/posts.html", http.StatusMovedPermanently)
	})
	r.Get(`/{^[a-zA-Z]*\.html}`, func(w http.ResponseWriter, r *http.Request){
		http.Redirect(w, r, "/web/posts.html", http.StatusMovedPermanently)
	})
	r.Get(`/web/`, func(w http.ResponseWriter, r *http.Request){
		http.Redirect(w, r, "/web/posts.html", http.StatusMovedPermanently)
	})

	// api routes
	r.Route("/api/v1", func(r chi.Router) {
		// REST endpoint for post
		r.Route("/post", func(r chi.Router) {
			// Public routes
			r.Group(func(r chi.Router) {
				// GET - retrieves all the posts
				r.Get("/", func(w http.ResponseWriter, r *http.Request) { RetrieveAllPost(w, r, s.db) })
				// GET - retrieves the post with the specific post id
				r.Get("/{PostID}", func(w http.ResponseWriter, r *http.Request) { RetrievePost(w, r, s.db) })

				// Private routes (requires the user to be logged in)
				r.Group(func(r chi.Router) {
					// Checks JWT value by running the middleware
					r.Use(func(handler http.Handler) http.Handler { return middleware.JWTAuthMiddleware(handler, s.secretKey) })
					// POST  - create a new post
					r.Post("/", func(w http.ResponseWriter, r *http.Request) { CreatePost(w, r, s.db) })
					// PATCH  - Update post
					r.Patch("/{PostID}", func(w http.ResponseWriter, r *http.Request) { UpdatePost(w, r, s.db) })
				})

			})
		})
		r.Route("/graph", func(r chi.Router) {
			// Public routes
			r.Group(func(r chi.Router) {
				// GET - retrieves the shortest (safest) path between two points
				r.Get("/", func(w http.ResponseWriter, r *http.Request) { RetrieveShortestPath(w, r, s.db) })
			})
		})
		r.Route("/comment", func(r chi.Router) {
			// Public routes
			r.Group(func(r chi.Router) {
				// GET - retrieves all the comments
				r.Get("/", func(w http.ResponseWriter, r *http.Request) { RetrieveAllComments(w, r, s.db) })
				// GET - retrieves the comment with the specific comment id
				r.Get("/{CommentID}", func(w http.ResponseWriter, r *http.Request) { RetrieveComment(w, r, s.db) })

				// Private routes
				r.Group(func(r chi.Router) {
					// Checks JWT value by running the middleware
					r.Use(func(handler http.Handler) http.Handler { return middleware.JWTAuthMiddleware(handler, s.secretKey) })
					// POST  - create a new comment
					r.Post("/", func(w http.ResponseWriter, r *http.Request) { CreateComment(w, r, s.db) })
					// PATCH  - Update post
					r.Patch("/{CommentID}", func(w http.ResponseWriter, r *http.Request) { UpdateComment(w, r, s.db) })
				})

			})
		})
		// Public routes
		r.Route("/user", func(r chi.Router) {

			r.Post("/", func(w http.ResponseWriter, r *http.Request){ CreateUser(w,r,s.db) })
			r.Get("/{UserID}", func(w http.ResponseWriter, r *http.Request){ RetrieveUser(w,r,s.db) })

			// Private routes
			r.Group(func(r chi.Router) {
				// Checks JWT value by running the middleware
				r.Use(func(handler http.Handler) http.Handler { return middleware.JWTAuthMiddleware(handler, s.secretKey) })
				// GET - uses JWT middleware to determine authorised user and return data for that specific user
				r.Get("/", func(w http.ResponseWriter, r *http.Request){ RetrieveUser(w,r,s.db) })
				// PATCH - updates user and thus must be protected by validating JWT Cookies
				r.Patch("/{UserID}", func(w http.ResponseWriter, r *http.Request){ UpdateUser(w,r,s.db) })
				// DELETE - delete user (must have auth level of `Admin`)
				r.Delete("/{UserID}", func(w http.ResponseWriter, r *http.Request){ admin.DeleteUser(w,r,s.db) })
			})
		})

		// endpoint to handle sessions - public so users can authenticate
		r.Route("/session", func(r chi.Router) {
			// POST - create session (login) using auth
			r.Post("/", func(w http.ResponseWriter, r *http.Request){ LoginUser(w,r,s.db,s.secretKey) })
		})

		// endpoint for admin functions - private
		r.Route("/admin", func(r chi.Router) {
			// Checks JWT value by running the middleware
			r.Use(func(handler http.Handler) http.Handler { return middleware.JWTAuthMiddleware(handler, s.secretKey) })
			// GET - Gets all users so they can perform admin functions
			r.Get("/user", func(w http.ResponseWriter, r *http.Request){ admin.RetrieveAllUsers(w,r,s.db) })
		})
		r.Route("/ws", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				// Checks JWT value by running the middleware
				r.Use(func(handler http.Handler) http.Handler { return middleware.JWTAuthMiddleware(handler, s.secretKey) })
				r.Get("/", func(w http.ResponseWriter, r *http.Request) { serveWs(w, r, s.db) })
			})
		})

	})

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "./static"))
	fmt.Println(filesDir)
	FileServer(r, "/web/", filesDir)
	log.Fatal(http.ListenAndServe(":8080", r))
}

// FileServer sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		fmt.Println("FileServer does not permit any URL parameters.")
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