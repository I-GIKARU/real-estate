-- Migration to create initial admin user
-- Creates the default admin user with credentials:
-- Email: admin@realestate.com
-- Password: Admin123!

-- Create admin user
INSERT INTO users (
    id,
    email,
    password_hash,
    first_name,
    last_name,
    phone_number,
    user_type,
    is_verified,
    is_approved,
    is_active,
    created_at,
    updated_at
) VALUES (
    gen_random_uuid(),
    'admin@realestate.com',
    '$2a$10$SQvttWpt6CqPukw.x00zkeV0b.8v6gELaOYQaN5IHssj/WkHzjlYq', -- correct bcrypt hash of 'Admin123!'
    'System',
    'Administrator',
    '+254700000000',
    'admin',
    true,
    true,
    true,
    NOW(),
    NOW()
) ON CONFLICT (email) DO NOTHING;

-- Add comment for documentation
COMMENT ON TABLE users IS 'Users table with admin, agent, and tenant user types';
