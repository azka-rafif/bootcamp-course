package course

import "github.com/evermos/boilerplate-go/infras"

type CourseRepository interface {
	Create(payload Course) (err error)
}

type CourseRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideCourseRepositoryMySQL(db *infras.MySQLConn) *CourseRepositoryMySQL {
	return &CourseRepositoryMySQL{DB: db}
}

func (*CourseRepositoryMySQL) Create(payload Course) (err error) {

	return
}
