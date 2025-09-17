package security

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestUserPermissions(t *testing.T) {
	tests := []struct {
		name       string
		roles      []Role
		permission Permission
		expected   bool
	}{
		{
			name:       "Admin has all permissions",
			roles:      []Role{RoleAdmin},
			permission: PermissionIssueVC,
			expected:   true,
		},
		{
			name:       "Issuer can issue VCs",
			roles:      []Role{RoleIssuer},
			permission: PermissionIssueVC,
			expected:   true,
		},
		{
			name:       "Verifier cannot issue VCs",
			roles:      []Role{RoleVerifier},
			permission: PermissionIssueVC,
			expected:   false,
		},
		{
			name:       "User can verify VCs",
			roles:      []Role{RoleUser},
			permission: PermissionVerifyVC,
			expected:   true,
		},
		{
			name:       "User cannot manage users",
			roles:      []Role{RoleUser},
			permission: PermissionManageUsers,
			expected:   false,
		},
		{
			name:       "Multiple roles combine permissions",
			roles:      []Role{RoleIssuer, RoleVerifier},
			permission: PermissionVerifyVC,
			expected:   true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				ID:       "test-user",
				Username: "testuser",
				Roles:    tt.roles,
				Active:   true,
			}
			
			result := user.HasPermission(tt.permission)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUserRoles(t *testing.T) {
	user := &User{
		ID:       "test-user",
		Username: "testuser",
		Roles:    []Role{RoleIssuer, RoleVerifier},
		Active:   true,
	}
	
	if !user.HasRole(RoleIssuer) {
		t.Error("Expected user to have issuer role")
	}
	
	if !user.HasRole(RoleVerifier) {
		t.Error("Expected user to have verifier role")
	}
	
	if user.HasRole(RoleAdmin) {
		t.Error("Expected user not to have admin role")
	}
}

func TestUserGetPermissions(t *testing.T) {
	user := &User{
		ID:       "test-user",
		Username: "testuser",
		Roles:    []Role{RoleIssuer},
		Active:   true,
	}
	
	permissions := user.GetPermissions()
	
	// Check that issuer permissions are included
	expectedPermissions := RolePermissions[RoleIssuer]
	for _, expectedPerm := range expectedPermissions {
		found := false
		for _, perm := range permissions {
			if perm == expectedPerm {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected permission %s not found", expectedPerm)
		}
	}
	
	// Check that admin-only permissions are not included
	if user.HasPermission(PermissionManageUsers) {
		t.Error("Issuer should not have admin permissions")
	}
}

func TestInMemoryUserStore(t *testing.T) {
	logger := zap.NewNop()
	store := NewInMemoryUserStore(logger)
	ctx := context.Background()
	
	user := &User{
		ID:        "test-user-1",
		Username:  "testuser1",
		Email:     "test@example.com",
		Roles:     []Role{RoleIssuer},
		APIKey:    "test-api-key-123",
		Active:    true,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	
	// Test CreateUser
	err := store.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Test duplicate user creation
	err = store.CreateUser(ctx, user)
	if err == nil {
		t.Error("Expected error when creating duplicate user")
	}
	
	// Test GetUserByID
	retrievedUser, err := store.GetUserByID(ctx, "test-user-1")
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}
	if retrievedUser.Username != "testuser1" {
		t.Errorf("Expected username 'testuser1', got '%s'", retrievedUser.Username)
	}
	
	// Test GetUserByAPIKey
	retrievedUser, err = store.GetUserByAPIKey(ctx, "test-api-key-123")
	if err != nil {
		t.Fatalf("Failed to get user by API key: %v", err)
	}
	if retrievedUser.ID != "test-user-1" {
		t.Errorf("Expected user ID 'test-user-1', got '%s'", retrievedUser.ID)
	}
	
	// Test invalid API key
	_, err = store.GetUserByAPIKey(ctx, "invalid-key")
	if err == nil {
		t.Error("Expected error for invalid API key")
	}
	
	// Test UpdateUser
	updatedUser := &User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     "updated@example.com",
		Roles:     user.Roles,
		APIKey:    "new-api-key-456",
		Active:    user.Active,
		CreatedAt: user.CreatedAt,
	}
	err = store.UpdateUser(ctx, updatedUser)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
	
	// Verify old API key is removed
	_, err = store.GetUserByAPIKey(ctx, "test-api-key-123")
	if err == nil {
		t.Error("Expected error for old API key after update")
	}
	
	// Verify new API key works
	retrievedUser, err = store.GetUserByAPIKey(ctx, "new-api-key-456")
	if err != nil {
		t.Fatalf("Failed to get user by new API key: %v", err)
	}
	if retrievedUser.Email != "updated@example.com" {
		t.Errorf("Expected email 'updated@example.com', got '%s'", retrievedUser.Email)
	}
	
	// Test inactive user
	updatedUser.Active = false
	store.UpdateUser(ctx, updatedUser)
	_, err = store.GetUserByAPIKey(ctx, "new-api-key-456")
	if err == nil {
		t.Error("Expected error for inactive user")
	}
	
	// Test ListUsers
	user2 := &User{
		ID:       "test-user-2",
		Username: "testuser2",
		Roles:    []Role{RoleVerifier},
		Active:   true,
	}
	store.CreateUser(ctx, user2)
	
	users, err := store.ListUsers(ctx, 10, 0)
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
	
	// Test pagination
	users, err = store.ListUsers(ctx, 1, 1)
	if err != nil {
		t.Fatalf("Failed to list users with pagination: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("Expected 1 user with pagination, got %d", len(users))
	}
	
	// Test DeleteUser
	err = store.DeleteUser(ctx, "test-user-1")
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}
	
	_, err = store.GetUserByID(ctx, "test-user-1")
	if err == nil {
		t.Error("Expected error when getting deleted user")
	}
}

func TestRBACMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	store := NewInMemoryUserStore(logger)
	ctx := context.Background()
	
	// Create test users
	adminUser := &User{
		ID:       "admin-1",
		Username: "admin",
		Roles:    []Role{RoleAdmin},
		APIKey:   "admin-key-123",
		Active:   true,
	}
	store.CreateUser(ctx, adminUser)
	
	issuerUser := &User{
		ID:       "issuer-1",
		Username: "issuer",
		Roles:    []Role{RoleIssuer},
		APIKey:   "issuer-key-456",
		Active:   true,
	}
	store.CreateUser(ctx, issuerUser)
	
	verifierUser := &User{
		ID:       "verifier-1",
		Username: "verifier",
		Roles:    []Role{RoleVerifier},
		APIKey:   "verifier-key-789",
		Active:   true,
	}
	store.CreateUser(ctx, verifierUser)
	
	rbac := NewRBACMiddleware(store, logger)
	
	// Test authentication middleware
	t.Run("Authentication", func(t *testing.T) {
		router := gin.New()
		router.Use(rbac.AuthenticateUser())
		router.GET("/test", func(c *gin.Context) {
			user, _ := GetCurrentUser(c)
			c.JSON(http.StatusOK, gin.H{"user_id": user.ID})
		})
		
		// Test missing API key
		req1 := httptest.NewRequest("GET", "/test", nil)
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		if w1.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 for missing API key, got %d", w1.Code)
		}
		
		// Test invalid API key
		req2 := httptest.NewRequest("GET", "/test", nil)
		req2.Header.Set("X-API-Key", "invalid-key")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		if w2.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 for invalid API key, got %d", w2.Code)
		}
		
		// Test valid API key
		req3 := httptest.NewRequest("GET", "/test", nil)
		req3.Header.Set("X-API-Key", "admin-key-123")
		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)
		if w3.Code != http.StatusOK {
			t.Errorf("Expected 200 for valid API key, got %d", w3.Code)
		}
		
		var response map[string]string
		json.Unmarshal(w3.Body.Bytes(), &response)
		if response["user_id"] != "admin-1" {
			t.Errorf("Expected user_id 'admin-1', got '%s'", response["user_id"])
		}
	})
	
	// Test permission middleware
	t.Run("Permission Check", func(t *testing.T) {
		router := gin.New()
		router.Use(rbac.AuthenticateUser())
		router.Use(rbac.RequirePermission(PermissionIssueVC))
		router.POST("/issue", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "VC issued"})
		})
		
		// Test admin (should have permission)
		req1 := httptest.NewRequest("POST", "/issue", nil)
		req1.Header.Set("X-API-Key", "admin-key-123")
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		if w1.Code != http.StatusOK {
			t.Errorf("Expected 200 for admin with permission, got %d", w1.Code)
		}
		
		// Test issuer (should have permission)
		req2 := httptest.NewRequest("POST", "/issue", nil)
		req2.Header.Set("X-API-Key", "issuer-key-456")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		if w2.Code != http.StatusOK {
			t.Errorf("Expected 200 for issuer with permission, got %d", w2.Code)
		}
		
		// Test verifier (should not have permission)
		req3 := httptest.NewRequest("POST", "/issue", nil)
		req3.Header.Set("X-API-Key", "verifier-key-789")
		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)
		if w3.Code != http.StatusForbidden {
			t.Errorf("Expected 403 for verifier without permission, got %d", w3.Code)
		}
	})
	
	// Test role middleware
	t.Run("Role Check", func(t *testing.T) {
		router := gin.New()
		router.Use(rbac.AuthenticateUser())
		router.Use(rbac.RequireRole(RoleAdmin))
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Admin area"})
		})
		
		// Test admin (should have role)
		req1 := httptest.NewRequest("GET", "/admin", nil)
		req1.Header.Set("X-API-Key", "admin-key-123")
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		if w1.Code != http.StatusOK {
			t.Errorf("Expected 200 for admin role, got %d", w1.Code)
		}
		
		// Test issuer (should not have admin role)
		req2 := httptest.NewRequest("GET", "/admin", nil)
		req2.Header.Set("X-API-Key", "issuer-key-456")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		if w2.Code != http.StatusForbidden {
			t.Errorf("Expected 403 for non-admin role, got %d", w2.Code)
		}
	})
}

func TestGetCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	user := &User{
		ID:       "test-user",
		Username: "testuser",
		Roles:    []Role{RoleIssuer},
	}
	
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		// Test without user in context
		_, err := GetCurrentUser(c)
		if err == nil {
			t.Error("Expected error when user not in context")
		}
		
		// Set user in context
		c.Set("user", user)
		
		// Test with user in context
		retrievedUser, err := GetCurrentUser(c)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if retrievedUser.ID != "test-user" {
			t.Errorf("Expected user ID 'test-user', got '%s'", retrievedUser.ID)
		}
		
		// Test with invalid user type in context
		c.Set("user", "invalid-user-type")
		_, err = GetCurrentUser(c)
		if err == nil {
			t.Error("Expected error for invalid user type")
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func BenchmarkUserHasPermission(b *testing.B) {
	user := &User{
		ID:       "test-user",
		Username: "testuser",
		Roles:    []Role{RoleIssuer, RoleVerifier},
		Active:   true,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user.HasPermission(PermissionIssueVC)
	}
}

func BenchmarkInMemoryUserStore(b *testing.B) {
	logger := zap.NewNop()
	store := NewInMemoryUserStore(logger)
	ctx := context.Background()
	
	// Create test user
	user := &User{
		ID:       "bench-user",
		Username: "benchuser",
		Roles:    []Role{RoleIssuer},
		APIKey:   "bench-api-key",
		Active:   true,
	}
	store.CreateUser(ctx, user)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.GetUserByAPIKey(ctx, "bench-api-key")
		}
	})
}