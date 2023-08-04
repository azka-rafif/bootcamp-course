package course

import "github.com/gofrs/uuid"

type CourseService interface {
	CreateCourse(payload CoursePayload, userId uuid.UUID) (res Course, err error)
	GetAll(limit, offset int, sort, field, title string, userId uuid.UUID) (res []Course, err error)
}

type CourseServiceImpl struct {
	Repo CourseRepository
}

func ProvideCourseServiceImpl(repo CourseRepository) *CourseServiceImpl {
	return &CourseServiceImpl{Repo: repo}
}

func (s *CourseServiceImpl) CreateCourse(payload CoursePayload, userId uuid.UUID) (res Course, err error) {
	res, err = res.NewFromPayload(payload, userId)
	if err != nil {
		return
	}
	err = s.Repo.Create(res)
	if err != nil {
		return
	}
	return
}

func (s *CourseServiceImpl) GetAll(limit, offset int, sort, field, title string, userId uuid.UUID) (res []Course, err error) {
	res, err = s.Repo.GetAll(limit, offset, sort, field, title, userId)
	if err != nil {
		return
	}
	return
}
