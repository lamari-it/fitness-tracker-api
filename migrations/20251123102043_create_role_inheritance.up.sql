-- Create role_inheritance table for multiple inheritance support
CREATE TABLE role_inheritance (
    child_role_id INT NOT NULL,
    parent_role_id INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (child_role_id, parent_role_id),
    CONSTRAINT fk_child_role FOREIGN KEY (child_role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_parent_role FOREIGN KEY (parent_role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT chk_no_self_reference CHECK (child_role_id != parent_role_id)
);

-- Index for efficient lookups
CREATE INDEX idx_role_inheritance_child ON role_inheritance(child_role_id);
CREATE INDEX idx_role_inheritance_parent ON role_inheritance(parent_role_id);
