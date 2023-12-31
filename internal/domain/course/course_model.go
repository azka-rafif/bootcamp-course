package course

import (
	"encoding/json"
	"time"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type Course struct {
	Id         uuid.UUID   `db:"id" validate:"true"`
	UserId     uuid.UUID   `db:"user_id" validate:"true"`
	Title      string      `db:"title" validate:"true"`
	Content    string      `db:"content" validate:"true"`
	Created_at time.Time   `db:"created_at" validate:"true"`
	Updated_at time.Time   `db:"updated_at" validate:"true"`
	Deleted_at null.Time   `db:"deleted_at"`
	Created_by uuid.UUID   `db:"created_by"`
	Updated_by uuid.UUID   `db:"updated_by"`
	Deleted_by nuuid.NUUID `db:"deleted_by"`
}

type CoursePayload struct {
	Title   string `json:"title" validate:"true"`
	Content string `json:"content" validate:"true"`
}

type CourseResponseFormat struct {
	Id         uuid.UUID   `json:"id"`
	UserId     uuid.UUID   `json:"userId"`
	Title      string      `json:"title"`
	Content    string      `json:"content"`
	Created_at time.Time   `json:"createdAt"`
	Updated_at time.Time   `json:"updatedAt"`
	Deleted_at null.Time   `json:"deletedAt,omitempty"`
	Created_by uuid.UUID   `json:"createdBy"`
	Updated_by uuid.UUID   `json:"updatedBy"`
	Deleted_by nuuid.NUUID `json:"deletedBy,omitempty"`
}

func (c *Course) NewFromPayload(load CoursePayload, userId uuid.UUID) (Course, error) {
	courseId, err := uuid.NewV4()
	if err != nil {
		return Course{}, err
	}
	newCours := Course{
		Id:         courseId,
		UserId:     userId,
		Title:      load.Title,
		Content:    load.Content,
		Created_at: time.Now().UTC(),
		Created_by: userId,
		Updated_at: time.Now().UTC(),
		Updated_by: userId,
	}
	err = newCours.Validate()
	return newCours, err
}

func (c Course) ToResponseFormat() CourseResponseFormat {
	resp := CourseResponseFormat(c)
	return resp
}

func (c Course) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.ToResponseFormat())
}

func (c *Course) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(c)
}
