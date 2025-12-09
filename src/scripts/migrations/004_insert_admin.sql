-- Insert default admin user
-- Email: admin@admin.com
-- Password: admin
-- This user will have the admin role assigned

DO $$
DECLARE
    admin_user_id UUID;
    admin_role_id UUID;
    admin_email VARCHAR(255) := 'admin@admin.com';
BEGIN
    -- Check if admin user already exists
    SELECT id INTO admin_user_id FROM users WHERE email = admin_email;
    
    IF admin_user_id IS NULL THEN
        -- Get admin role ID by name
        SELECT id INTO admin_role_id FROM roles WHERE name = 'admin';
        
        IF admin_role_id IS NULL THEN
            RAISE EXCEPTION 'Admin role not found. Please ensure migration 002_auth_schema.sql has been executed.';
        END IF;
        
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
            '$2a$10$HzMzEBsY7h9i4q4JNEkOcuYjppiLONMl4PljA9vJamATnBTgkNDyC', -- bcrypt hash of "admin"
            'local',
            true,
            NOW(),
            NOW()
        );
        
        -- Assign admin role to the user
        INSERT INTO user_roles (user_id, role_id, assigned_at)
        VALUES (admin_user_id, admin_role_id, NOW())
        ON CONFLICT (user_id, role_id) DO NOTHING;
        
        RAISE NOTICE 'Admin user created successfully with email: % and password: admin', admin_email;
    ELSE
        RAISE NOTICE 'Admin user already exists, skipping creation';
    END IF;
END $$;
