package security

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Role represents a user role
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleIssuer   Role = "issuer"
	RoleVerifier Role = "verifier"
	RoleUser     Role = "user"
)

// Permission represents a specific permission
type Permission string

const (
	// VC Permissions
	PermissionIssueVC     Permission = "vc:issue"
	PermissionVerifyVC    Permission = "vc:verify"
	PermissionRevokeVC    Permission = "vc:revoke"
	PermissionListVC      Permission = "vc:list"
	
	// DID Permissions
	PermissionCreateDID   Permission = "did:create"
	PermissionResolveDID  Permission = "did:resolve"
	PermissionUpdateDID   Permission = "did:update"
	PermissionDeactivateDID Permission = "did:deactivate"
	
	// KMS Permissions
	PermissionCreateKey   Permission = "kms:create_key"
	PermissionSignData    Permission = "kms:sign"
	PermissionRotateKey   Permission = "kms:rotate_key"
	PermissionDeleteKey   Permission = "kms:delete_key"
	
	// ML/KYC Permissions
	PermissionKYCVerify   Permission = "kyc:verify"
	PermissionKYCUpload   Permission = "kyc:upload"
	PermissionKYCCheck    Permission = "kyc:check"
	
	// Admin Permissions
	PermissionManageUsers Permission = "admin:manage_users"
	PermissionViewLogs    Permission = "admin:view_logs"
	PermissionManageKeys  Permission = "admin:manage_keys"
)

// RolePermissions maps roles to their permissions
var RolePermissions = map[Role][]Permission{
	RoleAdmin: {
		// Admin has all permissions
		PermissionIssueVC, PermissionVerifyVC, PermissionRevokeVC, PermissionListVC,
		PermissionCreateDID, PermissionResolveDID, PermissionUpdateDID, PermissionDeactivateDID,
		PermissionCreateKey, PermissionSignData, PermissionRotateKey, PermissionDeleteKey,
		PermissionKYCVerify, PermissionKYCUpload, PermissionKYCCheck,
		PermissionManageUsers, PermissionViewLogs, PermissionManageKeys,
	},
	RoleIssuer: {
		// Issuer can issue, revoke VCs and manage DIDs
		PermissionIssueVC, PermissionRevokeVC, PermissionListVC,
		PermissionCreateDID, PermissionResolveDID, PermissionUpdateDID,
		PermissionCreateKey, PermissionSignData, PermissionRotateKey,
		PermissionKYCVerify, PermissionKYCUpload, PermissionKYCCheck,
	},
	RoleVerifier: {
		// Verifier can verify VCs and resolve DIDs
		PermissionVerifyVC, PermissionListVC,
		PermissionResolveDID,
		PermissionKYCVerify, PermissionKYCCheck,
	},
	RoleUser: {
		// User can only resolve DIDs and verify VCs
		PermissionVerifyVC,
		PermissionResolveDID,
	},
}

// User represents a user with roles and permissions
type User struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Roles       []Role   `json:"roles"`
	APIKey      string   `json:"api_key,omitempty"`
	Active      bool     `json:"active"`
	CreatedAt   string   `json:"created_at"`
	LastLoginAt string   `json:"last_login_at,omitempty"`
}

// HasPermission checks if user has a specific permission
func (u *User) HasPermission(permission Permission) bool {
	for _, role := range u.Roles {
		if permissions, exists := RolePermissions[role]; exists {
			for _, p := range permissions {
				if p == permission {
					return true
				}
			}
		}
	}
	return false
}

// HasRole checks if user has a specific role
func (u *User) HasRole(role Role) bool {
	for _, r := range u.Roles {
		if r == role {
				return true
		}
	}
	return false
}

// GetPermissions returns all permissions for the user
func (u *User) GetPermissions() []Permission {
	permissionSet := make(map[Permission]bool)
	
	for _, role := range u.Roles {
		if permissions, exists := RolePermissions[role]; exists {
			for _, permission := range permissions {
				permissionSet[permission] = true
			}
		}
	}
	
	var permissions []Permission
	for permission := range permissionSet {
		permissions = append(permissions, permission)
	}
	
	return permissions
}

// UserStore interface for user management
type UserStore interface {
	GetUserByAPIKey(ctx context.Context, apiKey string) (*User, error)
	GetUserByID(ctx context.Context, userID string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, userID string) error
	ListUsers(ctx context.Context, limit, offset int) ([]*User, error)
}

// InMemoryUserStore is a simple in-memory implementation of UserStore
type InMemoryUserStore struct {
	users   map[string]*User
	apiKeys map[string]*User
	logger  *zap.Logger
}

// NewInMemoryUserStore creates a new in-memory user store
func NewInMemoryUserStore(logger *zap.Logger) *InMemoryUserStore {
	return &InMemoryUserStore{
		users:   make(map[string]*User),
		apiKeys: make(map[string]*User),
		logger:  logger,
	}
}

// GetUserByAPIKey retrieves a user by API key
func (s *InMemoryUserStore) GetUserByAPIKey(ctx context.Context, apiKey string) (*User, error) {
	user, exists := s.apiKeys[apiKey]
	if !exists {
		return nil, fmt.Errorf("user not found for API key")
	}
	
	if !user.Active {
		return nil, fmt.Errorf("user account is inactive")
	}
	
	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *InMemoryUserStore) GetUserByID(ctx context.Context, userID string) (*User, error) {
	user, exists := s.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// CreateUser creates a new user
func (s *InMemoryUserStore) CreateUser(ctx context.Context, user *User) error {
	if _, exists := s.users[user.ID]; exists {
		return fmt.Errorf("user already exists")
	}
	
	s.users[user.ID] = user
	if user.APIKey != "" {
		s.apiKeys[user.APIKey] = user
	}
	
	s.logger.Info("User created", zap.String("user_id", user.ID), zap.String("username", user.Username))
	return nil
}

// UpdateUser updates an existing user
func (s *InMemoryUserStore) UpdateUser(ctx context.Context, user *User) error {
	oldUser, exists := s.users[user.ID]
	if !exists {
		return fmt.Errorf("user not found")
	}
	
	// Remove old API key mapping
	if oldUser.APIKey != "" {
		delete(s.apiKeys, oldUser.APIKey)
	}
	
	s.users[user.ID] = user
	if user.APIKey != "" {
		s.apiKeys[user.APIKey] = user
	}
	
	s.logger.Info("User updated", zap.String("user_id", user.ID), zap.String("username", user.Username))
	return nil
}

// DeleteUser deletes a user
func (s *InMemoryUserStore) DeleteUser(ctx context.Context, userID string) error {
	user, exists := s.users[userID]
	if !exists {
		return fmt.Errorf("user not found")
	}
	
	delete(s.users, userID)
	if user.APIKey != "" {
		delete(s.apiKeys, user.APIKey)
	}
	
	s.logger.Info("User deleted", zap.String("user_id", userID))
	return nil
}

// ListUsers lists all users with pagination
func (s *InMemoryUserStore) ListUsers(ctx context.Context, limit, offset int) ([]*User, error) {
	var users []*User
	count := 0
	
	for _, user := range s.users {
		if count < offset {
			count++
			continue
		}
		
		if len(users) >= limit {
			break
		}
		
		users = append(users, user)
		count++
	}
	
	return users, nil
}

// RBACMiddleware provides role-based access control
type RBACMiddleware struct {
	userStore UserStore
	logger    *zap.Logger
}

// NewRBACMiddleware creates a new RBAC middleware
func NewRBACMiddleware(userStore UserStore, logger *zap.Logger) *RBACMiddleware {
	return &RBACMiddleware{
		userStore: userStore,
		logger:    logger,
	}
}

// AuthenticateUser middleware authenticates user and loads user info
func (rbac *RBACMiddleware) AuthenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			rbac.logger.Warn("Missing API key", zap.String("ip", getClientIP(c)))
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "API key required",
				"code":  "MISSING_API_KEY",
			})
			c.Abort()
			return
		}
		
		user, err := rbac.userStore.GetUserByAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			rbac.logger.Warn("Invalid API key", 
				zap.String("ip", getClientIP(c)),
				zap.String("error", err.Error()),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
				"code":  "INVALID_API_KEY",
			})
			c.Abort()
			return
		}
		
		// Store user in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Next()
	}
}

// RequirePermission middleware checks if user has required permission
func (rbac *RBACMiddleware) RequirePermission(permission Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
				"code":  "NOT_AUTHENTICATED",
			})
			c.Abort()
			return
		}
		
		user, ok := userInterface.(*User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user context",
				"code":  "INVALID_USER_CONTEXT",
			})
			c.Abort()
			return
		}
		
		if !user.HasPermission(permission) {
			rbac.logger.Warn("Permission denied", 
				zap.String("user_id", user.ID),
				zap.String("permission", string(permission)),
				zap.Strings("user_roles", rolesToStrings(user.Roles)),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
				"code":  "INSUFFICIENT_PERMISSIONS",
				"required_permission": string(permission),
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func (rbac *RBACMiddleware) RequireRole(role Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
				"code":  "NOT_AUTHENTICATED",
			})
			c.Abort()
			return
		}
		
		user, ok := userInterface.(*User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid user context",
				"code":  "INVALID_USER_CONTEXT",
			})
			c.Abort()
			return
		}
		
		if !user.HasRole(role) {
			rbac.logger.Warn("Role access denied", 
				zap.String("user_id", user.ID),
				zap.String("required_role", string(role)),
				zap.Strings("user_roles", rolesToStrings(user.Roles)),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient role",
				"code":  "INSUFFICIENT_ROLE",
				"required_role": string(role),
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// rolesToStrings converts roles to string slice
func rolesToStrings(roles []Role) []string {
	var result []string
	for _, role := range roles {
		result = append(result, string(role))
	}
	return result
}

// GetCurrentUser helper function to get current user from context
func GetCurrentUser(c *gin.Context) (*User, error) {
	userInterface, exists := c.Get("user")
	if !exists {
		return nil, fmt.Errorf("user not found in context")
	}
	
	user, ok := userInterface.(*User)
	if !ok {
		return nil, fmt.Errorf("invalid user type in context")
	}
	
	return user, nil
}