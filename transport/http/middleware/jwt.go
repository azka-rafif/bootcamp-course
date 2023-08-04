package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/transport/http/response"
)

type JwtAuthentication struct {
	config *configs.Config
}

type ClaimsKey string

const (
	HeaderJwt = "Authorization"
)

type ResponseJWT struct {
	Data jwt.Claims `json:"data"`
}

type CustomError struct {
	Error string `json:"error"`
}

func ProvideJwtAuthentication(conf *configs.Config) *JwtAuthentication {
	return &JwtAuthentication{config: conf}
}

func (a *JwtAuthentication) CheckJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(HeaderJwt)
		if token == "" {
			response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		req, err := http.NewRequest("GET", a.config.App.AuthBaseURL+"/validate", nil)
		if err != nil {
			response.WithError(w, err)
			return
		}
		req.Header.Set("Authorization", token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			response.WithError(w, err)
			return
		}
		if resp.StatusCode == 401 {
			decoder := json.NewDecoder(resp.Body)
			var payload CustomError
			err = decoder.Decode(&payload)
			if err != nil {
				response.WithError(w, err)
				return
			}
			response.WithError(w, failure.Unauthorized(strings.Split(payload.Error, ": ")[1]))
			return
		}
		decoder := json.NewDecoder(resp.Body)
		var payload ResponseJWT
		err = decoder.Decode(&payload)
		if err != nil {
			response.WithError(w, err)
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsKey("claims"), payload.Data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *JwtAuthentication) CheckRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey("claims")).(jwt.Claims)
		if !ok {
			response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if !strings.EqualFold(claims.Role, "teacher") {
			response.WithError(w, failure.Unauthorized("those who are not teachers, are unauthorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
