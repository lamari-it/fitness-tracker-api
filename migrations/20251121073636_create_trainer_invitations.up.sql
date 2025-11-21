-- Create trainer_invitations table for email-based invitations
CREATE TABLE IF NOT EXISTS trainer_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trainer_id UUID NOT NULL,
    invitee_email VARCHAR(255) NOT NULL,
    invitation_token VARCHAR(64) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    accepted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_trainer_invitations_trainer
        FOREIGN KEY (trainer_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- Indexes for efficient querying
CREATE INDEX idx_trainer_invitations_trainer_id ON trainer_invitations(trainer_id);
CREATE INDEX idx_trainer_invitations_invitee_email ON trainer_invitations(invitee_email);
CREATE INDEX idx_trainer_invitations_token ON trainer_invitations(invitation_token);
CREATE INDEX idx_trainer_invitations_status ON trainer_invitations(status);
CREATE INDEX idx_trainer_invitations_deleted_at ON trainer_invitations(deleted_at);

-- Unique constraint: one pending invitation per trainer-email combination
CREATE UNIQUE INDEX idx_trainer_invitations_unique_pending
    ON trainer_invitations(trainer_id, invitee_email)
    WHERE status = 'pending' AND deleted_at IS NULL;
