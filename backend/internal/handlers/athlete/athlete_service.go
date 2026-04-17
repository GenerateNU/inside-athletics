package athlete

import (
	"context"
	"inside-athletics/internal/handlers/college"
	"inside-athletics/internal/utils"

	"gorm.io/gorm"
)

type AthleteService struct {
	athleteDB *AthleteDB
	collegeDB *college.CollegeService
}

func NewAthleteService(db *gorm.DB) *AthleteService {
	return &AthleteService{
		athleteDB: NewAthleteDB(db),
	}
}

func (a *AthleteService) VerifyAthlete(context context.Context, input *VerifyAthleteParam) (*utils.ResponseBody[VerifyAthleteResponse], error) {
	athlete, found, err := a.athleteDB.GetAthlete(input.Name, input.College, input.Sport)
	if err != nil {
		return nil, err
	}
	return &utils.ResponseBody[VerifyAthleteResponse]{
		Body: &VerifyAthleteResponse{
			Verified:         found,
			Name:             athlete.Name,
			CollegeName:      athlete.College.Name,
			AthleticsWebsite: athlete.College.AthleticsWebsite,
		},
	}, nil
}
