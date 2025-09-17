package database

import (
	"context"
	"embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Migration represents a database migration
type Migration struct {
	Version     int
	Name        string
	SQL         string
	AppliedAt   *time.Time
	Checksum    string
}

// Migrator handles database migrations
type Migrator struct {
	db     *DB
	logger *zap.Logger
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *DB, logger *zap.Logger) *Migrator {
	if logger == nil {
		logger = zap.NewNop()
	}
	
	return &Migrator{
		db:     db,
		logger: logger,
	}
}

// createMigrationsTable creates the migrations tracking table
func (m *Migrator) createMigrationsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			checksum TEXT NOT NULL,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`
	
	_, err := m.db.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	
	return nil
}

// getAppliedMigrations returns a map of applied migrations
func (m *Migrator) getAppliedMigrations(ctx context.Context) (map[int]*Migration, error) {
	query := `SELECT version, name, checksum, applied_at FROM schema_migrations ORDER BY version`
	
	rows, err := m.db.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()
	
	applied := make(map[int]*Migration)
	for rows.Next() {
		mig := &Migration{}
		var appliedAt time.Time
		
		err := rows.Scan(&mig.Version, &mig.Name, &mig.Checksum, &appliedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan migration row: %w", err)
		}
		
		mig.AppliedAt = &appliedAt
		applied[mig.Version] = mig
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating migration rows: %w", err)
	}
	
	return applied, nil
}

// loadMigrations loads all migration files from the embedded filesystem
func (m *Migrator) loadMigrations() ([]*Migration, error) {
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}
	
	var migrations []*Migration
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		
		// Parse version from filename (e.g., "001_initial_schema.sql")
		parts := strings.SplitN(entry.Name(), "_", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid migration filename format: %s", entry.Name())
		}
		
		version, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid migration version in filename %s: %w", entry.Name(), err)
		}
		
		name := strings.TrimSuffix(parts[1], ".sql")
		
		// Read migration content
		content, err := migrationFiles.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", entry.Name(), err)
		}
		
		migrations = append(migrations, &Migration{
			Version: version,
			Name:    name,
			SQL:     string(content),
			Checksum: calculateChecksum(content),
		})
	}
	
	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
	
	return migrations, nil
}

// calculateChecksum calculates a simple checksum for migration content
func calculateChecksum(content []byte) string {
	// Simple hash - in production, consider using a proper hash function
	hash := 0
	for _, b := range content {
		hash = hash*31 + int(b)
	}
	return fmt.Sprintf("%x", hash)
}

// applyMigration applies a single migration
func (m *Migrator) applyMigration(ctx context.Context, migration *Migration) error {
	m.logger.Info("Applying migration",
		zap.Int("version", migration.Version),
		zap.String("name", migration.Name),
	)
	
	// Start transaction
	tx, err := m.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			m.logger.Error("Failed to rollback migration transaction", zap.Error(err))
		}
	}()
	
	// Execute migration SQL
	_, err = tx.Exec(ctx, migration.SQL)
	if err != nil {
		return fmt.Errorf("failed to execute migration %d: %w", migration.Version, err)
	}
	
	// Record migration as applied
	recordQuery := `
		INSERT INTO schema_migrations (version, name, checksum)
		VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(ctx, recordQuery, migration.Version, migration.Name, migration.Checksum)
	if err != nil {
		return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
	}
	
	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
	}
	
	m.logger.Info("Migration applied successfully",
		zap.Int("version", migration.Version),
		zap.String("name", migration.Name),
	)
	
	return nil
}

// ApplyMigrations applies all pending migrations
func (m *Migrator) ApplyMigrations(ctx context.Context) error {
	m.logger.Info("Starting database migration")
	
	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(ctx); err != nil {
		return err
	}
	
	// Load all migrations
	migrations, err := m.loadMigrations()
	if err != nil {
		return err
	}
	
	// Get applied migrations
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return err
	}
	
	// Apply pending migrations
	pendingCount := 0
	for _, migration := range migrations {
		appliedMig, exists := applied[migration.Version]
		
		if exists {
			// Verify checksum
			if appliedMig.Checksum != migration.Checksum {
				return fmt.Errorf("migration %d checksum mismatch: expected %s, got %s",
					migration.Version, appliedMig.Checksum, migration.Checksum)
			}
			continue
		}
		
		// Apply migration
		if err := m.applyMigration(ctx, migration); err != nil {
			return err
		}
		pendingCount++
	}
	
	m.logger.Info("Database migration completed",
		zap.Int("total_migrations", len(migrations)),
		zap.Int("applied_migrations", pendingCount),
	)
	
	return nil
}

// Status returns the current migration status
func (m *Migrator) Status(ctx context.Context) ([]*Migration, error) {
	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(ctx); err != nil {
		return nil, err
	}
	
	// Load all migrations
	migrations, err := m.loadMigrations()
	if err != nil {
		return nil, err
	}
	
	// Get applied migrations
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return nil, err
	}
	
	// Merge status
	for _, migration := range migrations {
		if appliedMig, exists := applied[migration.Version]; exists {
			migration.AppliedAt = appliedMig.AppliedAt
		}
	}
	
	return migrations, nil
}