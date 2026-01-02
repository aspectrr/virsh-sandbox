package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	serverError "virsh-sandbox/internal/error"
	serverJSON "virsh-sandbox/internal/json"
	"virsh-sandbox/internal/sshca"
)

// AccessHandler handles SSH access API requests.
type AccessHandler struct {
	accessSvc *sshca.AccessService
}

// NewAccessHandler creates a new access handler.
func NewAccessHandler(accessSvc *sshca.AccessService) *AccessHandler {
	return &AccessHandler{
		accessSvc: accessSvc,
	}
}

// RegisterRoutes registers the access routes on the given router.
func (h *AccessHandler) RegisterRoutes(r chi.Router) {
	r.Route("/access", func(r chi.Router) {
		// Request access to a sandbox
		r.Post("/request", h.handleRequestAccess)

		// Get CA public key for VM configuration
		r.Get("/ca-pubkey", h.handleGetCAPublicKey)

		// Certificate operations
		r.Route("/certificate/{certID}", func(r chi.Router) {
			r.Get("/", h.handleGetCertificate)
			r.Delete("/", h.handleRevokeCertificate)
		})

		// List certificates
		r.Get("/certificates", h.handleListCertificates)

		// Session operations
		r.Post("/session/start", h.handleRecordSessionStart)
		r.Post("/session/end", h.handleRecordSessionEnd)

		// List active sessions
		r.Get("/sessions", h.handleListSessions)
	})
}

// Request/Response types

// requestAccessRequest is the request body for requesting sandbox access.
type requestAccessRequest struct {
	// SandboxID is the target sandbox.
	SandboxID string `json:"sandbox_id"`

	// UserID identifies the requesting user.
	UserID string `json:"user_id"`

	// PublicKey is the user's SSH public key in OpenSSH format.
	PublicKey string `json:"public_key"`

	// TTLMinutes is the requested access duration (1-10 minutes).
	TTLMinutes int `json:"ttl_minutes,omitempty"`
}

// requestAccessResponse is the response for a successful access request.
type requestAccessResponse struct {
	// CertificateID is the ID of the issued certificate.
	CertificateID string `json:"certificate_id"`

	// Certificate is the SSH certificate content (save as key-cert.pub).
	Certificate string `json:"certificate"`

	// VMIPAddress is the IP address of the sandbox VM.
	VMIPAddress string `json:"vm_ip_address"`

	// SSHPort is the SSH port (usually 22).
	SSHPort int `json:"ssh_port"`

	// Username is the SSH username to use.
	Username string `json:"username"`

	// ValidUntil is when the certificate expires (RFC3339).
	ValidUntil string `json:"valid_until"`

	// TTLSeconds is the remaining validity in seconds.
	TTLSeconds int `json:"ttl_seconds"`

	// ConnectCommand is an example SSH command for connecting.
	ConnectCommand string `json:"connect_command"`

	// Instructions provides usage instructions.
	Instructions string `json:"instructions"`
}

// caPublicKeyResponse is the response for getting the CA public key.
type caPublicKeyResponse struct {
	// PublicKey is the CA public key in OpenSSH format.
	PublicKey string `json:"public_key"`

	// Usage explains how to use this key.
	Usage string `json:"usage"`
}

// certificateResponse is the response for certificate queries.
type certificateResponse struct {
	ID           string   `json:"id"`
	SandboxID    string   `json:"sandbox_id"`
	UserID       string   `json:"user_id"`
	VMID         string   `json:"vm_id"`
	Identity     string   `json:"identity"`
	SerialNumber uint64   `json:"serial_number"`
	Principals   []string `json:"principals"`
	ValidAfter   string   `json:"valid_after"`
	ValidBefore  string   `json:"valid_before"`
	Status       string   `json:"status"`
	IssuedAt     string   `json:"issued_at"`
	IsExpired    bool     `json:"is_expired"`
	TTLSeconds   int      `json:"ttl_seconds,omitempty"`
}

// listCertificatesResponse is the response for listing certificates.
type listCertificatesResponse struct {
	Certificates []certificateResponse `json:"certificates"`
	Total        int                   `json:"total"`
}

// revokeCertificateRequest is the request body for revoking a certificate.
type revokeCertificateRequest struct {
	Reason string `json:"reason"`
}

// sessionStartRequest is the request body for recording session start.
type sessionStartRequest struct {
	CertificateID string `json:"certificate_id"`
	SourceIP      string `json:"source_ip,omitempty"`
}

// sessionStartResponse is the response for session start.
type sessionStartResponse struct {
	SessionID string `json:"session_id"`
}

// sessionEndRequest is the request body for recording session end.
type sessionEndRequest struct {
	SessionID string `json:"session_id"`
	Reason    string `json:"reason,omitempty"`
}

// sessionResponse is the response for session queries.
type sessionResponse struct {
	ID              string `json:"id"`
	CertificateID   string `json:"certificate_id"`
	SandboxID       string `json:"sandbox_id"`
	UserID          string `json:"user_id"`
	VMID            string `json:"vm_id"`
	VMIPAddress     string `json:"vm_ip_address"`
	SourceIP        string `json:"source_ip,omitempty"`
	Status          string `json:"status"`
	StartedAt       string `json:"started_at"`
	EndedAt         string `json:"ended_at,omitempty"`
	DurationSeconds int    `json:"duration_seconds,omitempty"`
}

// listSessionsResponse is the response for listing sessions.
type listSessionsResponse struct {
	Sessions []sessionResponse `json:"sessions"`
	Total    int               `json:"total"`
}

// accessErrorResponse is a helper for error responses.
type accessErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

// writeError writes an error response.
func writeError(w http.ResponseWriter, status int, message, details string) {
	_ = serverJSON.RespondJSON(w, status, accessErrorResponse{
		Error:   message,
		Code:    status,
		Details: details,
	})
}

// Handlers

// handleRequestAccess handles POST /v1/access/request
// @Summary Request SSH access to a sandbox
// @Description Issues a short-lived SSH certificate for accessing a sandbox via tmux
// @Tags Access
// @Accept json
// @Produce json
// @Param request body requestAccessRequest true "Access request"
// @Success 200 {object} requestAccessResponse
// @Failure 400 {object} accessErrorResponse
// @Failure 404 {object} accessErrorResponse
// @Failure 500 {object} accessErrorResponse
// @Router /v1/access/request [post]
func (h *AccessHandler) handleRequestAccess(w http.ResponseWriter, r *http.Request) {
	var req requestAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	// Validate required fields
	if req.SandboxID == "" {
		writeError(w, http.StatusBadRequest, "sandbox_id is required", "")
		return
	}
	if req.UserID == "" {
		writeError(w, http.StatusBadRequest, "user_id is required", "")
		return
	}
	if req.PublicKey == "" {
		writeError(w, http.StatusBadRequest, "public_key is required", "")
		return
	}

	// Get source IP from request
	sourceIP := r.RemoteAddr
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		sourceIP = xff
	}

	// Build access request
	accessReq := &sshca.AccessRequest{
		SandboxID:   req.SandboxID,
		UserID:      req.UserID,
		PublicKey:   req.PublicKey,
		TTLMinutes:  req.TTLMinutes,
		SourceIP:    sourceIP,
		RequestTime: time.Now(),
	}

	// Request access
	resp, err := h.accessSvc.RequestAccess(r.Context(), accessReq)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to issue certificate", err.Error())
		return
	}

	// Build response with instructions
	instructions := `To connect to the sandbox:
1. Save your private key to a file (e.g., sandbox_key)
2. Save the certificate to sandbox_key-cert.pub
3. Run: chmod 600 sandbox_key
4. Connect using the command provided in connect_command`

	_ = serverJSON.RespondJSON(w, http.StatusOK, requestAccessResponse{
		CertificateID:  resp.CertificateID,
		Certificate:    resp.Certificate,
		VMIPAddress:    resp.VMIPAddress,
		SSHPort:        resp.SSHPort,
		Username:       resp.Username,
		ValidUntil:     resp.ValidUntil.Format(time.RFC3339),
		TTLSeconds:     resp.TTLSeconds,
		ConnectCommand: resp.ConnectCommand,
		Instructions:   instructions,
	})
}

// handleGetCAPublicKey handles GET /v1/access/ca-pubkey
// @Summary Get the SSH CA public key
// @Description Returns the CA public key that should be trusted by VMs
// @Tags Access
// @Produce json
// @Success 200 {object} caPublicKeyResponse
// @Failure 500 {object} accessErrorResponse
// @Router /v1/access/ca-pubkey [get]
func (h *AccessHandler) handleGetCAPublicKey(w http.ResponseWriter, r *http.Request) {
	pubKey, err := h.accessSvc.GetCAPublicKey()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get CA public key", err.Error())
		return
	}

	usage := `This CA public key should be installed in VM images at /etc/ssh/ssh_ca.pub
and referenced in sshd_config with: TrustedUserCAKeys /etc/ssh/ssh_ca.pub`

	_ = serverJSON.RespondJSON(w, http.StatusOK, caPublicKeyResponse{
		PublicKey: pubKey,
		Usage:     usage,
	})
}

// handleGetCertificate handles GET /v1/access/certificate/{certID}
// @Summary Get certificate details
// @Description Returns details about an issued certificate
// @Tags Access
// @Produce json
// @Param certID path string true "Certificate ID"
// @Success 200 {object} certificateResponse
// @Failure 404 {object} accessErrorResponse
// @Failure 500 {object} accessErrorResponse
// @Router /v1/access/certificate/{certID} [get]
func (h *AccessHandler) handleGetCertificate(w http.ResponseWriter, r *http.Request) {
	certID := chi.URLParam(r, "certID")
	if certID == "" {
		writeError(w, http.StatusBadRequest, "certificate ID is required", "")
		return
	}

	cert, err := h.accessSvc.GetCertificate(r.Context(), certID)
	if err != nil {
		writeError(w, http.StatusNotFound, "certificate not found", err.Error())
		return
	}

	ttlSeconds := 0
	if !cert.IsExpired() {
		ttlSeconds = int(cert.TimeToExpiry().Seconds())
	}

	_ = serverJSON.RespondJSON(w, http.StatusOK, certificateResponse{
		ID:           cert.ID,
		SandboxID:    cert.SandboxID,
		UserID:       cert.UserID,
		VMID:         cert.VMID,
		Identity:     cert.Identity,
		SerialNumber: cert.SerialNumber,
		Principals:   cert.Principals,
		ValidAfter:   cert.ValidAfter.Format(time.RFC3339),
		ValidBefore:  cert.ValidBefore.Format(time.RFC3339),
		Status:       string(cert.Status),
		IssuedAt:     cert.IssuedAt.Format(time.RFC3339),
		IsExpired:    cert.IsExpired(),
		TTLSeconds:   ttlSeconds,
	})
}

// handleRevokeCertificate handles DELETE /v1/access/certificate/{certID}
// @Summary Revoke a certificate
// @Description Immediately revokes a certificate, terminating any active sessions
// @Tags Access
// @Accept json
// @Produce json
// @Param certID path string true "Certificate ID"
// @Param request body revokeCertificateRequest false "Revocation reason"
// @Success 200 {object} map[string]string
// @Failure 400 {object} accessErrorResponse
// @Failure 404 {object} accessErrorResponse
// @Failure 500 {object} accessErrorResponse
// @Router /v1/access/certificate/{certID} [delete]
func (h *AccessHandler) handleRevokeCertificate(w http.ResponseWriter, r *http.Request) {
	certID := chi.URLParam(r, "certID")
	if certID == "" {
		writeError(w, http.StatusBadRequest, "certificate ID is required", "")
		return
	}

	var req revokeCertificateRequest
	if r.Body != nil && r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// Non-fatal, use empty reason
			req.Reason = "revoked via API"
		}
	}
	if req.Reason == "" {
		req.Reason = "revoked via API"
	}

	if err := h.accessSvc.RevokeAccess(r.Context(), certID, req.Reason); err != nil {
		if err == sshca.ErrCertAlreadyRevoked {
			writeError(w, http.StatusBadRequest, "certificate already revoked", "")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to revoke certificate", err.Error())
		return
	}

	_ = serverJSON.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "certificate revoked successfully",
		"id":      certID,
	})
}

// handleListCertificates handles GET /v1/access/certificates
// @Summary List certificates
// @Description Lists issued certificates with optional filtering
// @Tags Access
// @Produce json
// @Param sandbox_id query string false "Filter by sandbox ID"
// @Param user_id query string false "Filter by user ID"
// @Param status query string false "Filter by status (ACTIVE, EXPIRED, REVOKED)"
// @Param active_only query bool false "Only show active, non-expired certificates"
// @Param limit query int false "Maximum results to return"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} listCertificatesResponse
// @Failure 500 {object} accessErrorResponse
// @Router /v1/access/certificates [get]
func (h *AccessHandler) handleListCertificates(w http.ResponseWriter, r *http.Request) {
	// Build filter
	filter := sshca.CertificateFilter{}

	if sandboxID := r.URL.Query().Get("sandbox_id"); sandboxID != "" {
		filter.SandboxID = &sandboxID
	}
	if userID := r.URL.Query().Get("user_id"); userID != "" {
		filter.UserID = &userID
	}
	if status := r.URL.Query().Get("status"); status != "" {
		s := sshca.CertStatus(status)
		filter.Status = &s
	}
	if activeOnly := r.URL.Query().Get("active_only"); activeOnly == "true" {
		filter.ActiveOnly = true
	}

	// Build options
	opts := &sshca.ListOptions{}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			opts.Limit = l
		}
	}
	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			opts.Offset = o
		}
	}

	certs, err := h.accessSvc.ListCertificates(r.Context(), filter, opts)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list certificates", err.Error())
		return
	}

	// Convert to response format
	responses := make([]certificateResponse, len(certs))
	for i, cert := range certs {
		ttlSeconds := 0
		if !cert.IsExpired() {
			ttlSeconds = int(cert.TimeToExpiry().Seconds())
		}
		responses[i] = certificateResponse{
			ID:           cert.ID,
			SandboxID:    cert.SandboxID,
			UserID:       cert.UserID,
			VMID:         cert.VMID,
			Identity:     cert.Identity,
			SerialNumber: cert.SerialNumber,
			Principals:   cert.Principals,
			ValidAfter:   cert.ValidAfter.Format(time.RFC3339),
			ValidBefore:  cert.ValidBefore.Format(time.RFC3339),
			Status:       string(cert.Status),
			IssuedAt:     cert.IssuedAt.Format(time.RFC3339),
			IsExpired:    cert.IsExpired(),
			TTLSeconds:   ttlSeconds,
		}
	}

	_ = serverJSON.RespondJSON(w, http.StatusOK, listCertificatesResponse{
		Certificates: responses,
		Total:        len(responses),
	})
}

// handleRecordSessionStart handles POST /v1/access/session/start
// @Summary Record session start
// @Description Records the start of an SSH session (called by VM or auth service)
// @Tags Access
// @Accept json
// @Produce json
// @Param request body sessionStartRequest true "Session start request"
// @Success 200 {object} sessionStartResponse
// @Failure 400 {object} accessErrorResponse
// @Failure 500 {object} accessErrorResponse
// @Router /v1/access/session/start [post]
func (h *AccessHandler) handleRecordSessionStart(w http.ResponseWriter, r *http.Request) {
	var req sessionStartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if req.CertificateID == "" {
		writeError(w, http.StatusBadRequest, "certificate_id is required", "")
		return
	}

	sourceIP := req.SourceIP
	if sourceIP == "" {
		sourceIP = r.RemoteAddr
	}

	sessionID, err := h.accessSvc.RecordSessionStart(r.Context(), req.CertificateID, sourceIP)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to record session start", err.Error())
		return
	}

	_ = serverJSON.RespondJSON(w, http.StatusOK, sessionStartResponse{
		SessionID: sessionID,
	})
}

// handleRecordSessionEnd handles POST /v1/access/session/end
// @Summary Record session end
// @Description Records the end of an SSH session
// @Tags Access
// @Accept json
// @Produce json
// @Param request body sessionEndRequest true "Session end request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} accessErrorResponse
// @Failure 500 {object} accessErrorResponse
// @Router /v1/access/session/end [post]
func (h *AccessHandler) handleRecordSessionEnd(w http.ResponseWriter, r *http.Request) {
	var req sessionEndRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if req.SessionID == "" {
		writeError(w, http.StatusBadRequest, "session_id is required", "")
		return
	}

	reason := req.Reason
	if reason == "" {
		reason = "session ended normally"
	}

	if err := h.accessSvc.RecordSessionEnd(r.Context(), req.SessionID, reason); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to record session end", err.Error())
		return
	}

	_ = serverJSON.RespondJSON(w, http.StatusOK, map[string]string{
		"message":    "session ended successfully",
		"session_id": req.SessionID,
	})
}

// handleListSessions handles GET /v1/access/sessions
// @Summary List sessions
// @Description Lists access sessions with optional filtering
// @Tags Access
// @Produce json
// @Param sandbox_id query string false "Filter by sandbox ID"
// @Param certificate_id query string false "Filter by certificate ID"
// @Param user_id query string false "Filter by user ID"
// @Param active_only query bool false "Only show active sessions"
// @Param limit query int false "Maximum results to return"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} listSessionsResponse
// @Failure 500 {object} accessErrorResponse
// @Router /v1/access/sessions [get]
func (h *AccessHandler) handleListSessions(w http.ResponseWriter, r *http.Request) {
	// Build filter
	filter := sshca.SessionFilter{}

	if sandboxID := r.URL.Query().Get("sandbox_id"); sandboxID != "" {
		filter.SandboxID = &sandboxID
	}
	if certID := r.URL.Query().Get("certificate_id"); certID != "" {
		filter.CertificateID = &certID
	}
	if userID := r.URL.Query().Get("user_id"); userID != "" {
		filter.UserID = &userID
	}
	if activeOnly := r.URL.Query().Get("active_only"); activeOnly == "true" {
		filter.ActiveOnly = true
	}

	// Build options
	opts := &sshca.ListOptions{}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			opts.Limit = l
		}
	}
	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			opts.Offset = o
		}
	}

	var sessions []*sshca.AccessSession
	var err error

	if filter.SandboxID != nil {
		sessions, err = h.accessSvc.GetActiveSessionsForSandbox(r.Context(), *filter.SandboxID)
	} else {
		// Return empty list if no sandbox_id filter (would need a more general method)
		sessions = []*sshca.AccessSession{}
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list sessions", err.Error())
		return
	}

	// Convert to response format
	responses := make([]sessionResponse, len(sessions))
	for i, session := range sessions {
		resp := sessionResponse{
			ID:            session.ID,
			CertificateID: session.CertificateID,
			SandboxID:     session.SandboxID,
			UserID:        session.UserID,
			VMID:          session.VMID,
			VMIPAddress:   session.VMIPAddress,
			SourceIP:      session.SourceIP,
			Status:        string(session.Status),
			StartedAt:     session.StartedAt.Format(time.RFC3339),
		}
		if session.EndedAt != nil {
			resp.EndedAt = session.EndedAt.Format(time.RFC3339)
		}
		if session.DurationSeconds != nil {
			resp.DurationSeconds = *session.DurationSeconds
		} else if session.EndedAt != nil {
			resp.DurationSeconds = int(session.EndedAt.Sub(session.StartedAt).Seconds())
		}
		responses[i] = resp
	}

	_ = serverJSON.RespondJSON(w, http.StatusOK, listSessionsResponse{
		Sessions: responses,
		Total:    len(responses),
	})
}

// Ensure serverError is used to avoid unused import error
var _ = serverError.ErrorResponse{}
