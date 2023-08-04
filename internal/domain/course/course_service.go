package course

import "github.com/gofrs/uuid"

type CourseService interface {
	CreateCourse(payload CoursePayload, userId uuid.UUID) (res Course, err error)
}

type CourseServiceImpl struct {
	Repo CourseRepository
}

func ProvideCourseServiceImpl(repo CourseRepository) *CourseServiceImpl {
	return &CourseServiceImpl{Repo: repo}
}

func (s *CourseServiceImpl) CreateCourse(payload CoursePayload, userId uuid.UUID) (res Course, err error) {
	res, err = res.NewFromPayload(payload, userId)
	return
}
