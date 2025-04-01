package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"cognito-example/config"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	config    *config.Config
	oauthConf *oauth2.Config
	provider  *oidc.Provider
}

func NewAuthHandler(cfg *config.Config) (*AuthHandler, error) {
	provider, err := oidc.NewProvider(context.Background(), cfg.CognitoIssuerURL)
	if err != nil {
		return &AuthHandler{}, err
	}
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.CognitoClientID,
		ClientSecret: cfg.CognitoClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes: []string{
			"openid",
			"email",
			"phone",
		},
		Endpoint: provider.Endpoint(),
	}

	return &AuthHandler{
		config:    cfg,
		oauthConf: oauthConfig,
		provider:  provider,
	}, nil
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	url := h.oauthConf.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := h.oauthConf.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// You can now use the token to make authenticated requests
	response := map[string]interface{}{
		"access_token":  token.AccessToken,
		"token_type":    token.TokenType,
		"refresh_token": token.RefreshToken,
		"expiry":        token.Expiry,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		http.Error(w, "No authorization header", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token := &oauth2.Token{
		AccessToken: tokenString,
	}

	userInfo, err := h.provider.UserInfo(r.Context(), oauth2.StaticTokenSource(token))
	if err != nil {
		http.Error(w, fmt.Errorf("Failed to get user info: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	var claims map[string]interface{}
	if err := userInfo.Claims(&claims); err != nil {
		http.Error(w, fmt.Errorf("Failed to parse user info: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(claims)
}
