package main

import (
	"log"
	"net/http"
	"github.com/XIAOR1Ck/go-redis-voting/internal/voting"
		"github.com/XIAOR1Ck/go-redis-voting/internal/adapters/redis"
"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"os"
	"path/filepath"

	"github.com/go-chi/cors"
)

func CORSMiddleware() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		// Development
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		},

		// Production examples:
		// AllowedOrigins: []string{
		//     "https://app.example.com",
		//     "https://admin.example.com",
		// },

		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},

		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"Origin",
		},

		ExposedHeaders: []string{
			"Content-Length",
			"Content-Type",
		},

		// Set true only if using cookies/sessions.
		AllowCredentials: true,

		// Cache preflight requests for 5 minutes.
		MaxAge: 300,
	})
}

func main(){
	// Initialize router
	r := chi.NewRouter()
	
	// Add some useful middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(CORSMiddleware())	
	
	// Initialize services
	store := voting.NewRedisStore(redis.NewClient("your_redis_data_here"))
	svc := voting.NewService(store)
	handler := voting.NewHandler(svc)
 	// Auth middleware instance (implemented in service)
	authMiddleware := svc.AuthMiddleware
		
	// Static files for QR codes (public - no auth)
wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	qrDir := filepath.Join(wd, "QrCodes")

	if err := os.MkdirAll(qrDir, 0755); err != nil {
		log.Fatal(err)
	}

	log.Printf("Serving QR codes from: %s", qrDir)

	fs := http.FileServer(http.Dir(qrDir))

	r.Handle("/qrcodes/*", http.StripPrefix("/qrcodes/", fs))

	// Optional: redirect /qrcodes -> /qrcodes/
	r.Get("/qrcodes", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/qrcodes/", http.StatusMovedPermanently)
	})	// Public routes (no authentication required)
	r.Group(func(r chi.Router) {
		// These endpoints are publicly accessible
		r.Get("/getprojects", handler.GetProjects)           // Public
		r.Post("/castvote", handler.CastVote)                // Public
		r.Post("/expocreateadmin", handler.CreateAdmin)      // Public
		r.Post("/admin/login", handler.AdminLogin)           // Public
		r.Get("/bulkadd", handler.LoadBulk)
	})
	
	// Protected routes (authentication required)
	r.Group(func(r chi.Router) {
		// Apply auth middleware to this group
		r.Use(authMiddleware)
		
		// These endpoints require authentication
		r.Post("/addproject", handler.CreateProject)
		r.Post("/adduser", handler.CreateUser)
		r.Get("/getuser", handler.GetUser)
	})

		// Start server
	if err := http.ListenAndServe(":9090", r); err != nil {
		log.Fatal(err)
	}
	}


