package schema

import "github.com/opentreehole/backend/common/sensitive"

type SensitiveCheckRequest struct {
	Content string `json:"content" validate:"required"`
}

type SensitiveCheckResponse struct {
	sensitive.ResponseForCheck
}
