package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"unique;not null" json:"name"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []User       `gorm:"many2many:user_roles;" json:"users,omitempty"`
	ParentRoles []Role       `gorm:"many2many:role_inheritance;joinForeignKey:child_role_id;joinReferences:parent_role_id" json:"parent_roles,omitempty"`
	ChildRoles  []Role       `gorm:"many2many:role_inheritance;joinForeignKey:parent_role_id;joinReferences:child_role_id" json:"child_roles,omitempty"`
}

type RoleInheritance struct {
	ChildRoleID  uint      `gorm:"primaryKey" json:"child_role_id"`
	ParentRoleID uint      `gorm:"primaryKey" json:"parent_role_id"`
	CreatedAt    time.Time `json:"created_at"`

	ChildRole  Role `gorm:"foreignKey:ChildRoleID" json:"child_role,omitempty"`
	ParentRole Role `gorm:"foreignKey:ParentRoleID" json:"parent_role,omitempty"`
}

func (RoleInheritance) TableName() string {
	return "role_inheritance"
}

type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"unique;not null" json:"name"`
	Resource    string         `gorm:"not null" json:"resource"`
	Action      string         `gorm:"not null" json:"action"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Roles []Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

type RolePermission struct {
	RoleID       uint           `gorm:"primaryKey" json:"role_id"`
	PermissionID uint           `gorm:"primaryKey" json:"permission_id"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Role       Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permission Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

type UserRole struct {
	UserID    uuid.UUID      `gorm:"primaryKey" json:"user_id"`
	RoleID    uint           `gorm:"primaryKey" json:"role_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func (UserRole) TableName() string {
	return "user_roles"
}

// Response DTOs for RBAC

type PermissionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

type RoleResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	ParentRoles []string             `json:"parent_roles,omitempty"`
}

type EffectivePermissionsResponse struct {
	RoleID               uint                 `json:"role_id"`
	RoleName             string               `json:"role_name"`
	DirectPermissions    []PermissionResponse `json:"direct_permissions"`
	InheritedPermissions []PermissionResponse `json:"inherited_permissions"`
	EffectivePermissions []PermissionResponse `json:"effective_permissions"`
	InheritanceChain     []string             `json:"inheritance_chain"`
}

// ToResponse converts Permission to response format
func (p *Permission) ToResponse() PermissionResponse {
	return PermissionResponse{
		ID:          p.ID,
		Name:        p.Name,
		Resource:    p.Resource,
		Action:      p.Action,
		Description: p.Description,
	}
}

// ToResponse converts Role to response format
func (r *Role) ToResponse() RoleResponse {
	permissions := make([]PermissionResponse, 0, len(r.Permissions))
	for _, p := range r.Permissions {
		permissions = append(permissions, p.ToResponse())
	}

	parentRoleNames := make([]string, 0, len(r.ParentRoles))
	for _, pr := range r.ParentRoles {
		parentRoleNames = append(parentRoleNames, pr.Name)
	}

	return RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: permissions,
		ParentRoles: parentRoleNames,
	}
}
