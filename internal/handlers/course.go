package handlers

import (
	"net/http"

	"github.com/evermos/boilerplate-go/internal/domain/course"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
)

type CourseHandler struct {
	Service course.CourseService
	JwtAuth *middleware.JwtAuthentication
}

func ProvideCourseHandler(service course.CourseService, jwtAuth *middleware.JwtAuthentication) CourseHandler {
	return CourseHandler{Service: service, JwtAuth: jwtAuth}
}

func (h *CourseHandler) Router(r chi.Router) {
	r.Route("/course", func(r chi.Router) {
		r.Use(h.JwtAuth.CheckJwt)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			response.WithMessage(w, http.StatusOK, "ok")
		})
		r.Post("/", h.HandleCreate)
	})
}

func (h CourseHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {

}
