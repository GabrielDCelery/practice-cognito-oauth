package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type contextKey string

const (
	CookieDataKey contextKey = "cookie_data"
)

type CookieConfig struct {
	Name     string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}

func DefaultCookieConfig() CookieConfig {
	return CookieConfig{
		MaxAge:   3600,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func SetCookie(w http.ResponseWriter, name string, value interface{}, config CookieConfig) {
	valueAsBytes, err := json.Marshal(value)
	if err != nil {
		http.Error(w, fmt.Sprintf("something unexpected happened"), http.StatusInternalServerError)
		return
	}
	cookie := &http.Cookie{
		Name:     name,
		Value:    string(valueAsBytes),
		MaxAge:   config.MaxAge,
		Path:     config.Path,
		Domain:   config.Domain,
		Secure:   config.Secure,
		HttpOnly: config.HttpOnly,
		SameSite: config.SameSite,
		Expires:  time.Now().Add(time.Duration(config.MaxAge) * time.Second),
	}
	http.SetCookie(w, cookie)
}

func GetCookie(r *http.Request, name string) (*http.Cookie, error) {
	return r.Cookie(name)
}

func DeleteCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func CreateAuthMiddleware() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := GetCookie(r, "profile")
			if err != nil {
				if err == http.ErrNoCookie {
					http.Redirect(w, r, "/login", http.StatusFound)
					return
				}
				http.Error(w, "something unexpected happened", http.StatusInternalServerError)
			}
			var cookieData interface{}
			if err := json.Unmarshal([]byte(cookie.Value), cookieData); err != nil {
				http.Error(w, "something unexpected happened", http.StatusInternalServerError)
				return
			}
			r = SetCookieDataToContext(r, cookieData)
			next(w, r)
		}
	}
}

func GetCookieDataFromContext(r *http.Request) interface{} {
	return r.Context().Value(CookieDataKey)
}

func SetCookieDataToContext(r *http.Request, data interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), CookieDataKey, data))
}
