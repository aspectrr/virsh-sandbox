package ansible

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"sync"

	"github.com/google/uuid"
)

// JobStatus represents the current state of an Ansible job.
type JobStatus string

const (
	JobStatusPending  JobStatus = "pending"
	JobStatusRunning  JobStatus = "running"
	JobStatusFinished JobStatus = "finished"
	JobStatusFailed   JobStatus = "failed"
)

// Job represents an Ansible playbook execution job.
type Job struct {
	ID       string    `json:"id"`
	VMName   string    `json:"vm_name"`
	Playbook string    `json:"playbook"`
	Check    bool      `json:"check"`
	Status   JobStatus `json:"status"`
}

// JobRequest contains parameters for creating a new Ansible job.
type JobRequest struct {
	VMName   string `json:"vm_name"`
	Playbook string `json:"playbook"`
	Check    bool   `json:"check"`
}

// JobResponse is returned when a job is created.
type JobResponse struct {
	JobID string `json:"job_id"`
	WSURL string `json:"ws_url"`
}

// Runner manages Ansible job execution.
type Runner struct {
	mu               sync.RWMutex
	jobs             map[string]*Job
	allowedPlaybooks map[string]struct{}
	inventoryPath    string
	ansibleImage     string
}

// NewRunner creates a new Ansible runner.
func NewRunner(inventoryPath, ansibleImage string, allowedPlaybooks []string) *Runner {
	allowed := make(map[string]struct{}, len(allowedPlaybooks))
	for _, p := range allowedPlaybooks {
		allowed[p] = struct{}{}
	}

	return &Runner{
		jobs:             make(map[string]*Job),
		allowedPlaybooks: allowed,
		inventoryPath:    inventoryPath,
		ansibleImage:     ansibleImage,
	}
}

// CreateJob creates a new Ansible job and returns its ID.
func (r *Runner) CreateJob(req JobRequest) (*JobResponse, error) {
	if _, ok := r.allowedPlaybooks[req.Playbook]; !ok {
		return nil, fmt.Errorf("playbook not allowed: %s", req.Playbook)
	}

	jobID := uuid.New().String()
	job := &Job{
		ID:       jobID,
		VMName:   req.VMName,
		Playbook: req.Playbook,
		Check:    req.Check,
		Status:   JobStatusPending,
	}

	r.mu.Lock()
	r.jobs[jobID] = job
	r.mu.Unlock()

	return &JobResponse{
		JobID: jobID,
		WSURL: fmt.Sprintf("/ws/ansible/jobs/%s", jobID),
	}, nil
}

// GetJob retrieves a job by ID.
func (r *Runner) GetJob(jobID string) (*Job, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	job, ok := r.jobs[jobID]
	return job, ok
}

// SetJobStatus updates a job's status.
func (r *Runner) SetJobStatus(jobID string, status JobStatus) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if job, ok := r.jobs[jobID]; ok {
		job.Status = status
	}
}

// OutputWriter is an interface for writing job output lines.
type OutputWriter interface {
	WriteLine(line string) error
}

// RunJob executes an Ansible job via Docker and streams output to the writer.
// This is a blocking call that returns when the job completes.
func (r *Runner) RunJob(ctx context.Context, jobID string, writer OutputWriter) error {
	job, ok := r.GetJob(jobID)
	if !ok {
		return fmt.Errorf("job not found: %s", jobID)
	}

	// Build ansible command
	ansibleCmd := fmt.Sprintf(
		"ansible-playbook -i %s playbooks/%s --limit %s",
		r.inventoryPath,
		job.Playbook,
		job.VMName,
	)
	if job.Check {
		ansibleCmd += " --check"
	}

	if err := writer.WriteLine(fmt.Sprintf("Running: %s\n", ansibleCmd)); err != nil {
		return err
	}

	// Build docker command
	dockerArgs := []string{
		"run",
		"--rm",
		"--network", "host",
		"--read-only",
		"--pids-limit", "128",
		"--memory", "512m",
		"-e", fmt.Sprintf("ANSIBLE_CMD=%s", ansibleCmd),
		"-v", "/ansible:/runner:ro",
		"-v", "/var/run/libvirt:/var/run/libvirt",
		r.ansibleImage,
	}

	cmd := exec.CommandContext(ctx, "docker", dockerArgs...)

	// Capture stdout and stderr together
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	cmd.Stderr = cmd.Stdout // Merge stderr into stdout

	r.SetJobStatus(jobID, JobStatusRunning)

	if err := cmd.Start(); err != nil {
		r.SetJobStatus(jobID, JobStatusFailed)
		return fmt.Errorf("failed to start docker: %w", err)
	}

	// Stream output line by line
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if err := writer.WriteLine(scanner.Text()); err != nil {
			// Client disconnected, but let the process continue
			break
		}
	}

	// Wait for the process to complete
	exitErr := cmd.Wait()
	exitCode := 0
	if exitErr != nil {
		if exitError, ok := exitErr.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			r.SetJobStatus(jobID, JobStatusFailed)
			return fmt.Errorf("failed to wait for docker: %w", exitErr)
		}
	}

	r.SetJobStatus(jobID, JobStatusFinished)

	_ = writer.WriteLine(fmt.Sprintf("\nJob finished (rc=%d)", exitCode))

	return nil
}
