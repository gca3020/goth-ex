package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/go-playground/validator/v10"

	"github.com/gca3020/goth-ex/internal/services"
	"github.com/gca3020/goth-ex/internal/templates/pages"
)

type AuthHandlers struct {
	service *services.AuthService
}

func NewAuthHandlers(authService *services.AuthService) *AuthHandlers {
	return &AuthHandlers{
		service: authService,
	}
}

func (h *AuthHandlers) GetJWT() *jwtauth.JWTAuth {
	return tokenAuth
}

func (h *AuthHandlers) Login(r chi.Router) {
	// r.Use() // some middleware..
	r.Get("/", templ.Handler(pages.LoginPage(false, nil)).ServeHTTP)
	r.Post("/", h.handleLoginPost)
}

func (h *AuthHandlers) Signup(r chi.Router) {
	r.Get("/", templ.Handler(pages.LoginPage(true, nil)).ServeHTTP)
	r.Post("/", h.handleSignupPost)
}

func (h *AuthHandlers) Logout(r chi.Router) {
	r.Post("/", h.handleLogoutPost)
}

func makeToken(name string, isAdmin bool) string {
	_, token, _ := tokenAuth.Encode(map[string]interface{}{"name": name, "isAdmin": isAdmin})
	return token
}

type signupData struct {
	Name         string `validate:"min=1,max=64"`
	Email        string `validate:"email"`
	Password     string `validate:"min=4,max=72"`
	Confirmation string `validate:"eqfield=Password"`
}

func (h *AuthHandlers) handleLoginPost(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.LoginUser(r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		// handle errors
		return
	}

	// Create the JWT Token and set the cookie in the response
	token := makeToken(user.Name, user.IsAdmin)
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		// Uncomment below for HTTPS:
		// Secure: true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *AuthHandlers) handleSignupPost(w http.ResponseWriter, r *http.Request) {
	sd := signupData{
		Name:         r.FormValue("name"),
		Email:        r.FormValue("email"),
		Password:     r.FormValue("password"),
		Confirmation: r.FormValue("confirm"),
	}

	// Validate the data from the form, and re-render the form with errors if present
	err := validate.Struct(sd)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		translatedErrors := errs.Translate(translator)
		slog.Warn("Validation Errors", "translatedErrors", translatedErrors)
		templ.Handler(pages.LoginPage(true, translatedErrors)).ServeHTTP(w, r)
		return
	}

	// Add the user
	user, err := h.service.AddUser(sd.Name, sd.Email, sd.Password)
	if err != nil {
		templ.Handler(pages.LoginPage(true, map[string]string{"submit": err.Error()})).ServeHTTP(w, r)
		return
	}

	// Create the JWT Token and set the cookie in the response
	token := makeToken(user.Name, user.IsAdmin)
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		// Uncomment below for HTTPS:
		// Secure: true,
	})

	// Redirect back to the main page
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *AuthHandlers) handleLogoutPost(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		MaxAge:   -1, // Delete the cookie.
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		// Uncomment below for HTTPS:
		// Secure: true,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
