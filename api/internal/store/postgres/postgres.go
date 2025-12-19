package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"virsh-sandbox/internal/store"
)

// Ensure interface compliance.
var (
	_ store.Store     = (*postgresStore)(nil)
	_ store.DataStore = (*postgresStore)(nil)
)

type postgresStore struct {
	db   *gorm.DB
	conf store.Config
}

// New creates a Store backed by Postgres + GORM.
func New(ctx context.Context, cfg store.Config) (store.Store, error) {
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("postgres: missing DatabaseURL")
	}

	db, err := gorm.Open(
		postgres.Open(cfg.DatabaseURL),
		&gorm.Config{
			NowFunc: func() time.Time { return time.Now().UTC() },
			Logger:  logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("postgres: open: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("postgres: sql.DB handle: %w", err)
	}

	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	pg := &postgresStore{
		db:   db.WithContext(ctx),
		conf: cfg,
	}

	if cfg.AutoMigrate && !cfg.ReadOnly {
		if err := pg.autoMigrate(ctx); err != nil {
			_ = sqlDB.Close()
			return nil, err
		}
	}

	if err := pg.Ping(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	return pg, nil
}

// NewWithDB wraps an existing *gorm.DB (useful for tests).
func NewWithDB(db *gorm.DB, cfg store.Config) store.Store {
	return &postgresStore{db: db, conf: cfg}
}

func (s *postgresStore) Config() store.Config {
	return s.conf
}

func (s *postgresStore) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s *postgresStore) Ping(ctx context.Context) error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func (s *postgresStore) WithTx(ctx context.Context, fn func(tx store.DataStore) error) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&postgresStore{db: tx, conf: s.conf})
	})
}

// --- Sandbox ---

func (s *postgresStore) CreateSandbox(ctx context.Context, sb *store.Sandbox) error {
	if s.conf.ReadOnly {
		return fmt.Errorf("postgres: CreateSandbox: %w", store.ErrInvalid)
	}
	if sb == nil || sb.ID == "" || sb.JobID == "" || sb.AgentID == "" || sb.VMName == "" ||
		sb.BaseImage == "" || sb.Network == "" || sb.State == "" {
		return fmt.Errorf("postgres: CreateSandbox: %w", store.ErrInvalid)
	}

	now := time.Now().UTC()
	sb.CreatedAt = now
	sb.UpdatedAt = now

	if err := s.db.WithContext(ctx).Create(sandboxToModel(sb)).Error; err != nil {
		return mapDBError(err)
	}
	return nil
}

func (s *postgresStore) GetSandbox(ctx context.Context, id string) (*store.Sandbox, error) {
	var model SandboxModel
	if err := s.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&model).Error; err != nil {
		return nil, mapDBError(err)
	}
	return sandboxFromModel(&model), nil
}

func (s *postgresStore) GetSandboxByVMName(ctx context.Context, vmName string) (*store.Sandbox, error) {
	var model SandboxModel
	if err := s.db.WithContext(ctx).
		Where("vm_name = ? AND deleted_at IS NULL", vmName).
		First(&model).Error; err != nil {
		return nil, mapDBError(err)
	}
	return sandboxFromModel(&model), nil
}

func (s *postgresStore) ListSandboxes(ctx context.Context, filter store.SandboxFilter, opt *store.ListOptions) ([]*store.Sandbox, error) {
	tx := s.db.WithContext(ctx).Model(&SandboxModel{}).Where("deleted_at IS NULL")
	if filter.AgentID != nil {
		tx = tx.Where("agent_id = ?", *filter.AgentID)
	}
	if filter.JobID != nil {
		tx = tx.Where("job_id = ?", *filter.JobID)
	}
	if filter.BaseImage != nil {
		tx = tx.Where("base_image = ?", *filter.BaseImage)
	}
	if filter.State != nil {
		tx = tx.Where("state = ?", string(*filter.State))
	}
	if filter.VMName != nil {
		tx = tx.Where("vm_name = ?", *filter.VMName)
	}

	tx = applyListOptions(tx, opt, map[string]string{
		"created_at": "created_at",
		"updated_at": "updated_at",
		"vm_name":    "vm_name",
	})

	var models []SandboxModel
	if err := tx.Find(&models).Error; err != nil {
		return nil, mapDBError(err)
	}

	out := make([]*store.Sandbox, 0, len(models))
	for i := range models {
		out = append(out, sandboxFromModel(&models[i]))
	}
	return out, nil
}

func (s *postgresStore) UpdateSandbox(ctx context.Context, sb *store.Sandbox) error {
	if s.conf.ReadOnly {
		return fmt.Errorf("postgres: UpdateSandbox: %w", store.ErrInvalid)
	}
	if sb == nil || sb.ID == "" {
		return fmt.Errorf("postgres: UpdateSandbox: %w", store.ErrInvalid)
	}
	sb.UpdatedAt = time.Now().UTC()
	model := sandboxToModel(sb)

	res := s.db.WithContext(ctx).
		Model(&SandboxModel{}).
		Where("id = ? AND deleted_at IS NULL", sb.ID).
		Updates(map[string]any{
			"job_id":      model.JobID,
			"agent_id":    model.AgentID,
			"vm_name":     model.VMName,
			"base_image":  model.BaseImage,
			"network":     model.Network,
			"ip":          model.IPAddress,
			"state":       model.State,
			"ttl_seconds": model.TTLSeconds,
			"updated_at":  model.UpdatedAt,
		})

	if err := mapDBError(res.Error); err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return store.ErrNotFound
	}
	return nil
}

func (s *postgresStore) UpdateSandboxState(ctx context.Context, id string, newState store.SandboxState, ipAddr *string) error {
	if s.conf.ReadOnly {
		return fmt.Errorf("postgres: UpdateSandboxState: %w", store.ErrInvalid)
	}
	if id == "" {
		return fmt.Errorf("postgres: UpdateSandboxState: %w", store.ErrInvalid)
	}

	res := s.db.WithContext(ctx).Model(&SandboxModel{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"state":      string(newState),
			"ip":         copyString(ipAddr),
			"updated_at": time.Now().UTC(),
		})
	if err := mapDBError(res.Error); err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return store.ErrNotFound
	}
	return nil
}

func (s *postgresStore) DeleteSandbox(ctx context.Context, id string) error {
	if s.conf.ReadOnly {
		return fmt.Errorf("postgres: DeleteSandbox: %w", store.ErrInvalid)
	}
	if id == "" {
		return fmt.Errorf("postgres: DeleteSandbox: %w", store.ErrInvalid)
	}
	now := time.Now().UTC()
	res := s.db.WithContext(ctx).Model(&SandboxModel{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"state":      string(store.SandboxStateDestroyed),
			"deleted_at": &now,
			"updated_at": now,
		})
	if err := mapDBError(res.Error); err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return store.ErrNotFound
	}
	return nil
}

// --- Snapshot ---

func (s *postgresStore) CreateSnapshot(ctx context.Context, sn *store.Snapshot) error {
	if s.conf.ReadOnly {
		return fmt.Errorf("postgres: CreateSnapshot: %w", store.ErrInvalid)
	}
	if sn == nil || sn.ID == "" || sn.SandboxID == "" || sn.Name == "" || sn.Ref == "" || sn.Kind == "" {
		return fmt.Errorf("postgres: CreateSnapshot: %w", store.ErrInvalid)
	}
	if sn.CreatedAt.IsZero() {
		sn.CreatedAt = time.Now().UTC()
	}
	if err := s.db.WithContext(ctx).Create(snapshotToModel(sn)).Error; err != nil {
		return mapDBError(err)
	}
	return nil
}

func (s *postgresStore) GetSnapshot(ctx context.Context, id string) (*store.Snapshot, error) {
	var model SnapshotModel
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, mapDBError(err)
	}
	return snapshotFromModel(&model), nil
}

func (s *postgresStore) GetSnapshotByName(ctx context.Context, sandboxID, name string) (*store.Snapshot, error) {
	var model SnapshotModel
	if err := s.db.WithContext(ctx).
		Where("sandbox_id = ? AND name = ?", sandboxID, name).
		First(&model).Error; err != nil {
		return nil, mapDBError(err)
	}
	return snapshotFromModel(&model), nil
}

func (s *postgresStore) ListSnapshots(ctx context.Context, sandboxID string, opt *store.ListOptions) ([]*store.Snapshot, error) {
	tx := s.db.WithContext(ctx).Model(&SnapshotModel{}).Where("sandbox_id = ?", sandboxID)
	tx = applyListOptions(tx, opt, map[string]string{
		"created_at": "created_at",
		"name":       "name",
	})

	var models []SnapshotModel
	if err := tx.Find(&models).Error; err != nil {
		return nil, mapDBError(err)
	}
	out := make([]*store.Snapshot, 0, len(models))
	for i := range models {
		out = append(out, snapshotFromModel(&models[i]))
	}
	return out, nil
}

// --- Command ---

func (s *postgresStore) SaveCommand(ctx context.Context, cmd *store.Command) error {
	if s.conf.ReadOnly {
		return fmt.Errorf("postgres: SaveCommand: %w", store.ErrInvalid)
	}
	if cmd == nil || cmd.ID == "" || cmd.SandboxID == "" || cmd.Command == "" {
		return fmt.Errorf("postgres: SaveCommand: %w", store.ErrInvalid)
	}
	if cmd.StartedAt.IsZero() {
		cmd.StartedAt = time.Now().UTC()
	}
	if cmd.EndedAt.IsZero() {
		cmd.EndedAt = time.Now().UTC()
	}

	if err := s.db.WithContext(ctx).Create(commandToModel(cmd)).Error; err != nil {
		return mapDBError(err)
	}
	return nil
}

func (s *postgresStore) GetCommand(ctx context.Context, id string) (*store.Command, error) {
	var model CommandModel
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, mapDBError(err)
	}
	return commandFromModel(&model), nil
}

func (s *postgresStore) ListCommands(ctx context.Context, sandboxID string, opt *store.ListOptions) ([]*store.Command, error) {
	tx := s.db.WithContext(ctx).Model(&CommandModel{}).Where("sandbox_id = ?", sandboxID)
	tx = applyListOptions(tx, opt, map[string]string{
		"started_at": "started_at",
		"ended_at":   "ended_at",
	})

	var models []CommandModel
	if err := tx.Find(&models).Error; err != nil {
		return nil, mapDBError(err)
	}
	out := make([]*store.Command, 0, len(models))
	for i := range models {
		out = append(out, commandFromModel(&models[i]))
	}
	return out, nil
}

// --- Diff ---

func (s *postgresStore) SaveDiff(ctx context.Context, d *store.Diff) error {
	if s.conf.ReadOnly {
		return fmt.Errorf("postgres: SaveDiff: %w", store.ErrInvalid)
	}
	if d == nil || d.ID == "" || d.SandboxID == "" || d.FromSnapshot == "" || d.ToSnapshot == "" {
		return fmt.Errorf("postgres: SaveDiff: %w", store.ErrInvalid)
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now().UTC()
	}
	model, err := diffToModel(d)
	if err != nil {
		return err
	}
	if err := s.db.WithContext(ctx).Create(model).Error; err != nil {
		return mapDBError(err)
	}
	return nil
}

func (s *postgresStore) GetDiff(ctx context.Context, id string) (*store.Diff, error) {
	var model DiffModel
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, mapDBError(err)
	}
	return diffFromModel(&model)
}

func (s *postgresStore) GetDiffBySnapshots(ctx context.Context, sandboxID, fromSnapshot, toSnapshot string) (*store.Diff, error) {
	var model DiffModel
	if err := s.db.WithContext(ctx).
		Where("sandbox_id = ? AND from_snapshot = ? AND to_snapshot = ?", sandboxID, fromSnapshot, toSnapshot).
		First(&model).Error; err != nil {
		return nil, mapDBError(err)
	}
	return diffFromModel(&model)
}

// --- ChangeSet ---

func (s *postgresStore) CreateChangeSet(ctx context.Context, cs *store.ChangeSet) error {
	if s.conf.ReadOnly {
		return fmt.Errorf("postgres: CreateChangeSet: %w", store.ErrInvalid)
	}
	if cs == nil || cs.ID == "" || cs.JobID == "" || cs.SandboxID == "" || cs.DiffID == "" ||
		cs.PathAnsible == "" || cs.PathPuppet == "" {
		return fmt.Errorf("postgres: CreateChangeSet: %w", store.ErrInvalid)
	}
	if cs.CreatedAt.IsZero() {
		cs.CreatedAt = time.Now().UTC()
	}
	if err := s.db.WithContext(ctx).Create(changeSetToModel(cs)).Error; err != nil {
		return mapDBError(err)
	}
	return nil
}

func (s *postgresStore) GetChangeSet(ctx context.Context, id string) (*store.ChangeSet, error) {
	var model ChangeSetModel
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, mapDBError(err)
	}
	return changeSetFromModel(&model), nil
}

func (s *postgresStore) GetChangeSetByJob(ctx context.Context, jobID string) (*store.ChangeSet, error) {
	var model ChangeSetModel
	if err := s.db.WithContext(ctx).Where("job_id = ?", jobID).First(&model).Error; err != nil {
		return nil, mapDBError(err)
	}
	return changeSetFromModel(&model), nil
}

// --- Publication ---

func (s *postgresStore) CreatePublication(ctx context.Context, p *store.Publication) error {
	if s.conf.ReadOnly {
		return fmt.Errorf("postgres: CreatePublication: %w", store.ErrInvalid)
	}
	if p == nil || p.ID == "" || p.JobID == "" || p.RepoURL == "" || p.Branch == "" || p.Status == "" {
		return fmt.Errorf("postgres: CreatePublication: %w", store.ErrInvalid)
	}
	now := time.Now().UTC()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = now
	}
	if err := s.db.WithContext(ctx).Create(publicationToModel(p)).Error; err != nil {
		return mapDBError(err)
	}
	return nil
}

func (s *postgresStore) UpdatePublicationStatus(ctx context.Context, id string, status store.PublicationStatus, commitSHA, prURL, errMsg *string) error {
	if s.conf.ReadOnly {
		return fmt.Errorf("postgres: UpdatePublicationStatus: %w", store.ErrInvalid)
	}
	now := time.Now().UTC()
	res := s.db.WithContext(ctx).Model(&PublicationModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     string(status),
			"commit_sha": copyString(commitSHA),
			"pr_url":     copyString(prURL),
			"error_msg":  copyString(errMsg),
			"updated_at": now,
		})
	if err := mapDBError(res.Error); err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return store.ErrNotFound
	}
	return nil
}

func (s *postgresStore) GetPublication(ctx context.Context, id string) (*store.Publication, error) {
	var model PublicationModel
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, mapDBError(err)
	}
	return publicationFromModel(&model), nil
}

// --- Migration ---

func (s *postgresStore) autoMigrate(ctx context.Context) error {
	return s.db.WithContext(ctx).AutoMigrate(
		&SandboxModel{},
		&SnapshotModel{},
		&CommandModel{},
		&DiffModel{},
		&ChangeSetModel{},
		&PublicationModel{},
	)
}

// --- Models & Converters ---

type SandboxModel struct {
	ID         string     `gorm:"primaryKey;column:id"`
	JobID      string     `gorm:"column:job_id;not null;index"`
	AgentID    string     `gorm:"column:agent_id;not null;index"`
	VMName     string     `gorm:"column:vm_name;not null;uniqueIndex"`
	BaseImage  string     `gorm:"column:base_image;not null;index"`
	Network    string     `gorm:"column:network;not null"`
	IPAddress  *string    `gorm:"column:ip"`
	State      string     `gorm:"column:state;not null;index"`
	TTLSeconds *int       `gorm:"column:ttl_seconds"`
	CreatedAt  time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;not null"`
	DeletedAt  *time.Time `gorm:"column:deleted_at;index"`
}

func (SandboxModel) TableName() string { return "sandboxes" }

type SnapshotModel struct {
	ID        string    `gorm:"primaryKey;column:id"`
	SandboxID string    `gorm:"column:sandbox_id;not null;index;index:idx_snapshots_sandbox_name,unique"`
	Name      string    `gorm:"column:name;not null;index:idx_snapshots_sandbox_name,unique"`
	Kind      string    `gorm:"column:kind;not null"`
	Ref       string    `gorm:"column:ref;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	MetaJSON  *string   `gorm:"column:meta_json;type:jsonb"`
}

func (SnapshotModel) TableName() string { return "snapshots" }

type CommandModel struct {
	ID        string    `gorm:"primaryKey;column:id"`
	SandboxID string    `gorm:"column:sandbox_id;not null;index"`
	Command   string    `gorm:"column:command;not null"`
	EnvJSON   *string   `gorm:"column:env_json;type:jsonb"`
	Stdout    string    `gorm:"column:stdout;not null"`
	Stderr    string    `gorm:"column:stderr;not null"`
	ExitCode  int       `gorm:"column:exit_code;not null"`
	StartedAt time.Time `gorm:"column:started_at;not null;index"`
	EndedAt   time.Time `gorm:"column:ended_at;not null"`
}

func (CommandModel) TableName() string { return "commands" }

type DiffModel struct {
	ID           string         `gorm:"primaryKey;column:id"`
	SandboxID    string         `gorm:"column:sandbox_id;not null;index;index:idx_diffs_sandbox_snapshots,unique"`
	FromSnapshot string         `gorm:"column:from_snapshot;not null;index:idx_diffs_sandbox_snapshots,unique"`
	ToSnapshot   string         `gorm:"column:to_snapshot;not null;index:idx_diffs_sandbox_snapshots,unique"`
	DiffJSON     datatypes.JSON `gorm:"column:diff_json;type:jsonb;not null"`
	CreatedAt    time.Time      `gorm:"column:created_at;not null"`
}

func (DiffModel) TableName() string { return "diffs" }

type ChangeSetModel struct {
	ID          string    `gorm:"primaryKey;column:id"`
	JobID       string    `gorm:"column:job_id;not null;uniqueIndex"`
	SandboxID   string    `gorm:"column:sandbox_id;not null;index"`
	DiffID      string    `gorm:"column:diff_id;not null;index"`
	PathAnsible string    `gorm:"column:path_ansible;not null"`
	PathPuppet  string    `gorm:"column:path_puppet;not null"`
	MetaJSON    *string   `gorm:"column:meta_json;type:jsonb"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
}

func (ChangeSetModel) TableName() string { return "changesets" }

type PublicationModel struct {
	ID        string    `gorm:"primaryKey;column:id"`
	JobID     string    `gorm:"column:job_id;not null;index"`
	RepoURL   string    `gorm:"column:repo_url;not null"`
	Branch    string    `gorm:"column:branch;not null"`
	CommitSHA *string   `gorm:"column:commit_sha"`
	PRURL     *string   `gorm:"column:pr_url"`
	Status    string    `gorm:"column:status;not null;index"`
	ErrorMsg  *string   `gorm:"column:error_msg"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (PublicationModel) TableName() string { return "publications" }

func sandboxToModel(sb *store.Sandbox) *SandboxModel {
	return &SandboxModel{
		ID:         sb.ID,
		JobID:      sb.JobID,
		AgentID:    sb.AgentID,
		VMName:     sb.VMName,
		BaseImage:  sb.BaseImage,
		Network:    sb.Network,
		IPAddress:  copyString(sb.IPAddress),
		State:      string(sb.State),
		TTLSeconds: copyInt(sb.TTLSeconds),
		CreatedAt:  sb.CreatedAt,
		UpdatedAt:  sb.UpdatedAt,
		DeletedAt:  copyTime(sb.DeletedAt),
	}
}

func sandboxFromModel(m *SandboxModel) *store.Sandbox {
	return &store.Sandbox{
		ID:         m.ID,
		JobID:      m.JobID,
		AgentID:    m.AgentID,
		VMName:     m.VMName,
		BaseImage:  m.BaseImage,
		Network:    m.Network,
		IPAddress:  copyString(m.IPAddress),
		State:      store.SandboxState(m.State),
		TTLSeconds: copyInt(m.TTLSeconds),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		DeletedAt:  copyTime(m.DeletedAt),
	}
}

func snapshotToModel(sn *store.Snapshot) *SnapshotModel {
	return &SnapshotModel{
		ID:        sn.ID,
		SandboxID: sn.SandboxID,
		Name:      sn.Name,
		Kind:      string(sn.Kind),
		Ref:       sn.Ref,
		CreatedAt: sn.CreatedAt,
		MetaJSON:  copyString(sn.MetaJSON),
	}
}

func snapshotFromModel(m *SnapshotModel) *store.Snapshot {
	return &store.Snapshot{
		ID:        m.ID,
		SandboxID: m.SandboxID,
		Name:      m.Name,
		Kind:      store.SnapshotKind(m.Kind),
		Ref:       m.Ref,
		CreatedAt: m.CreatedAt,
		MetaJSON:  copyString(m.MetaJSON),
	}
}

func commandToModel(cmd *store.Command) *CommandModel {
	return &CommandModel{
		ID:        cmd.ID,
		SandboxID: cmd.SandboxID,
		Command:   cmd.Command,
		EnvJSON:   copyString(cmd.EnvJSON),
		Stdout:    cmd.Stdout,
		Stderr:    cmd.Stderr,
		ExitCode:  cmd.ExitCode,
		StartedAt: cmd.StartedAt,
		EndedAt:   cmd.EndedAt,
	}
}

func commandFromModel(m *CommandModel) *store.Command {
	return &store.Command{
		ID:        m.ID,
		SandboxID: m.SandboxID,
		Command:   m.Command,
		EnvJSON:   copyString(m.EnvJSON),
		Stdout:    m.Stdout,
		Stderr:    m.Stderr,
		ExitCode:  m.ExitCode,
		StartedAt: m.StartedAt,
		EndedAt:   m.EndedAt,
	}
}

func diffToModel(d *store.Diff) (*DiffModel, error) {
	payload, err := json.Marshal(d.DiffJSON)
	if err != nil {
		return nil, fmt.Errorf("postgres: marshal diff_json: %w", err)
	}
	return &DiffModel{
		ID:           d.ID,
		SandboxID:    d.SandboxID,
		FromSnapshot: d.FromSnapshot,
		ToSnapshot:   d.ToSnapshot,
		DiffJSON:     datatypes.JSON(payload),
		CreatedAt:    d.CreatedAt,
	}, nil
}

func diffFromModel(m *DiffModel) (*store.Diff, error) {
	var diff store.Diff
	diff.ID = m.ID
	diff.SandboxID = m.SandboxID
	diff.FromSnapshot = m.FromSnapshot
	diff.ToSnapshot = m.ToSnapshot
	diff.CreatedAt = m.CreatedAt
	if err := json.Unmarshal([]byte(m.DiffJSON), &diff.DiffJSON); err != nil {
		return nil, fmt.Errorf("postgres: unmarshal diff_json: %w", err)
	}
	return &diff, nil
}

func changeSetToModel(cs *store.ChangeSet) *ChangeSetModel {
	return &ChangeSetModel{
		ID:          cs.ID,
		JobID:       cs.JobID,
		SandboxID:   cs.SandboxID,
		DiffID:      cs.DiffID,
		PathAnsible: cs.PathAnsible,
		PathPuppet:  cs.PathPuppet,
		MetaJSON:    copyString(cs.MetaJSON),
		CreatedAt:   cs.CreatedAt,
	}
}

func changeSetFromModel(m *ChangeSetModel) *store.ChangeSet {
	return &store.ChangeSet{
		ID:          m.ID,
		JobID:       m.JobID,
		SandboxID:   m.SandboxID,
		DiffID:      m.DiffID,
		PathAnsible: m.PathAnsible,
		PathPuppet:  m.PathPuppet,
		MetaJSON:    copyString(m.MetaJSON),
		CreatedAt:   m.CreatedAt,
	}
}

func publicationToModel(p *store.Publication) *PublicationModel {
	return &PublicationModel{
		ID:        p.ID,
		JobID:     p.JobID,
		RepoURL:   p.RepoURL,
		Branch:    p.Branch,
		CommitSHA: copyString(p.CommitSHA),
		PRURL:     copyString(p.PRURL),
		Status:    string(p.Status),
		ErrorMsg:  copyString(p.ErrorMsg),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func publicationFromModel(m *PublicationModel) *store.Publication {
	return &store.Publication{
		ID:        m.ID,
		JobID:     m.JobID,
		RepoURL:   m.RepoURL,
		Branch:    m.Branch,
		CommitSHA: copyString(m.CommitSHA),
		PRURL:     copyString(m.PRURL),
		Status:    store.PublicationStatus(m.Status),
		ErrorMsg:  copyString(m.ErrorMsg),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// --- Helpers ---

func applyListOptions(tx *gorm.DB, opt *store.ListOptions, whitelist map[string]string) *gorm.DB {
	orderApplied := false
	if opt != nil {
		if col, ok := whitelist[opt.OrderBy]; ok {
			dir := "DESC"
			if opt.Asc {
				dir = "ASC"
			}
			tx = tx.Order(fmt.Sprintf("%s %s", col, dir))
			orderApplied = true
		}
		if opt.Limit > 0 {
			tx = tx.Limit(opt.Limit)
			if opt.Offset > 0 {
				tx = tx.Offset(opt.Offset)
			}
		}
	}
	if !orderApplied {
		tx = tx.Order("created_at DESC")
	}
	return tx
}

func copyString(src *string) *string {
	if src == nil {
		return nil
	}
	val := *src
	return &val
}

func copyInt(src *int) *int {
	if src == nil {
		return nil
	}
	val := *src
	return &val
}

func copyTime(src *time.Time) *time.Time {
	if src == nil {
		return nil
	}
	val := *src
	return &val
}

func mapDBError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return store.ErrNotFound
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return store.ErrAlreadyExists
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return store.ErrAlreadyExists
		case "23503":
			return store.ErrInvalid
		}
	}
	return err
}
