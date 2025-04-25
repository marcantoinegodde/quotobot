package server

import (
	"context"
	"embed"
	"errors"
	"net/http"
	"quotobot/pkg/config"
	"quotobot/pkg/database"
	"quotobot/pkg/logger"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

//go:embed templates/*.html
var templates embed.FS

type Server struct {
	Logger   *logger.Logger
	Config   *config.Config
	Database *gorm.DB
}

func NewServer() *Server {
	l := logger.NewLogger()
	c := config.LoadConfig(l)
	db := database.LoadDatabase(l)

	return &Server{
		Logger:   l,
		Config:   c,
		Database: db,
	}
}

func (s *Server) Start() {
	ctx := context.Background()

	store := sessions.NewCookieStore([]byte(s.Config.Server.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // TODO: set to true in production
		MaxAge:   3600,
	}

	provider, err := oidc.NewProvider(ctx, s.Config.Server.ProviderURL)
	if err != nil {
		s.Logger.Error.Fatalf("Failed to create OIDC provider: %v", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     s.Config.Server.ClientID,
		ClientSecret: s.Config.Server.ClientSecret,
		RedirectURL:  s.Config.Server.RedirectURL,

		Endpoint: provider.Endpoint(),

		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: s.Config.Server.ClientID})

	http.HandleFunc("/oauth/authorize", s.AuthorizeHandler(store, oauth2Config))
	http.HandleFunc("/oauth/callback", s.CallbackHandler(ctx, store, oauth2Config, verifier))
	http.HandleFunc("/register", s.RegisterHandler(store))

	s.Logger.Info.Println("Server listening on :8080")

	http.ListenAndServe(":8080", nil)
}

func (s *Server) AuthorizeHandler(store *sessions.CookieStore, oauth2Config *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := generateState()
		codeVerifier := oauth2.GenerateVerifier()

		session, _ := store.Get(r, "session")

		session.Values["state"] = state
		session.Values["verifier"] = codeVerifier

		if err := session.Save(r, w); err != nil {
			s.Logger.Error.Printf("Failed to save session: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, oauth2Config.AuthCodeURL(state, oauth2.S256ChallengeOption(codeVerifier)), http.StatusFound)
	}
}

func (s *Server) CallbackHandler(ctx context.Context, store *sessions.CookieStore, oauth2Config *oauth2.Config, verifier *oidc.IDTokenVerifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify state
		session, _ := store.Get(r, "session")

		if r.URL.Query().Get("state") != session.Values["state"] {
			s.Logger.Error.Println("Invalid state parameter.")
			http.Error(w, "invalid state", http.StatusBadRequest)
			return
		}

		// Exchange code for token
		codeVerifier, ok := session.Values["verifier"].(string)
		if !ok {
			s.Logger.Error.Println("Invalid verifier parameter.")
			http.Error(w, "invalid verifier", http.StatusBadRequest)
			return
		}

		oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"), oauth2.VerifierOption(codeVerifier))
		if err != nil {
			s.Logger.Error.Printf("Failed to exchange token: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// Extract the ID Token from OAuth2 token
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			s.Logger.Error.Println("No id_token field in oauth2 token.")
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// Parse and verify ID Token payload
		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			s.Logger.Error.Printf("Failed to verify ID Token: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// Extract custom claims
		var claims struct {
			Email     string `json:"email"`
			Name      string `json:"name"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Avatar    string `json:"avatar"`
		}
		if err := idToken.Claims(&claims); err != nil {
			s.Logger.Error.Printf("Failed to parse ID Token claims: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// Store user information in session
		u := User{
			ID:        idToken.Subject,
			Email:     claims.Email,
			Name:      claims.Name,
			FirstName: claims.FirstName,
			LastName:  claims.LastName,
			Avatar:    claims.Avatar,
		}

		session.Values["authenticated"] = true
		session.Values["user"] = u
		if err := session.Save(r, w); err != nil {
			s.Logger.Error.Printf("Failed to save session: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// Redirect to the register page
		http.Redirect(w, r, "/register", http.StatusFound)
	}
}

func (s *Server) RegisterHandler(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ok := validateURL(r.URL.Query(), s.Config.Server.HMACSecret); !ok {
			s.renderTemplate(w, []string{"templates/register.html"}, RegisterTemplateData{Status: "error"})
			return
		}

		username := r.URL.Query().Get("username")
		id := r.URL.Query().Get("id")

		session, _ := store.Get(r, "session")

		authenticated, ok := session.Values["authenticated"].(bool)
		if !ok || !authenticated {
			s.renderTemplate(w, []string{"templates/register.html"}, RegisterTemplateData{Username: username, Status: "unauthenticated"})
			return
		}

		user, ok := session.Values["user"].(User)
		if !ok {
			s.renderTemplate(w, []string{"templates/register.html"}, RegisterTemplateData{Username: username, Status: "unauthenticated"})
			return
		}

		u := database.User{
			TelegramID: id,
			ViaRezoID:  user.ID,
		}

		if err := s.Database.Create(&u).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				s.renderTemplate(w, []string{"templates/register.html"}, RegisterTemplateData{FirstName: user.FirstName, Status: "already_registered"})
				return
			} else {
				s.Logger.Error.Printf("Failed to create user in database: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		}

		s.Logger.Info.Printf("User registered: %s - %s", u.ViaRezoID, u.TelegramID)
		s.renderTemplate(w, []string{"templates/register.html"}, RegisterTemplateData{FirstName: user.FirstName, Status: "success"})
	}
}
