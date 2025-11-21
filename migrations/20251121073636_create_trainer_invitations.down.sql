-- Drop indexes first
DROP INDEX IF EXISTS idx_trainer_invitations_unique_pending;
DROP INDEX IF EXISTS idx_trainer_invitations_deleted_at;
DROP INDEX IF EXISTS idx_trainer_invitations_status;
DROP INDEX IF EXISTS idx_trainer_invitations_token;
DROP INDEX IF EXISTS idx_trainer_invitations_invitee_email;
DROP INDEX IF EXISTS idx_trainer_invitations_trainer_id;

-- Drop the table
DROP TABLE IF EXISTS trainer_invitations;
