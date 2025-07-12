-- Migration: 005_rename_landlord_to_agent.sql
-- Drop landlord_id column and keep agent_id in properties table

-- First, drop the foreign key constraint on landlord_id
ALTER TABLE properties DROP CONSTRAINT IF EXISTS fk_users_properties;

-- Drop the landlord_id column (agent_id already exists)
ALTER TABLE properties DROP COLUMN IF EXISTS landlord_id;

-- Add foreign key constraint for agent_id if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'fk_users_agent_properties' 
        AND conrelid = 'properties'::regclass
    ) THEN
        ALTER TABLE properties ADD CONSTRAINT fk_users_agent_properties 
            FOREIGN KEY (agent_id) REFERENCES users(id);
    END IF;
END $$;

-- Update the user_type check constraint to include 'agent' instead of 'landlord'
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_user_type_check;
ALTER TABLE users ADD CONSTRAINT users_user_type_check 
    CHECK (user_type IN ('landlord', 'tenant', 'agent', 'admin'));
