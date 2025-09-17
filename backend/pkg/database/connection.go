package database

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Config holds database connection configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string

	// Connection pool settings
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// DefaultConfig returns a default database configuration
func DefaultConfig() *Config {
	return &Config{
		Host:            "localhost",
		Port:            5432,
		User:            "verza",
		Password:        "verza",
		Database:        "verza",
		SSLMode:         "disable",
		MaxConns:        25,
		MinConns:        5,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: time.Minute * 30,
	}
}

// DSN returns the database connection string
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode,
	)
}

// DB wraps the database connection pool and queries
type DB struct {
	pool    *pgxpool.Pool
	queries *Queries
	logger  *zap.Logger
}

// NewDB creates a new database connection
func NewDB(ctx context.Context, config *Config, logger *zap.Logger) (*DB, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	// Configure connection pool
	poolConfig, err := pgxpool.ParseConfig(config.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	poolConfig.MaxConns = config.MaxConns
	poolConfig.MinConns = config.MinConns
	poolConfig.MaxConnLifetime = config.MaxConnLifetime
	poolConfig.MaxConnIdleTime = config.MaxConnIdleTime

	// Add query tracer for logging
	poolConfig.ConnConfig.Tracer = &queryTracer{logger: logger}

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		pool:    pool,
		queries: New(pool),
		logger:  logger,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() {
	db.pool.Close()
}

// Pool returns the underlying connection pool
func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

// Queries returns the generated queries
func (db *DB) Queries() *Queries {
	return db.queries
}

// Health checks database connectivity
func (db *DB) Health(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// Stats returns connection pool statistics
func (db *DB) Stats() *pgxpool.Stat {
	return db.pool.Stat()
}

// WithTx executes a function within a database transaction
func (db *DB) WithTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	queries := db.queries.WithTx(tx)
	if err := fn(queries); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// queryTracer implements pgx.QueryTracer for logging SQL queries
type queryTracer struct {
	logger *zap.Logger
}

type contextKey string

const queryStartTimeKey contextKey = "query_start_time"

func (t *queryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	t.logger.Debug("executing query",
		zap.String("sql", data.SQL),
		zap.Any("args", data.Args),
	)
	return context.WithValue(ctx, queryStartTimeKey, time.Now())
}

func (t *queryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	var duration time.Duration
	if startTime, ok := ctx.Value(queryStartTimeKey).(time.Time); ok {
		duration = time.Since(startTime)
	}

	if data.Err != nil {
		t.logger.Error("query failed",
			zap.Error(data.Err),
			zap.Duration("duration", duration),
		)
	} else {
		t.logger.Debug("query completed",
			zap.Duration("duration", duration),
		)
	}
}

// CalculateChecksum calculates SHA256 checksum of data
func CalculateChecksum(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}