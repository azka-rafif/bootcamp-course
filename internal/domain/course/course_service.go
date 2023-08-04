package course

type CourseService interface {
}

type CourseServiceImpl struct {
	Repo CourseRepository
}

func ProvideCourseServiceImpl(repo CourseRepository) *CourseServiceImpl {
	return &CourseServiceImpl{Repo: repo}
}
