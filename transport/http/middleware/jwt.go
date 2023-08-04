package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/transport/http/response"
)

type JwtAuthentication struct {
}

type ClaimsKey string

const (
	HeaderJwt = "Authorization"
)

type ResponseJWT struct {
	Data jwt.Claims `json:"data"`
}

func ProvideJwtAuthentication() *JwtAuthentication {
	return &JwtAuthentication{}
}

func (a *JwtAuthentication) CheckJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(HeaderJwt)
		if token == "" {
			response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		req, err := http.NewRequest("GET", "http://127.0.0.1:8080/v1/auth/validate", nil)
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

		if resp.StatusCode != 200 {
			response.WithMessage(w, http.StatusBadRequest, "Invalid Jwt Token")
			return
		}
		decoder := json.NewDecoder(resp.Body)
		var payload ResponseJWT
		err = decoder.Decode(&payload)
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsKey("claims"), payload)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
