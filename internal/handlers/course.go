package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/evermos/boilerplate-go/internal/domain/course"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/shared/pagination"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
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
		r.Get("/", h.HandleGetAll)
		r.Group(func(r chi.Router) {
			r.Use(h.JwtAuth.CheckRole)
			r.Post("/", h.HandleCreate)
		})
	})
}

func (h *CourseHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {

	page, err := pagination.ConvertToInt(pagination.ParseQueryParams(r, "page"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	limit, err := pagination.ConvertToInt(pagination.ParseQueryParams(r, "limit"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	sort := pagination.GetSortDirection(pagination.ParseQueryParams(r, "sort"))
	field := pagination.CheckFieldQuery(pagination.ParseQueryParams(r, "field"), "id")
	title := pagination.ParseQueryParams(r, "title")
	offset := (page - 1) * limit
	claims, ok := r.Context().Value(middleware.ClaimsKey("claims")).(jwt.Claims)
	if !ok {
		response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userId, err := uuid.FromString(claims.UserId)
	if err != nil {
		response.WithError(w, err)
		return
	}
	res, err := h.Service.GetAll(limit, offset, sort, field, title, userId)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, res)
}

func (h *CourseHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload course.CoursePayload
	err := decoder.Decode(&payload)
	err = shared.GetValidator().Struct(payload)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	claims, ok := r.Context().Value(middleware.ClaimsKey("claims")).(jwt.Claims)
	if !ok {
		response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userId, err := uuid.FromString(claims.UserId)
	if err != nil {
		response.WithError(w, failure.InternalError(err))
	}
	res, err := h.Service.CreateCourse(payload, userId)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusCreated, res)
}
