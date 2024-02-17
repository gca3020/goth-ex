package handlers

import (
	"os"

	"github.com/go-chi/jwtauth"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	uni        *ut.UniversalTranslator
	validate   *validator.Validate
	translator ut.Translator

	tokenAuth *jwtauth.JWTAuth
)

func Init() {
	// Build the locale and register it as the default (fallback) locale
	en := en.New()
	uni = ut.New(en, en)

	// Construct the validator that our handlers will use
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Build a translator and register the default validation translations into it
	translator, _ = uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, translator)

	// Construct the JWT Auth
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenAuth = jwtauth.New("HS256", []byte(jwtSecret), nil)
}
