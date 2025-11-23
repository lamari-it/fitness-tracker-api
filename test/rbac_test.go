package test

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
	"testing"
)

func TestRoleInheritance(t *testing.T) {
	// Setup test app (this initializes the database)
	_ = SetupTestApp(t)

	// Seed roles with inheritance
	SeedTestRoles(t)

	// Seed RBAC data (permissions)
	if err := database.SeedRBACData(database.DB); err != nil {
		t.Fatalf("Failed to seed RBAC data: %v", err)
	}

	t.Run("Trainer inherits from User", func(t *testing.T) {
		// Get trainer role
		var trainerRole models.Role
		if err := database.DB.Where("name = ?", "trainer").First(&trainerRole).Error; err != nil {
			t.Fatalf("Failed to get trainer role: %v", err)
		}

		// Get effective permissions
		effective, err := utils.GetEffectivePermissions(database.DB, trainerRole.ID)
		if err != nil {
			t.Fatalf("Failed to get effective permissions: %v", err)
		}

		// Trainer should have user permissions via inheritance
		userPermissions := []string{
			"auth.login",
			"auth.profile.read",
			"dashboard.view",
			"exercises.read",
		}

		for _, permName := range userPermissions {
			found := false
			for _, p := range effective.EffectivePermissions {
				if p.Name == permName {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Trainer should have inherited permission %s from user", permName)
			}
		}

		// Trainer should have its own direct permissions
		trainerPermissions := []string{
			"exercises.create",
			"exercises.update",
			"workout_plans.create",
		}

		for _, permName := range trainerPermissions {
			found := false
			for _, p := range effective.DirectPermissions {
				if p.Name == permName {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Trainer should have direct permission %s", permName)
			}
		}

		// Check inheritance chain
		if len(effective.InheritanceChain) == 0 {
			t.Error("Trainer should have inheritance chain")
		}

		foundUser := false
		for _, roleName := range effective.InheritanceChain {
			if roleName == "user" {
				foundUser = true
				break
			}
		}
		if !foundUser {
			t.Error("Trainer's inheritance chain should include 'user'")
		}
	})

	t.Run("Admin inherits from Trainer and User", func(t *testing.T) {
		// Get admin role
		var adminRole models.Role
		if err := database.DB.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
			t.Fatalf("Failed to get admin role: %v", err)
		}

		// Get effective permissions
		effective, err := utils.GetEffectivePermissions(database.DB, adminRole.ID)
		if err != nil {
			t.Fatalf("Failed to get effective permissions: %v", err)
		}

		// Admin should have user permissions via inheritance (through trainer)
		userPermissions := []string{
			"auth.login",
			"auth.profile.read",
			"exercises.read",
		}

		for _, permName := range userPermissions {
			found := false
			for _, p := range effective.EffectivePermissions {
				if p.Name == permName {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Admin should have inherited permission %s from user (via trainer)", permName)
			}
		}

		// Admin should have trainer permissions via inheritance
		trainerPermissions := []string{
			"exercises.create",
			"exercises.update",
			"workout_plans.create",
		}

		for _, permName := range trainerPermissions {
			found := false
			for _, p := range effective.EffectivePermissions {
				if p.Name == permName {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Admin should have inherited permission %s from trainer", permName)
			}
		}

		// Admin should have its own direct permissions
		adminPermissions := []string{
			"auth.register",
			"muscle_groups.create",
			"translations.create",
		}

		for _, permName := range adminPermissions {
			found := false
			for _, p := range effective.DirectPermissions {
				if p.Name == permName {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Admin should have direct permission %s", permName)
			}
		}

		// Check inheritance chain includes both trainer and user
		if len(effective.InheritanceChain) < 2 {
			t.Error("Admin should have at least 2 roles in inheritance chain")
		}

		foundTrainer := false
		foundUser := false
		for _, roleName := range effective.InheritanceChain {
			if roleName == "trainer" {
				foundTrainer = true
			}
			if roleName == "user" {
				foundUser = true
			}
		}
		if !foundTrainer {
			t.Error("Admin's inheritance chain should include 'trainer'")
		}
		if !foundUser {
			t.Error("Admin's inheritance chain should include 'user'")
		}
	})

	t.Run("User has no inheritance", func(t *testing.T) {
		// Get user role
		var userRole models.Role
		if err := database.DB.Where("name = ?", "user").First(&userRole).Error; err != nil {
			t.Fatalf("Failed to get user role: %v", err)
		}

		// Get effective permissions
		effective, err := utils.GetEffectivePermissions(database.DB, userRole.ID)
		if err != nil {
			t.Fatalf("Failed to get effective permissions: %v", err)
		}

		// User should have no inherited permissions
		if len(effective.InheritedPermissions) > 0 {
			t.Errorf("User should have no inherited permissions, got %d", len(effective.InheritedPermissions))
		}

		// User should have no inheritance chain
		if len(effective.InheritanceChain) > 0 {
			t.Errorf("User should have no inheritance chain, got %v", effective.InheritanceChain)
		}

		// Effective permissions should equal direct permissions
		if len(effective.EffectivePermissions) != len(effective.DirectPermissions) {
			t.Errorf("User's effective permissions should equal direct permissions")
		}
	})

	t.Run("HasPermission works with inheritance", func(t *testing.T) {
		// Get trainer role
		var trainerRole models.Role
		if err := database.DB.Where("name = ?", "trainer").First(&trainerRole).Error; err != nil {
			t.Fatalf("Failed to get trainer role: %v", err)
		}

		// Trainer should have auth.login (inherited from user)
		hasLogin, err := utils.HasPermission(database.DB, trainerRole.ID, "auth.login")
		if err != nil {
			t.Fatalf("HasPermission failed: %v", err)
		}
		if !hasLogin {
			t.Error("Trainer should have auth.login permission (inherited)")
		}

		// Trainer should have exercises.create (direct)
		hasCreate, err := utils.HasPermission(database.DB, trainerRole.ID, "exercises.create")
		if err != nil {
			t.Fatalf("HasPermission failed: %v", err)
		}
		if !hasCreate {
			t.Error("Trainer should have exercises.create permission (direct)")
		}

		// Trainer should NOT have translations.create (admin only)
		hasTranslations, err := utils.HasPermission(database.DB, trainerRole.ID, "translations.create")
		if err != nil {
			t.Fatalf("HasPermission failed: %v", err)
		}
		if hasTranslations {
			t.Error("Trainer should NOT have translations.create permission")
		}
	})

	t.Run("Permissions are deduplicated", func(t *testing.T) {
		// Get admin role
		var adminRole models.Role
		if err := database.DB.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
			t.Fatalf("Failed to get admin role: %v", err)
		}

		// Get effective permissions
		effective, err := utils.GetEffectivePermissions(database.DB, adminRole.ID)
		if err != nil {
			t.Fatalf("Failed to get effective permissions: %v", err)
		}

		// Check for duplicates
		permissionCounts := make(map[string]int)
		for _, p := range effective.EffectivePermissions {
			permissionCounts[p.Name]++
		}

		for name, count := range permissionCounts {
			if count > 1 {
				t.Errorf("Permission %s appears %d times in effective permissions (should be deduplicated)", name, count)
			}
		}
	})
}

func TestGetUserEffectivePermissions(t *testing.T) {
	// Setup test app
	e := SetupTestApp(t)

	// Seed roles with inheritance
	SeedTestRoles(t)

	// Seed RBAC data
	if err := database.SeedRBACData(database.DB); err != nil {
		t.Fatalf("Failed to seed RBAC data: %v", err)
	}

	t.Run("User gets all effective permissions from their roles", func(t *testing.T) {
		// Register a trainer user
		SeedTestSpecialties(t)

		e.POST("/api/v1/auth/register").
			WithJSON(map[string]interface{}{
				"email":            "trainer@test.com",
				"password":         "TestPassword123!",
				"password_confirm": "TestPassword123!",
				"first_name":       "Test",
				"last_name":        "Trainer",
				"is_trainer":       true,
				"trainer_profile": map[string]interface{}{
					"bio":        "Test bio",
					"visibility": "public",
				},
			}).
			Expect().
			Status(201)

		// Get user from database
		var user models.User
		if err := database.DB.Where("email = ?", "trainer@test.com").First(&user).Error; err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}

		// Get user's effective permissions
		permissions, err := utils.GetUserEffectivePermissions(database.DB, user.ID)
		if err != nil {
			t.Fatalf("Failed to get user effective permissions: %v", err)
		}

		// User should have trainer permissions (direct) and user permissions (inherited)
		expectedPermissions := []string{
			"auth.login",           // inherited from user
			"exercises.read",       // inherited from user
			"exercises.create",     // direct trainer permission
			"workout_plans.create", // direct trainer permission
		}

		for _, permName := range expectedPermissions {
			found := false
			for _, p := range permissions {
				if p.Name == permName {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("User should have permission %s", permName)
			}
		}
	})
}
