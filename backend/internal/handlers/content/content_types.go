package content

// Request body for requesting a presigned upload URL.
type GetUploadURLRequest struct {
	Key      string `json:"key" required:"true" doc:"Full S3 object key, e.g. premium/image/content-123/photo.jpg"`
	FileType string `json:"fileType" required:"true" doc:"MIME type, e.g. image/jpeg, application/pdf"`
	FileName string `json:"fileName" doc:"Optional; used as documentId in response (defaults to last segment of key)"`
}

// Input for POST /upload-url (body only).
type GetUploadURLInput struct {
	Body GetUploadURLRequest
}

// Query params for requesting a presigned download URL.
type GetDownloadURLParams struct {
	Key string `query:"key" required:"true" doc:"S3 object key"`
}

// Request body for confirming an upload (verify object exists, get download URL).
type ConfirmUploadRequest struct {
	Key string `json:"key" required:"true" doc:"S3 object key from upload-url response"`
}

// Input for POST /confirm-upload.
type ConfirmUploadInput struct {
	Body ConfirmUploadRequest
}

// Query params for deleting content by key.
type DeleteContentParams struct {
	Key string `query:"key" required:"true" doc:"S3 object key to delete"`
}

// Response after deleting content.
type DeleteContentResponse struct {
	Message string `json:"message" doc:"Success message"`
}
