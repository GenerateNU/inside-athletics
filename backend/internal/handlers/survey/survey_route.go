package survey

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	svc := NewSurveyService(db)

	{
		grp := huma.NewGroup(api, "/api/v1/survey")
		huma.Post(grp, "/", svc.CreateSurvey)                          // POST   /api/v1/survey/          — submit a survey
		huma.Delete(grp, "/{id}", svc.DeleteSurvey)                    // DELETE /api/v1/survey/{id}      — delete a survey
		huma.Get(grp, "/user/{user_id}", svc.GetSurveysByUser)         // GET    /api/v1/survey/user/{id} — own user's surveys
		huma.Get(grp, "/averages", svc.GetAverageRatings)              // GET    /api/v1/survey/averages  — averages (sport/college filters)
	}
}