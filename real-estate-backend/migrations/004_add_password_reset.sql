-- Migration: 004_add_password_reset.sql
-- Add password reset functionality

-- Create password_resets table
CREATE TABLE password_resets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX idx_password_resets_token ON password_resets(token);
CREATE INDEX idx_password_resets_user ON password_resets(user_id);
CREATE INDEX idx_password_resets_expires ON password_resets(expires_at);
CREATE INDEX idx_password_resets_deleted ON password_resets(deleted_at);

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_password_resets_updated_at BEFORE UPDATE ON password_resets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Update users table to fix user_type values (remove 'landlord', add 'admin', 'agent')
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_user_type_check;
ALTER TABLE users ADD CONSTRAINT users_user_type_check 
    CHECK (user_type IN ('admin', 'tenant', 'agent'));

-- Add approval fields for agents
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_approved BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS approved_at TIMESTAMP;
ALTER TABLE users ADD COLUMN IF NOT EXISTS approved_by UUID;
ALTER TABLE users ADD CONSTRAINT fk_users_approved_by 
    FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL;

-- Update properties table to use agent_id instead of landlord_id
ALTER TABLE properties RENAME COLUMN landlord_id TO agent_id;
ALTER TABLE properties DROP CONSTRAINT IF EXISTS properties_landlord_id_fkey;
ALTER TABLE properties ADD CONSTRAINT properties_agent_id_fkey
    FOREIGN KEY (agent_id) REFERENCES users(id) ON DELETE CASCADE;
