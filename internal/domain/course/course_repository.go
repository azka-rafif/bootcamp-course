package course

import (
	"fmt"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type CourseRepository interface {
	Create(payload Course) (err error)
	GetAll(limit, offset int, sort, field, title string, userId uuid.UUID) (res []Course, err error)
}

type CourseRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideCourseRepositoryMySQL(db *infras.MySQLConn) *CourseRepositoryMySQL {
	return &CourseRepositoryMySQL{DB: db}
}

func (r *CourseRepositoryMySQL) Create(payload Course) (err error) {
	return r.DB.WithTransaction(func(db *sqlx.Tx, c chan error) {
		if err := r.txCreate(db, payload); err != nil {
			c <- err
			return
		}
		c <- nil
	})
}

func (r *CourseRepositoryMySQL) txCreate(tx *sqlx.Tx, payload Course) (err error) {
	query := `INSERT INTO course (id,user_id,title,content,created_at,created_by,updated_at,updated_by)
    VALUES (:id,:user_id,:title,:content,:created_at,:created_by,:updated_at,:updated_by)`

	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(payload)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *CourseRepositoryMySQL) GetAll(limit, offset int, sort, field, title string, userId uuid.UUID) (res []Course, err error) {
	query := `SELECT * FROM course WHERE `
	query += fmt.Sprintf(" user_id = '%s' ", userId.String())

	if title != "" {
		query += `AND title LIKE `
		query += fmt.Sprintf("'%%%s%%'", title)
	}
	query += fmt.Sprintf(" ORDER BY %s %s LIMIT %d OFFSET %d", field, sort, limit, offset)
	println(query)
	err = r.DB.Read.Select(&res, query)
	if err != nil {
		err = failure.InternalError(err)
		logger.ErrorWithStack(err)
	}
	return
}
