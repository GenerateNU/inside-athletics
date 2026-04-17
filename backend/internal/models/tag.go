
package models

import (
	"time"

	"github.com/google/uuid"
)

type TagType string

const (
    TagTypeSports               TagType = "sports"
    TagTypeSchools              TagType = "schools"
    TagTypeDivisions            TagType = "divisions"
    TagTypeAthleticsPerformance TagType = "athletics_performance"
    TagTypeHealthWellness       TagType = "health_wellness"
    TagTypeStudentAthleteLife   TagType = "student_athlete_life"
    TagTypeRecruitingLogistics  TagType = "recruiting_logistics"
)

type Tag struct {
    ID        uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `sql:"index" json:"deleted_at"`
    Name      string     `json:"name" example:"Hockey" doc:"The name of the tag" gorm:"type:varchar(100);not null"`
    Type      TagType    `json:"type" example:"sports" doc:"The type of the tag" gorm:"type:varchar(50);not null"`
}
