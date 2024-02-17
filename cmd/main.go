package main

import (
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/jwtauth"
	_ "github.com/joho/godotenv/autoload"

	"github.com/gca3020/goth-ex/assets"
	"github.com/gca3020/goth-ex/internal/handlers"
	"github.com/gca3020/goth-ex/internal/services"
	"github.com/gca3020/goth-ex/internal/store"
	"github.com/gca3020/goth-ex/internal/templates/pages"
)

func main() {
	logger := httplog.NewLogger("example-app")
	addr := os.Getenv("SERVE_ADDR")

	// Initialize and Construct the Database Layer
	userStore := store.NewMemUserStore()

	// Initialize and Construct the Service Layer
	authService := services.NewAuthService(userStore)

	// Initialize and Construct the Handler Layer
	handlers.Init()
	authHandler := handlers.NewAuthHandlers(authService)

	// Set up routing and Middleware
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(jwtauth.Verifier(authHandler.GetJWT()))

	// Set up routes
	r.Get("/", templ.Handler(pages.HomePage()).ServeHTTP)
	r.Route("/login", authHandler.Login)
	r.Route("/signup", authHandler.Signup)
	r.Route("/logout", authHandler.Logout)

	// Add a route for static assets
	r.Handle("/assets/*", assets.GetServer())

	// Start serving
	http.ListenAndServe(addr, r)
}
