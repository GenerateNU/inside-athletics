package models

import (
	"time"

	"github.com/google/uuid"
)

type Sport struct {
	ID         uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `sql:"index" json:"deleted_at"`
	Name       string     `json:"name" example:"Women's Soccer" doc:"The name of a sport" gorm:"type:varchar(100);not null"`
	Popularity int32      `json:"popularity" example:"20000" doc:"How many people play that sport in the U.S." gorm:"type:int32;"`
}

func (s *Sport) SetPopularity(popularity int32) *Sport {
	s.Popularity = popularity
	return s
}

func (s *Sport) setName(name string) *Sport {
	s.Name = name
	return s
}
