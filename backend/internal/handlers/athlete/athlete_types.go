package athlete

import "github.com/google/uuid"

type VerifyAthleteParam struct {
	Name    string     `query:"name" required:"true" maxLength:"200" example:"Joe Adams" doc:"Name of athlete to verify"`
	College *uuid.UUID `query:"college" required:"true" doc:"uuid for college to verify user at"`
	Sport   *uuid.UUID `query:"sport" required:"true" doc:"uuid for sport to verify user at"`
}

type VerifyAthleteResponse struct {
	Verified         bool   `json:"verified" doc:"boolean represented if the athlete was verified" example:"true"`
	Name             string `json:"name" doc:"name of the athlete" example:"Hannah Bang"`
	CollegeName      string `json:"college_name" doc:"Name of the college athlete goes to"`
	AthleticsWebsite string `json:"athletics_website" doc:"Athletics website for the given college" example:"https//:www.huskies.edu"`
}
