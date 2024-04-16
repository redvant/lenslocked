package cookies

import (
	"fmt"
	"net/http"
)

const (
	CookieSession = "session"
)

func NewCookie(name, value string) *http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
	return &cookie
}

func SetCookie(w http.ResponseWriter, name, value string) {
	cookie := NewCookie(name, value)
	http.SetCookie(w, cookie)
}

func ReadCookie(r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("%s: %w", name, err)
	}
	return c.Value, nil
}

func DeleteCookie(w http.ResponseWriter, name string) {
	cookie := NewCookie(name, "")
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}
