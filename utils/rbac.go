package utils

import (
	"fit-flow-api/models"

	"gorm.io/gorm"
)

// GetEffectivePermissions computes all permissions for a role including inherited ones
func GetEffectivePermissions(db *gorm.DB, roleID uint) (*models.EffectivePermissionsResponse, error) {
	// Load the role with its direct permissions and parent roles
	var role models.Role
	if err := db.Preload("Permissions").Preload("ParentRoles").First(&role, roleID).Error; err != nil {
		return nil, err
	}

	// Get direct permissions
	directPermissions := make([]models.PermissionResponse, 0, len(role.Permissions))
	for _, p := range role.Permissions {
		directPermissions = append(directPermissions, p.ToResponse())
	}

	// Get all parent roles recursively
	allParentRoles, inheritanceChain := getAllParentRolesRecursive(db, role.ID, make(map[uint]bool))

	// Collect inherited permissions
	inheritedPermMap := make(map[uint]models.PermissionResponse)
	for _, parentRole := range allParentRoles {
		var parent models.Role
		if err := db.Preload("Permissions").First(&parent, parentRole.ID).Error; err != nil {
			continue
		}
		for _, p := range parent.Permissions {
			if _, exists := inheritedPermMap[p.ID]; !exists {
				inheritedPermMap[p.ID] = p.ToResponse()
			}
		}
	}

	inheritedPermissions := make([]models.PermissionResponse, 0, len(inheritedPermMap))
	for _, p := range inheritedPermMap {
		inheritedPermissions = append(inheritedPermissions, p)
	}

	// Combine direct and inherited for effective permissions (deduplicated)
	effectivePermMap := make(map[uint]models.PermissionResponse)
	for _, p := range directPermissions {
		effectivePermMap[p.ID] = p
	}
	for _, p := range inheritedPermissions {
		effectivePermMap[p.ID] = p
	}

	effectivePermissions := make([]models.PermissionResponse, 0, len(effectivePermMap))
	for _, p := range effectivePermMap {
		effectivePermissions = append(effectivePermissions, p)
	}

	return &models.EffectivePermissionsResponse{
		RoleID:               role.ID,
		RoleName:             role.Name,
		DirectPermissions:    directPermissions,
		InheritedPermissions: inheritedPermissions,
		EffectivePermissions: effectivePermissions,
		InheritanceChain:     inheritanceChain,
	}, nil
}

// getAllParentRolesRecursive recursively gets all parent roles avoiding cycles
func getAllParentRolesRecursive(db *gorm.DB, roleID uint, visited map[uint]bool) ([]models.Role, []string) {
	if visited[roleID] {
		return nil, nil
	}
	visited[roleID] = true

	var role models.Role
	if err := db.Preload("ParentRoles").First(&role, roleID).Error; err != nil {
		return nil, nil
	}

	var allParents []models.Role
	var chain []string

	for _, parent := range role.ParentRoles {
		if !visited[parent.ID] {
			allParents = append(allParents, parent)
			chain = append(chain, parent.Name)

			// Recursively get parents of this parent
			grandParents, grandChain := getAllParentRolesRecursive(db, parent.ID, visited)
			allParents = append(allParents, grandParents...)
			chain = append(chain, grandChain...)
		}
	}

	return allParents, chain
}

// GetAllParentRoles returns all ancestor roles for a given role
func GetAllParentRoles(db *gorm.DB, roleID uint) ([]models.Role, error) {
	parents, _ := getAllParentRolesRecursive(db, roleID, make(map[uint]bool))
	return parents, nil
}

// HasPermission checks if a role has a specific permission (including inherited)
func HasPermission(db *gorm.DB, roleID uint, permissionName string) (bool, error) {
	effective, err := GetEffectivePermissions(db, roleID)
	if err != nil {
		return false, err
	}

	for _, p := range effective.EffectivePermissions {
		if p.Name == permissionName {
			return true, nil
		}
	}

	return false, nil
}

// HasAnyPermission checks if a role has any of the specified permissions
func HasAnyPermission(db *gorm.DB, roleID uint, permissionNames []string) (bool, error) {
	effective, err := GetEffectivePermissions(db, roleID)
	if err != nil {
		return false, err
	}

	permSet := make(map[string]bool)
	for _, name := range permissionNames {
		permSet[name] = true
	}

	for _, p := range effective.EffectivePermissions {
		if permSet[p.Name] {
			return true, nil
		}
	}

	return false, nil
}

// GetUserEffectivePermissions gets all effective permissions for a user across all their roles
func GetUserEffectivePermissions(db *gorm.DB, userID interface{}) ([]models.PermissionResponse, error) {
	var user models.User
	if err := db.Preload("Roles").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	effectivePermMap := make(map[uint]models.PermissionResponse)

	for _, role := range user.Roles {
		effective, err := GetEffectivePermissions(db, role.ID)
		if err != nil {
			continue
		}

		for _, p := range effective.EffectivePermissions {
			effectivePermMap[p.ID] = p
		}
	}

	permissions := make([]models.PermissionResponse, 0, len(effectivePermMap))
	for _, p := range effectivePermMap {
		permissions = append(permissions, p)
	}

	return permissions, nil
}
