package content

import (
	"inside-athletics/internal/s3"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

// Registers premium content (S3 upload/download URL) routes; s3Service must be non-nil.
func Route(api huma.API, db *gorm.DB, s3Service *s3.Service) {
	if s3Service == nil {
		return
	}
	svc := NewContentService(s3Service)
	grp := huma.NewGroup(api, "/api/v1/content")
	huma.Post(grp, "/upload-url", svc.GetUploadURL)
	huma.Get(grp, "/download-url", svc.GetDownloadURL)
	huma.Post(grp, "/confirm-upload", svc.ConfirmUpload)
	huma.Delete(grp, "", svc.DeleteContent)
}
