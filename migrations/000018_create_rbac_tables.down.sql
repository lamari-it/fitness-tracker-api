-- Drop indexes
DROP INDEX IF EXISTS idx_user_roles_role_id;
DROP INDEX IF EXISTS idx_user_roles_user_id;
DROP INDEX IF EXISTS idx_role_permissions_permission_id;
DROP INDEX IF EXISTS idx_role_permissions_role_id;
DROP INDEX IF EXISTS idx_permissions_resource_action;
DROP INDEX IF EXISTS idx_permissions_deleted_at;
DROP INDEX IF EXISTS idx_roles_deleted_at;

-- Drop junction tables first (due to foreign key constraints)
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS role_permissions;

-- Drop main tables
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;