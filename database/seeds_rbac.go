package database

import (
	"fit-flow-api/models"
	"log"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) error {

	roles := []models.Role{
		{Name: "admin", Description: "Admin role"},
		{Name: "trainer", Description: "Trainer role"},
		{Name: "user", Description: "User role"},
	}

	for _, role := range roles {
		var existingRole models.Role
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&role).Error; err != nil {
					log.Printf("Error creating role %s: %v", role.Name, err)
					return err
				}
				log.Printf("Created role: %s", role.Name)
			} else {
				return err
			}
		} else {
			log.Printf("Role already exists: %s", role.Name)
		}
	}

	// Set up role inheritance
	if err := seedRoleInheritance(db); err != nil {
		return err
	}

	return nil
}

// seedRoleInheritance sets up the role inheritance hierarchy
// trainer inherits from user
// admin inherits from trainer
func seedRoleInheritance(db *gorm.DB) error {
	// Get roles
	var userRole, trainerRole, adminRole models.Role
	if err := db.Where("name = ?", "user").First(&userRole).Error; err != nil {
		log.Printf("User role not found: %v", err)
		return err
	}
	if err := db.Where("name = ?", "trainer").First(&trainerRole).Error; err != nil {
		log.Printf("Trainer role not found: %v", err)
		return err
	}
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		log.Printf("Admin role not found: %v", err)
		return err
	}

	// Define inheritance relationships
	inheritanceMap := []struct {
		childID  uint
		parentID uint
		childName string
		parentName string
	}{
		{trainerRole.ID, userRole.ID, "trainer", "user"},
		{adminRole.ID, trainerRole.ID, "admin", "trainer"},
	}

	for _, inheritance := range inheritanceMap {
		var existing models.RoleInheritance
		if err := db.Where("child_role_id = ? AND parent_role_id = ?", inheritance.childID, inheritance.parentID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				roleInheritance := models.RoleInheritance{
					ChildRoleID:  inheritance.childID,
					ParentRoleID: inheritance.parentID,
				}
				if err := db.Create(&roleInheritance).Error; err != nil {
					log.Printf("Error creating role inheritance %s -> %s: %v", inheritance.childName, inheritance.parentName, err)
					return err
				}
				log.Printf("Created role inheritance: %s inherits from %s", inheritance.childName, inheritance.parentName)
			} else {
				return err
			}
		} else {
			log.Printf("Role inheritance already exists: %s inherits from %s", inheritance.childName, inheritance.parentName)
		}
	}

	return nil
}

func SeedRBACData(db *gorm.DB) error {
	// Define all permissions based on existing endpoints
	permissions := []models.Permission{
		// Auth permissions
		{Name: "auth.register", Resource: "auth", Action: "register", Description: "Register new users"},
		{Name: "auth.login", Resource: "auth", Action: "login", Description: "User login"},
		{Name: "auth.profile.read", Resource: "auth", Action: "read_profile", Description: "View user profile"},

		// Dashboard permissions
		{Name: "dashboard.view", Resource: "dashboard", Action: "view", Description: "View dashboard"},

		// Muscle Groups permissions
		{Name: "muscle_groups.create", Resource: "muscle_groups", Action: "create", Description: "Create muscle groups"},
		{Name: "muscle_groups.read", Resource: "muscle_groups", Action: "read", Description: "View muscle groups"},
		{Name: "muscle_groups.update", Resource: "muscle_groups", Action: "update", Description: "Update muscle groups"},
		{Name: "muscle_groups.delete", Resource: "muscle_groups", Action: "delete", Description: "Delete muscle groups"},

		// Exercises permissions
		{Name: "exercises.create", Resource: "exercises", Action: "create", Description: "Create exercises"},
		{Name: "exercises.read", Resource: "exercises", Action: "read", Description: "View exercises"},
		{Name: "exercises.update", Resource: "exercises", Action: "update", Description: "Update exercises"},
		{Name: "exercises.delete", Resource: "exercises", Action: "delete", Description: "Delete exercises"},
		{Name: "exercises.muscle_groups.manage", Resource: "exercises", Action: "manage_muscle_groups", Description: "Manage exercise muscle groups"},
		{Name: "exercises.equipment.manage", Resource: "exercises", Action: "manage_equipment", Description: "Manage exercise equipment"},

		// Equipment permissions
		{Name: "equipment.create", Resource: "equipment", Action: "create", Description: "Create equipment"},
		{Name: "equipment.read", Resource: "equipment", Action: "read", Description: "View equipment"},
		{Name: "equipment.update", Resource: "equipment", Action: "update", Description: "Update equipment"},
		{Name: "equipment.delete", Resource: "equipment", Action: "delete", Description: "Delete equipment"},

		// Fitness Levels permissions
		{Name: "fitness_levels.create", Resource: "fitness_levels", Action: "create", Description: "Create fitness levels"},
		{Name: "fitness_levels.read", Resource: "fitness_levels", Action: "read", Description: "View fitness levels"},
		{Name: "fitness_levels.update", Resource: "fitness_levels", Action: "update", Description: "Update fitness levels"},
		{Name: "fitness_levels.delete", Resource: "fitness_levels", Action: "delete", Description: "Delete fitness levels"},

		// Fitness Goals permissions
		{Name: "fitness_goals.create", Resource: "fitness_goals", Action: "create", Description: "Create fitness goals"},
		{Name: "fitness_goals.read", Resource: "fitness_goals", Action: "read", Description: "View fitness goals"},
		{Name: "fitness_goals.update", Resource: "fitness_goals", Action: "update", Description: "Update fitness goals"},
		{Name: "fitness_goals.delete", Resource: "fitness_goals", Action: "delete", Description: "Delete fitness goals"},

		// User Fitness Settings permissions
		{Name: "user_fitness.goals.read", Resource: "user_fitness", Action: "read_goals", Description: "View user fitness goals"},
		{Name: "user_fitness.goals.update", Resource: "user_fitness", Action: "update_goals", Description: "Update user fitness goals"},
		{Name: "user_fitness.level.update", Resource: "user_fitness", Action: "update_level", Description: "Update user fitness level"},

		// User Equipment permissions
		{Name: "user_equipment.create", Resource: "user_equipment", Action: "create", Description: "Add user equipment"},
		{Name: "user_equipment.read", Resource: "user_equipment", Action: "read", Description: "View user equipment"},
		{Name: "user_equipment.update", Resource: "user_equipment", Action: "update", Description: "Update user equipment"},
		{Name: "user_equipment.delete", Resource: "user_equipment", Action: "delete", Description: "Remove user equipment"},

		// Workout Plans permissions
		{Name: "workout_plans.create", Resource: "workout_plans", Action: "create", Description: "Create workout plans"},
		{Name: "workout_plans.read", Resource: "workout_plans", Action: "read", Description: "View workout plans"},
		{Name: "workout_plans.update", Resource: "workout_plans", Action: "update", Description: "Update workout plans"},
		{Name: "workout_plans.delete", Resource: "workout_plans", Action: "delete", Description: "Delete workout plans"},

		// Friends permissions
		{Name: "friends.request.send", Resource: "friends", Action: "send_request", Description: "Send friend requests"},
		{Name: "friends.request.view", Resource: "friends", Action: "view_requests", Description: "View friend requests"},
		{Name: "friends.request.respond", Resource: "friends", Action: "respond_request", Description: "Respond to friend requests"},
		{Name: "friends.view", Resource: "friends", Action: "view", Description: "View friends list"},
		{Name: "friends.remove", Resource: "friends", Action: "remove", Description: "Remove friends"},

		// Translations permissions (Admin only)
		{Name: "translations.create", Resource: "translations", Action: "create", Description: "Create translations"},
		{Name: "translations.read", Resource: "translations", Action: "read", Description: "View translations"},
		{Name: "translations.update", Resource: "translations", Action: "update", Description: "Update translations"},
		{Name: "translations.delete", Resource: "translations", Action: "delete", Description: "Delete translations"},

		// Legacy endpoints permissions
		{Name: "workouts.view", Resource: "workouts", Action: "view", Description: "View workouts"},
		{Name: "nutrition.view", Resource: "nutrition", Action: "view", Description: "View nutrition data"},
	}

	// Create permissions
	for _, permission := range permissions {
		var existingPermission models.Permission
		if err := db.Where("name = ?", permission.Name).First(&existingPermission).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&permission).Error; err != nil {
					log.Printf("Error creating permission %s: %v", permission.Name, err)
					return err
				}
				log.Printf("Created permission: %s", permission.Name)
			} else {
				return err
			}
		} else {
			log.Printf("Permission already exists: %s", permission.Name)
		}
	}

	// Define role-permission mappings
	// With inheritance: trainer inherits from user, admin inherits from trainer
	// Only assign permissions that are unique to each role (not inherited)
	rolePermissionMappings := map[string][]string{
		"admin": {
			// Admin-specific permissions (inherits trainer + user permissions)
			"auth.register",
			"muscle_groups.create", "muscle_groups.update", "muscle_groups.delete",
			"exercises.delete",
			"equipment.create", "equipment.update", "equipment.delete",
			"fitness_levels.create", "fitness_levels.update", "fitness_levels.delete",
			"fitness_goals.create", "fitness_goals.update", "fitness_goals.delete",
			"translations.create", "translations.read", "translations.update", "translations.delete",
		},
		"trainer": {
			// Trainer-specific permissions (inherits user permissions)
			"exercises.create", "exercises.update",
			"exercises.muscle_groups.manage", "exercises.equipment.manage",
			"workout_plans.create", "workout_plans.update", "workout_plans.delete",
		},
		"user": {
			// Base user permissions
			"auth.login", "auth.profile.read",
			"dashboard.view",
			"muscle_groups.read",
			"exercises.read",
			"equipment.read",
			"fitness_levels.read",
			"fitness_goals.read",
			"user_fitness.goals.read", "user_fitness.goals.update", "user_fitness.level.update",
			"user_equipment.create", "user_equipment.read", "user_equipment.update", "user_equipment.delete",
			"workout_plans.read",
			"friends.request.send", "friends.request.view", "friends.request.respond", "friends.view", "friends.remove",
			"workouts.view", "nutrition.view",
		},
	}

	// Assign permissions to roles
	for roleName, permissionNames := range rolePermissionMappings {
		var role models.Role
		if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
			log.Printf("Role %s not found, skipping permission assignment", roleName)
			continue
		}

		for _, permName := range permissionNames {
			var permission models.Permission
			if err := db.Where("name = ?", permName).First(&permission).Error; err != nil {
				log.Printf("Permission %s not found, skipping", permName)
				continue
			}

			// Check if association already exists
			var rolePermission models.RolePermission
			if err := db.Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).First(&rolePermission).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					// Create the association
					rolePermission = models.RolePermission{
						RoleID:       role.ID,
						PermissionID: permission.ID,
					}
					if err := db.Create(&rolePermission).Error; err != nil {
						log.Printf("Error creating role-permission association: %v", err)
						return err
					}
					log.Printf("Assigned permission %s to role %s", permName, roleName)
				} else {
					return err
				}
			}
		}
	}

	log.Println("RBAC seed completed successfully")
	return nil
}

// MigrateExistingUsersToRoles assigns default roles to existing users based on is_admin flag
func MigrateExistingUsersToRoles(db *gorm.DB) error {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	var adminRole, userRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "user").First(&userRole).Error; err != nil {
		return err
	}

	for _, user := range users {
		var roleToAssign models.Role
		if user.IsAdmin {
			roleToAssign = adminRole
		} else {
			roleToAssign = userRole
		}

		// Check if user already has this role
		var userRole models.UserRole
		if err := db.Where("user_id = ? AND role_id = ?", user.ID, roleToAssign.ID).First(&userRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Assign role to user
				userRole = models.UserRole{
					UserID: user.ID,
					RoleID: roleToAssign.ID,
				}
				if err := db.Create(&userRole).Error; err != nil {
					log.Printf("Error assigning role to user %s: %v", user.Email, err)
					return err
				}
				log.Printf("Assigned %s role to user %s", roleToAssign.Name, user.Email)
			} else {
				return err
			}
		}
	}

	log.Println("User role migration completed successfully")
	return nil
}
