package server

import (
	"mime/multipart"

	"github.com/f4hrenh9it/seismograph/back/types"
)

func ValidateAddDataRequest(mpf *multipart.Form) string {
	// TODO: validate from js
	if _, ok := mpf.Value["meta"]; !ok {
		return NoMetaFieldErr
	}
	if len(mpf.File[types.MPartDataKey]) > 1 {
		return SingleBlobAllowedErr
	}
	if mpf.File[types.MPartDataKey][0].Size == 0 {
		return EmptyBlobErr
	}
	return ""
}

func ValidateClusterSpec(spec AttackClusterRequest) string {
	if spec.ProviderName != "aws" {
		return UnsupportedProviderErr
	}
	return ""
}
