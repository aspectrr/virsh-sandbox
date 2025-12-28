package ansible

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	serverError "virsh-sandbox/internal/error"
	serverJSON "virsh-sandbox/internal/json"
)

// Handler provides HTTP handlers for Ansible operations.
type Handler struct {
	runner   *Runner
	upgrader websocket.Upgrader
}

// NewHandler creates a new Ansible HTTP handler.
func NewHandler(runner *Runner) *Handler {
	return &Handler{
		runner: runner,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for now; tighten in production
				return true
			},
		},
	}
}

// wsOutputWriter implements OutputWriter for WebSocket connections.
type wsOutputWriter struct {
	conn *websocket.Conn
}

func (w *wsOutputWriter) WriteLine(line string) error {
	return w.conn.WriteMessage(websocket.TextMessage, []byte(line))
}

// HandleCreateJob creates a new Ansible job.
// @Summary Create Ansible job
// @Description Creates a new Ansible playbook execution job
// @Tags Ansible
// @Accept json
// @Produce json
// @Param request body JobRequest true "Job creation parameters"
// @Success 200 {object} JobResponse
// @Failure 400 {object} serverError.ErrorResponse
// @Router /v1/ansible/jobs [post]
func (h *Handler) HandleCreateJob(w http.ResponseWriter, r *http.Request) {
	var req JobRequest
	if err := serverJSON.DecodeJSON(r.Context(), r, &req); err != nil {
		serverError.RespondError(w, http.StatusBadRequest, err)
		return
	}

	if req.VMName == "" {
		serverError.RespondError(w, http.StatusBadRequest,
			&validationError{field: "vm_name", message: "vm_name is required"})
		return
	}
	if req.Playbook == "" {
		serverError.RespondError(w, http.StatusBadRequest,
			&validationError{field: "playbook", message: "playbook is required"})
		return
	}

	resp, err := h.runner.CreateJob(req)
	if err != nil {
		serverError.RespondError(w, http.StatusBadRequest, err)
		return
	}

	_ = serverJSON.RespondJSON(w, http.StatusOK, resp)
}

// HandleGetJob retrieves job status.
// @Summary Get Ansible job
// @Description Gets the status of an Ansible job
// @Tags Ansible
// @Accept json
// @Produce json
// @Param job_id path string true "Job ID"
// @Success 200 {object} Job
// @Failure 404 {object} serverError.ErrorResponse
// @Router /v1/ansible/jobs/{job_id} [get]
func (h *Handler) HandleGetJob(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "job_id")

	job, ok := h.runner.GetJob(jobID)
	if !ok {
		serverError.RespondError(w, http.StatusNotFound,
			&notFoundError{resource: "job", id: jobID})
		return
	}

	_ = serverJSON.RespondJSON(w, http.StatusOK, job)
}

// HandleJobWebSocket handles WebSocket connections for job output streaming.
// @Summary Stream Ansible job output
// @Description Connects via WebSocket to run an Ansible job and stream output
// @Tags Ansible
// @Param job_id path string true "Job ID"
// @Success 101 {string} string "Switching Protocols - WebSocket connection established"
// @Failure 404 {string} string "Invalid job ID"
// @Failure 409 {string} string "Job already started or finished"
// @Router /v1/ansible/jobs/{job_id}/stream [get]
func (h *Handler) HandleJobWebSocket(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "job_id")

	job, ok := h.runner.GetJob(jobID)
	if !ok {
		http.Error(w, "Invalid job ID", http.StatusNotFound)
		return
	}

	if job.Status != JobStatusPending {
		http.Error(w, "Job already started or finished", http.StatusConflict)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		// Upgrade already sends the error response
		return
	}
	defer conn.Close()

	// Set a reasonable deadline for the entire job
	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Minute)); err != nil {
		return
	}

	writer := &wsOutputWriter{conn: conn}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
	defer cancel()

	if err := h.runner.RunJob(ctx, jobID, writer); err != nil {
		_ = conn.WriteMessage(websocket.TextMessage, []byte("Error: "+err.Error()))
	}

	_ = conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

// RegisterRoutes registers Ansible routes on the given router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/ansible", func(r chi.Router) {
		r.Post("/jobs", h.HandleCreateJob)
		r.Get("/jobs/{job_id}", h.HandleGetJob)
		r.Get("/jobs/{job_id}/stream", h.HandleJobWebSocket)
	})
}

// validationError represents a validation error.
type validationError struct {
	field   string
	message string
}

func (e *validationError) Error() string {
	return e.message
}

// notFoundError represents a resource not found error.
type notFoundError struct {
	resource string
	id       string
}

func (e *notFoundError) Error() string {
	return e.resource + " not found: " + e.id
}
