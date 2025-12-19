package rest

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	serverError "virsh-sandbox/internal/error"
	serverJSON "virsh-sandbox/internal/json"
	"virsh-sandbox/internal/store"
	"virsh-sandbox/internal/vm"
)

// Server wires the HTTP layer to application services.
type Server struct {
	Router chi.Router
	vmSvc  *vm.Service
}

// NewServer constructs a REST server with routes registered.
func NewServer(vmSvc *vm.Service) *Server {
	s := &Server{
		Router: chi.NewRouter(),
		vmSvc:  vmSvc,
	}
	s.routes()
	return s
}

// StartHTTP runs the HTTP server on the given address.
func (s *Server) StartHTTP(addr string) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           s.Router,
		ReadHeaderTimeout: 10 * time.Second,
	}
	return srv.ListenAndServe()
}

func (s *Server) routes() {
	r := s.Router

	// Basic liveness endpoint
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_ = serverJSON.RespondJSON(w, http.StatusOK, map[string]any{"status": "ok"})
	})

	// Sandbox lifecycle
	r.Route("/sandbox", func(r chi.Router) {
		r.Post("/create", s.handleCreateSandbox)

		r.Route("/{id}", func(r chi.Router) {
			r.Post("/sshkey", s.handleInjectSSHKey)
			r.Post("/start", s.handleStartSandbox)
			r.Post("/run", s.handleRunCommand)
			r.Post("/snapshot", s.handleCreateSnapshot)
			r.Post("/diff", s.handleDiffSnapshots)

			r.Post("/generate/{tool}", s.handleGenerate) // tool ∈ {ansible, puppet}
			r.Post("/publish", s.handlePublish)

			r.Delete("/", s.handleDestroySandbox)
		})
	})
}

// --- Request/Response DTOs ---

type createSandboxRequest struct {
	BaseImage string `json:"base_image"`          // required; e.g. "ubuntu-22.04.qcow2"
	AgentID   string `json:"agent_id"`            // required
	VMName    string `json:"vm_name,omitempty"`   // optional; generated if empty
	CPU       int    `json:"cpu,omitempty"`       // optional; default from service config if <=0
	MemoryMB  int    `json:"memory_mb,omitempty"` // optional; default from service config if <=0
}

type createSandboxResponse struct {
	Sandbox *store.Sandbox `json:"sandbox"`
}

type injectSSHKeyRequest struct {
	PublicKey string `json:"public_key"`         // required
	Username  string `json:"username,omitempty"` // required (explicit); typical: "ubuntu" or "centos"
}

type startSandboxRequest struct {
	WaitForIP bool `json:"wait_for_ip"` // optional; default false
}

type startSandboxResponse struct {
	IPAddress string `json:"ip_address,omitempty"`
}

type runCommandRequest struct {
	Username       string            `json:"username"`              // required
	PrivateKeyPath string            `json:"private_key_path"`      // required; path on API host
	Command        string            `json:"command"`               // required
	TimeoutSec     int               `json:"timeout_sec,omitempty"` // optional; default from service config
	Env            map[string]string `json:"env,omitempty"`         // optional
}

type runCommandResponse struct {
	Command *store.Command `json:"command"`
}

type snapshotRequest struct {
	Name     string `json:"name"`               // required
	External bool   `json:"external,omitempty"` // optional; default false (internal snapshot)
}

type snapshotResponse struct {
	Snapshot *store.Snapshot `json:"snapshot"`
}

type diffRequest struct {
	FromSnapshot string `json:"from_snapshot"` // required
	ToSnapshot   string `json:"to_snapshot"`   // required
}

type diffResponse struct {
	Diff *store.Diff `json:"diff"`
}

type generateResponse struct {
	Message string `json:"message"`
	Note    string `json:"note,omitempty"`
}

type publishRequest struct {
	JobID     string   `json:"job_id"`              // required
	Message   string   `json:"message,omitempty"`   // optional commit/PR message
	Reviewers []string `json:"reviewers,omitempty"` // optional
}

type publishResponse struct {
	Message string `json:"message"`
	Note    string `json:"note,omitempty"`
}

// --- Handlers ---

func (s *Server) handleCreateSandbox(w http.ResponseWriter, r *http.Request) {
	var req createSandboxRequest
	if err := serverJSON.DecodeJSON(r.Context(), r, &req); err != nil {
		serverError.RespondError(w, http.StatusBadRequest, err)
		return
	}
	if req.BaseImage == "" || req.AgentID == "" {
		serverError.RespondError(w, http.StatusBadRequest, errors.New("base_image and agent_id are required"))
		return
	}

	sb, err := s.vmSvc.CreateSandbox(r.Context(), req.BaseImage, req.AgentID, req.VMName, req.CPU, req.MemoryMB)
	if err != nil {
		serverError.RespondError(w, http.StatusInternalServerError, fmt.Errorf("create sandbox: %w", err))
		return
	}
	_ = serverJSON.RespondJSON(w, http.StatusCreated, createSandboxResponse{Sandbox: sb})
}

func (s *Server) handleInjectSSHKey(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req injectSSHKeyRequest
	if err := serverJSON.DecodeJSON(r.Context(), r, &req); err != nil {
		serverError.RespondError(w, http.StatusBadRequest, err)
		return
	}
	if req.PublicKey == "" || req.Username == "" {
		serverError.RespondError(w, http.StatusBadRequest, errors.New("public_key and username are required"))
		return
	}

	if err := s.vmSvc.InjectSSHKey(r.Context(), id, req.Username, req.PublicKey); err != nil {
		serverError.RespondError(w, http.StatusInternalServerError, fmt.Errorf("inject ssh key: %w", err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleStartSandbox(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req startSandboxRequest
	// tolerate empty body; default WaitForIP=false
	if r.ContentLength > 0 {
		if err := serverJSON.DecodeJSON(r.Context(), r, &req); err != nil {
			serverError.RespondError(w, http.StatusBadRequest, err)
			return
		}
	}

	ip, err := s.vmSvc.StartSandbox(r.Context(), id, req.WaitForIP)
	if err != nil {
		serverError.RespondError(w, http.StatusInternalServerError, fmt.Errorf("start sandbox: %w", err))
		return
	}
	_ = serverJSON.RespondJSON(w, http.StatusOK, startSandboxResponse{IPAddress: ip})
}

func (s *Server) handleRunCommand(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req runCommandRequest
	if err := serverJSON.DecodeJSON(r.Context(), r, &req); err != nil {
		serverError.RespondError(w, http.StatusBadRequest, err)
		return
	}
	if req.Username == "" || req.PrivateKeyPath == "" || req.Command == "" {
		serverError.RespondError(w, http.StatusBadRequest, errors.New("username, private_key_path and command are required"))
		return
	}
	timeout := time.Duration(req.TimeoutSec) * time.Second
	cmd, err := s.vmSvc.RunCommand(r.Context(), id, req.Username, req.PrivateKeyPath, req.Command, timeout, req.Env)
	if err != nil {
		serverError.RespondError(w, http.StatusInternalServerError, fmt.Errorf("run command: %w", err))
		return
	}
	_ = serverJSON.RespondJSON(w, http.StatusOK, runCommandResponse{Command: cmd})
}

func (s *Server) handleCreateSnapshot(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req snapshotRequest
	if err := serverJSON.DecodeJSON(r.Context(), r, &req); err != nil {
		serverError.RespondError(w, http.StatusBadRequest, err)
		return
	}
	if req.Name == "" {
		serverError.RespondError(w, http.StatusBadRequest, errors.New("name is required"))
		return
	}
	snap, err := s.vmSvc.CreateSnapshot(r.Context(), id, req.Name, req.External)
	if err != nil {
		serverError.RespondError(w, http.StatusInternalServerError, fmt.Errorf("create snapshot: %w", err))
		return
	}
	_ = serverJSON.RespondJSON(w, http.StatusCreated, snapshotResponse{Snapshot: snap})
}

func (s *Server) handleDiffSnapshots(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req diffRequest
	if err := serverJSON.DecodeJSON(r.Context(), r, &req); err != nil {
		serverError.RespondError(w, http.StatusBadRequest, err)
		return
	}
	if req.FromSnapshot == "" || req.ToSnapshot == "" {
		serverError.RespondError(w, http.StatusBadRequest, errors.New("from_snapshot and to_snapshot are required"))
		return
	}
	d, err := s.vmSvc.DiffSnapshots(r.Context(), id, req.FromSnapshot, req.ToSnapshot)
	if err != nil {
		serverError.RespondError(w, http.StatusInternalServerError, fmt.Errorf("diff snapshots: %w", err))
		return
	}
	_ = serverJSON.RespondJSON(w, http.StatusOK, diffResponse{Diff: d})
}

func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tool := chi.URLParam(r, "tool")
	if id == "" {
		serverError.RespondError(w, http.StatusBadRequest, errors.New("sandbox id is required"))
		return
	}
	switch tool {
	case "ansible", "puppet":
		// Stub: these will be implemented when ansible/puppet generators are wired.
		_ = serverJSON.RespondJSON(w, http.StatusNotImplemented, generateResponse{
			Message: "generation not implemented yet",
			Note:    "tool=" + tool + " for sandbox " + id,
		})
	default:
		serverError.RespondError(w, http.StatusBadRequest, fmt.Errorf("unsupported tool %q; expected 'ansible' or 'puppet'", tool))
	}
}

func (s *Server) handlePublish(w http.ResponseWriter, r *http.Request) {
	var req publishRequest
	if err := serverJSON.DecodeJSON(r.Context(), r, &req); err != nil {
		serverError.RespondError(w, http.StatusBadRequest, err)
		return
	}
	if req.JobID == "" {
		serverError.RespondError(w, http.StatusBadRequest, errors.New("job_id is required"))
		return
	}
	// Stub: implement when GitOps publisher is wired.
	_ = serverJSON.RespondJSON(w, http.StatusNotImplemented, publishResponse{
		Message: "publish not implemented yet",
		Note:    "job_id=" + req.JobID,
	})
}

func (s *Server) handleDestroySandbox(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		serverError.RespondError(w, http.StatusBadRequest, errors.New("sandbox id is required"))
		return
	}
	if err := s.vmSvc.DestroySandbox(r.Context(), id); err != nil {
		serverError.RespondError(w, http.StatusInternalServerError, fmt.Errorf("destroy sandbox: %w", err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
