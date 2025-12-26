-- Insert default admin user
-- Email: admin@admin.com
-- Password: admin@123
-- This user will have the admin role assigned

DO $$
DECLARE
    admin_user_id UUID;
    admin_role_id UUID;
    admin_email VARCHAR(255) := 'admin@admin.com';
BEGIN
    -- Get admin role ID by name
    SELECT id INTO admin_role_id FROM roles WHERE name = 'admin';
    
    IF admin_role_id IS NULL THEN
        RAISE EXCEPTION 'Admin role not found. Please ensure migration 002_auth_schema.sql has been executed.';
    END IF;
    
    -- Check if admin user already exists
    SELECT id INTO admin_user_id FROM users WHERE email = admin_email;
    
    IF admin_user_id IS NULL THEN
        -- Create admin user with random UUID
        admin_user_id := gen_random_uuid();
        
        INSERT INTO users (
            id, 
            email, 
            name, 
            password_hash, 
            auth_type, 
            active, 
            created_at, 
            updated_at
        ) VALUES (
            admin_user_id,
            admin_email,
            'Administrator',
            '$2a$10$IIbu/Hx8lQJanbd0Rr3OeunWWVDF.m6PdRErfcFpZbaJkSsNoJX0.', -- bcrypt hash of "admin@123"
            'local',
            true,
            NOW(),
            NOW()
        );
        
        RAISE NOTICE 'Admin user created successfully with email: % and password: admin@123', admin_email;
    ELSE
        RAISE NOTICE 'Admin user already exists, ensuring admin role is assigned';
    END IF;
    
    -- Assign admin role to the user (works for both new and existing users)
    INSERT INTO user_roles (user_id, role_id, assigned_at)
    VALUES (admin_user_id, admin_role_id, NOW())
    ON CONFLICT (user_id, role_id) DO NOTHING;
    
    RAISE NOTICE 'Admin role assigned to user with email: %', admin_email;
END $$;
