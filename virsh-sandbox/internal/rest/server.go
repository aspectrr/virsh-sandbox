package rest

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"virsh-sandbox/internal/ansible"
	serverError "virsh-sandbox/internal/error"
	serverJSON "virsh-sandbox/internal/json"
	"virsh-sandbox/internal/libvirt"
	"virsh-sandbox/internal/store"
	"virsh-sandbox/internal/vm"
)

// Server wires the HTTP layer to application services.
type Server struct {
	Router         chi.Router
	vmSvc          *vm.Service
	domainMgr      *libvirt.DomainManager
	ansibleHandler *ansible.Handler
}

// NewServer constructs a REST server with routes registered.
func NewServer(vmSvc *vm.Service, domainMgr *libvirt.DomainManager, ansibleRunner *ansible.Runner) *Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	var ansibleHandler *ansible.Handler
	if ansibleRunner != nil {
		ansibleHandler = ansible.NewHandler(ansibleRunner)
	}

	s := &Server{
		Router:         router,
		vmSvc:          vmSvc,
		domainMgr:      domainMgr,
		ansibleHandler: ansibleHandler,
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

	// @Summary API reference
	// @Description Returns HTML API reference documentation
	// @Accept json
	// @Produce html
	// @Success 200 {string} string
	// @Router /docs [get]
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			// SpecURL: "https://generator3.swagger.io/openapi.json",// allow external URL or local path file
			SpecURL: "./docs/openapi.yaml",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Virsh Sandbox API",
			},
			DarkMode: true,
		})
		if err != nil {
			fmt.Printf("%v", err)
		}

		fmt.Fprintln(w, htmlContent)
	})

	// API v1 routes
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", s.handleHealth)
		r.Get("/vms", s.handleListVMs)

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

		// Ansible job management
		if s.ansibleHandler != nil {
			s.ansibleHandler.RegisterRoutes(r)
		}
	})
}

// --- Request/Response DTOs ---

type createSandboxRequest struct {
	SourceVMName string `json:"source_vm_name"`      // required; name of existing VM in libvirt to clone from
	AgentID      string `json:"agent_id"`            // required
	VMName       string `json:"vm_name,omitempty"`   // optional; generated if empty
	CPU          int    `json:"cpu,omitempty"`       // optional; default from service config if <=0
	MemoryMB     int    `json:"memory_mb,omitempty"` // optional; default from service config if <=0
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

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

type vmInfo struct {
	Name       string `json:"name"`
	UUID       string `json:"uuid"`
	State      string `json:"state"`
	Persistent bool   `json:"persistent"`
	DiskPath   string `json:"disk_path,omitempty"`
}

type listVMsResponse struct {
	VMs []vmInfo `json:"vms"`
}

// --- Handlers ---

// handleHealth returns service health status.
// @Summary Health check
// @Description Returns service health status
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Id getHealth
// @Router /v1/health [get]
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	_ = serverJSON.RespondJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

// @Summary Create a new sandbox
// @Description Creates a new virtual machine sandbox by cloning from an existing VM
// @Tags Sandbox
// @Accept json
// @Produce json
// @Param request body createSandboxRequest true "Sandbox creation parameters"
// @Success 201 {object} createSandboxResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Id createSandbox
// @Router /v1/sandbox/create [post]
func (s *Server) handleCreateSandbox(w http.ResponseWriter, r *http.Request) {
	var req createSandboxRequest
	if err := serverJSON.DecodeJSON(r.Context(), r, &req); err != nil {
		serverError.RespondError(w, http.StatusBadRequest, err)
		return
	}
	if req.SourceVMName == "" || req.AgentID == "" {
		serverError.RespondError(w, http.StatusBadRequest, errors.New("source_vm_name and agent_id are required"))
		return
	}

	sb, err := s.vmSvc.CreateSandbox(r.Context(), req.SourceVMName, req.AgentID, req.VMName, req.CPU, req.MemoryMB)
	if err != nil {
		serverError.RespondError(w, http.StatusInternalServerError, fmt.Errorf("create sandbox: %w", err))
		return
	}
	_ = serverJSON.RespondJSON(w, http.StatusCreated, createSandboxResponse{Sandbox: sb})
}

// @Summary Inject SSH key into sandbox
// @Description Injects a public SSH key for a user in the sandbox
// @Tags Sandbox
// @Accept json
// @Produce json
// @Param id path string true "Sandbox ID"
// @Param request body injectSSHKeyRequest true "SSH key injection parameters"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Id injectSshKey
// @Router /v1/sandbox/{id}/sshkey [post]
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

// @Summary Start sandbox
// @Description Starts the virtual machine sandbox
// @Tags Sandbox
// @Accept json
// @Produce json
// @Param id path string true "Sandbox ID"
// @Param request body startSandboxRequest false "Start parameters"
// @Success 200 {object} startSandboxResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Id startSandbox
// @Router /v1/sandbox/{id}/start [post]
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

// @Summary Run command in sandbox
// @Description Executes a command inside the sandbox via SSH
// @Tags Sandbox
// @Accept json
// @Produce json
// @Param id path string true "Sandbox ID"
// @Param request body runCommandRequest true "Command execution parameters"
// @Success 200 {object} runCommandResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Id runSandboxCommand
// @Router /v1/sandbox/{id}/run [post]
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

// @Summary Create snapshot
// @Description Creates a snapshot of the sandbox
// @Tags Sandbox
// @Accept json
// @Produce json
// @Param id path string true "Sandbox ID"
// @Param request body snapshotRequest true "Snapshot parameters"
// @Success 201 {object} snapshotResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Id createSnapshot
// @Router /v1/sandbox/{id}/snapshot [post]
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

// @Summary Diff snapshots
// @Description Computes differences between two snapshots
// @Tags Sandbox
// @Accept json
// @Produce json
// @Param id path string true "Sandbox ID"
// @Param request body diffRequest true "Diff parameters"
// @Success 200 {object} diffResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Id diffSnapshots
// @Router /v1/sandbox/{id}/diff [post]
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

// @Summary Generate configuration
// @Description Generates Ansible or Puppet configuration from sandbox changes
// @Tags Sandbox
// @Accept json
// @Produce json
// @Param id path string true "Sandbox ID"
// @Param tool path string true "Tool type (ansible or puppet)"
// @Success 501 {object} generateResponse
// @Failure 400 {object} ErrorResponse
// @Id generateConfiguration
// @Router /v1/sandbox/{id}/generate/{tool} [post]
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

// @Summary Publish changes
// @Description Publishes sandbox changes to GitOps repository
// @Tags Sandbox
// @Accept json
// @Produce json
// @Param id path string true "Sandbox ID"
// @Param request body publishRequest true "Publish parameters"
// @Success 501 {object} publishResponse
// @Failure 400 {object} ErrorResponse
// @Id publishChanges
// @Router /v1/sandbox/{id}/publish [post]
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

// @Summary List all VMs
// @Description Returns a list of all virtual machines from the libvirt instance
// @Tags VMs
// @Accept json
// @Produce json
// @Success 200 {object} listVMsResponse
// @Failure 500 {object} ErrorResponse
// @Id listVirtualMachines
// @Router /v1/vms [get]
func (s *Server) handleListVMs(w http.ResponseWriter, r *http.Request) {
	domains, err := s.domainMgr.ListDomains(r.Context())
	if err != nil {
		serverError.RespondError(w, http.StatusInternalServerError, fmt.Errorf("list vms: %w", err))
		return
	}

	vms := make([]vmInfo, 0, len(domains))
	for _, d := range domains {
		vms = append(vms, vmInfo{
			Name:       d.Name,
			UUID:       d.UUID,
			State:      d.State.String(),
			Persistent: d.Persistent,
			DiskPath:   d.DiskPath,
		})
	}

	_ = serverJSON.RespondJSON(w, http.StatusOK, listVMsResponse{VMs: vms})
}

// @Summary Destroy sandbox
// @Description Destroys the sandbox and cleans up resources
// @Tags Sandbox
// @Accept json
// @Produce json
// @Param id path string true "Sandbox ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Id destroySandbox
// @Router /v1/sandbox/{id} [delete]
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
