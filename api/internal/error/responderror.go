package error

import (
	"net/http"
	serverJSON "virsh-sandbox/internal/json"
)

type errorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

func RespondError(w http.ResponseWriter, status int, err error) {
	_ = serverJSON.RespondJSON(w, status, errorResponse{
		Error: err.Error(),
		Code:  status,
		// details intentionally omitted to avoid leaking internals; add as needed
	})
}
