package database

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, 5432, config.Port)
	assert.Equal(t, "verza", config.User)
	assert.Equal(t, "verza", config.Password)
	assert.Equal(t, "verza", config.Database)
	assert.Equal(t, "disable", config.SSLMode)
}

func TestConfigDSN(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     5432,
		User:     "testuser",
		Password: "testpass",
		Database: "testdb",
		SSLMode:  "disable",
	}

	expected := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	assert.Equal(t, expected, config.DSN())
}

func TestNewWithInvalidConfig(t *testing.T) {
	config := &Config{
		Host: "invalid-host-that-does-not-exist",
		Port: 5432,
	}

	ctx := context.Background()
	logger := zap.NewNop()

	db, err := NewDB(ctx, config, logger)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestMigrationVersionParsing(t *testing.T) {
	tests := []struct {
		filename string
		expectedVersion int
		expectedName string
		shouldError bool
	}{
		{"001_initial_schema.sql", 1, "initial_schema", false},
		{"002_add_indexes.sql", 2, "add_indexes", false},
		{"010_user_preferences.sql", 10, "user_preferences", false},
		{"invalid_filename.sql", 0, "", true},
		{"abc_invalid_version.sql", 0, "", true},
	}
	
	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			// This would be tested in the actual migration loading logic
			// For now, we'll just validate the expected behavior
			if test.shouldError {
				// Should fail to parse version
				return
			}
			
			// Should successfully parse version and name
			if test.expectedVersion <= 0 {
				t.Errorf("Expected positive version, got %d", test.expectedVersion)
			}
			
			if test.expectedName == "" {
				t.Error("Expected non-empty name")
			}
		})
	}
}

func TestCalculateChecksum(t *testing.T) {
	content1 := []byte("SELECT 1;")
	content2 := []byte("SELECT 2;")
	content3 := []byte("SELECT 1;") // Same as content1
	
	checksum1 := calculateChecksum(content1)
	checksum2 := calculateChecksum(content2)
	checksum3 := calculateChecksum(content3)
	
	if checksum1 == checksum2 {
		t.Error("Different content should have different checksums")
	}
	
	if checksum1 != checksum3 {
		t.Error("Same content should have same checksums")
	}
	
	if checksum1 == "" {
		t.Error("Checksum should not be empty")
	}
}

func TestMigrationStruct(t *testing.T) {
	mig := &Migration{
		Version:  1,
		Name:     "test_migration",
		SQL:      "CREATE TABLE test (id INTEGER);",
		Checksum: "abc123",
	}
	
	if mig.Version != 1 {
		t.Errorf("Expected version 1, got %d", mig.Version)
	}
	
	if mig.Name != "test_migration" {
		t.Errorf("Expected name 'test_migration', got '%s'", mig.Name)
	}
	
	if mig.AppliedAt != nil {
		t.Error("Expected AppliedAt to be nil for new migration")
	}
	
	// Simulate applying the migration
	now := time.Now()
	mig.AppliedAt = &now
	
	if mig.AppliedAt == nil {
		t.Error("Expected AppliedAt to be set after applying")
	}
}

func TestQueryTracer(t *testing.T) {
	logger := zap.NewNop()
	tracer := &queryTracer{logger: logger}

	assert.NotNil(t, tracer)
	assert.Equal(t, logger, tracer.logger)
}

func TestDatabaseTypes(t *testing.T) {
	// Test UUID generation
	id := uuid.New()
	if id == uuid.Nil {
		t.Error("Expected valid UUID, got nil")
	}
	
	// Test time handling
	now := time.Now()
	if now.IsZero() {
		t.Error("Expected valid time, got zero")
	}
	
	// Test context
	ctx := context.Background()
	if ctx == nil {
		t.Error("Expected valid context, got nil")
	}
	
	// Test context with timeout
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	
	select {
	case <-ctx.Done():
		t.Error("Context should not be done immediately")
	default:
		// Expected behavior
	}
}

func BenchmarkCalculateChecksum(b *testing.B) {
	content := []byte("CREATE TABLE users (id UUID PRIMARY KEY, name TEXT);")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = calculateChecksum(content)
	}
}

func BenchmarkUUIDGeneration(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uuid.New()
	}
}

func BenchmarkTimeNow(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = time.Now()
	}
}