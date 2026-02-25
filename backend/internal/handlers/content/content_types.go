package content

// Request body for requesting a presigned upload URL.
type GetUploadURLRequest struct {
	FileName    string `json:"fileName" required:"true" doc:"Original filename"`
	FileType    string `json:"fileType" required:"true" doc:"MIME type, e.g. image/jpeg, application/pdf"`
	ContentKind string `json:"contentKind" required:"true" doc:"image, video, or pdf"`
	ContentID   string `json:"contentId" doc:"Optional; preferred for key path if set"`
	UserID      string `json:"userId" doc:"Optional; used for key path if ContentID empty"`
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
