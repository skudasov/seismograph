package server

const (
	SingleBlobAllowedErr = "only single csv blob is allowed"
	MultipartReadErr     = "failed to read multipart file"
	NoMetaFieldErr       = "no meta field in multipart, must be json"
	MetaParsingErr       = "error parsing meta field"
	EmptyBlobErr         = "data blob is empty"
	MetaInfoCreationErr  = "error creating meta info"
	TestDataCreationErr  = "error putting blob"

	// cluster
	UnsupportedProviderErr = "provider is not supported"
	UnsupportedImageErr    = "image is not supported"
)

type ErrorResponse struct {
	Errors []string `json:"errors"`
}

func DefaultErrorResponse(error string) ErrorResponse {
	return ErrorResponse{
		Errors: []string{error},
	}
}
