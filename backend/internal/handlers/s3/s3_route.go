package s3

import (
	"context"
	"inside-athletics/internal/s3"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	cfg, ok := s3.LoadConfigFromEnv()
	if !ok {

	}

	client, err := s3.NewClient(context.Background(), cfg)
	if err != nil {

	}

	s3Service := s3.NewService(client, cfg)
	{
		grp := huma.NewGroup(api, "/api/v1/s3")
		huma.Get(grp, "/upload-url", s3Service.GetUploadURL)
		huma.Get(grp, "/download-url", s3Service.GetDownloadURL)
	}
}
