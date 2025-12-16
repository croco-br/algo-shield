package services

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.getUserByEmail(ctx, email, false)
}

func (s *UserService) GetUserByEmailWithPassword(ctx context.Context, email string) (*models.User, error) {
	return s.getUserByEmail(ctx, email, true)
}

func (s *UserService) getUserByEmail(ctx context.Context, email string, includePassword bool) (*models.User, error) {
	var user models.User
	var passwordHash, googleID, pictureURL sql.NullString
	var lastLoginAt sql.NullTime

	query := `
		SELECT id, email, name, password_hash, google_id, picture_url, auth_type, active, created_at, updated_at, last_login_at
		FROM users
		WHERE email = $1
	`

	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Name, &passwordHash, &googleID, &pictureURL,
		&user.AuthType, &user.Active, &user.CreatedAt, &user.UpdatedAt, &lastLoginAt,
	)
	if err != nil {
		return nil, err
	}

	if includePassword && passwordHash.Valid {
		user.PasswordHash = &passwordHash.String
	}
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if pictureURL.Valid {
		user.PictureURL = &pictureURL.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	// Load roles and groups
	if err := s.loadUserRoles(ctx, &user); err != nil {
		return nil, err
	}
	if err := s.loadUserGroups(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	var passwordHash, googleID, pictureURL sql.NullString
	var lastLoginAt sql.NullTime

	query := `
		SELECT id, email, name, password_hash, google_id, picture_url, auth_type, active, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1
	`

	err := s.db.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.Name, &passwordHash, &googleID, &pictureURL,
		&user.AuthType, &user.Active, &user.CreatedAt, &user.UpdatedAt, &lastLoginAt,
	)
	if err != nil {
		return nil, err
	}

	// Don't include password hash in normal queries
	if googleID.Valid {
		user.GoogleID = &googleID.String
	}
	if pictureURL.Valid {
		user.PictureURL = &pictureURL.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	// Load roles and groups
	if err := s.loadUserRoles(ctx, &user); err != nil {
		return nil, err
	}
	if err := s.loadUserGroups(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) CreateUser(ctx context.Context, email, name, passwordHash string) (*models.User, error) {
	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		Name:         name,
		PasswordHash: &passwordHash,
		AuthType:     models.AuthTypeLocal,
		Active:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Insert user
	query := `
		INSERT INTO users (id, email, name, password_hash, auth_type, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.Exec(ctx, query,
		user.ID, user.Email, user.Name, user.PasswordHash,
		user.AuthType, user.Active, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Assign default viewer role to new users (by name)
	// This is non-critical - user creation succeeds even if role assignment fails
	var viewerRoleID uuid.UUID
	roleQuery := `SELECT id FROM roles WHERE name = 'viewer' LIMIT 1`
	err = tx.QueryRow(ctx, roleQuery).Scan(&viewerRoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Warning: 'viewer' role not found, user %s created without default role", user.ID)
		} else {
			log.Printf("Error querying for 'viewer' role for user %s: %v", user.ID, err)
		}
	} else {
		// Role found, assign it within the transaction
		assignRoleQuery := `
			INSERT INTO user_roles (user_id, role_id, assigned_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, role_id) DO NOTHING
		`
		_, err = tx.Exec(ctx, assignRoleQuery, user.ID, viewerRoleID, time.Now())
		if err != nil {
			log.Printf("Error assigning default 'viewer' role to user %s: %v", user.ID, err)
			// Continue with user creation even if role assignment fails
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// Load roles and groups (read-only operations, outside transaction)
	if err := s.loadUserRoles(ctx, user); err != nil {
		return nil, err
	}
	if err := s.loadUserGroups(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateLastLogin(ctx context.Context, userID uuid.UUID, lastLoginAt *time.Time) error {
	query := `UPDATE users SET last_login_at = $1, updated_at = $2 WHERE id = $3`
	_, err := s.db.Exec(ctx, query, lastLoginAt, time.Now(), userID)
	return err
}

func (s *UserService) loadUserRoles(ctx context.Context, user *models.User) error {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
		UNION
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN group_roles gr ON r.id = gr.role_id
		INNER JOIN user_groups ug ON gr.group_id = ug.group_id
		WHERE ug.user_id = $1
	`

	rows, err := s.db.Query(ctx, query, user.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	roles := make([]models.Role, 0)
	roleMap := make(map[uuid.UUID]bool) // To avoid duplicates

	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
			continue
		}
		if !roleMap[role.ID] {
			roles = append(roles, role)
			roleMap[role.ID] = true
		}
	}

	user.Roles = roles
	return nil
}

func (s *UserService) loadUserGroups(ctx context.Context, user *models.User) error {
	query := `
		SELECT g.id, g.name, g.description, g.created_at, g.updated_at
		FROM groups g
		INNER JOIN user_groups ug ON g.id = ug.group_id
		WHERE ug.user_id = $1
	`

	rows, err := s.db.Query(ctx, query, user.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	groups := make([]models.Group, 0)
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatedAt, &group.UpdatedAt); err != nil {
			continue
		}
		groups = append(groups, group)
	}

	user.Groups = groups
	return nil
}

func (s *UserService) HasRole(ctx context.Context, userID uuid.UUID, roleName string) (bool, error) {
	query := `
		SELECT COUNT(*) > 0
		FROM (
			SELECT r.id
			FROM roles r
			INNER JOIN user_roles ur ON r.id = ur.role_id
			WHERE ur.user_id = $1 AND r.name = $2
			UNION
			SELECT r.id
			FROM roles r
			INNER JOIN group_roles gr ON r.id = gr.role_id
			INNER JOIN user_groups ug ON gr.group_id = ug.group_id
			WHERE ug.user_id = $1 AND r.name = $2
		) AS user_roles
	`

	var hasRole bool
	err := s.db.QueryRow(ctx, query, userID, roleName).Scan(&hasRole)
	return hasRole, err
}

func (s *UserService) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, assigned_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`
	_, err := s.db.Exec(ctx, query, userID, roleID, time.Now())
	return err
}

func (s *UserService) RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	_, err := s.db.Exec(ctx, query, userID, roleID)
	return err
}

func (s *UserService) ListUsers(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT id, email, name, password_hash, google_id, picture_url, auth_type, active, created_at, updated_at, last_login_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		var passwordHash, googleID, pictureURL sql.NullString
		var lastLoginAt sql.NullTime

		if err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &passwordHash, &googleID, &pictureURL,
			&user.AuthType, &user.Active, &user.CreatedAt, &user.UpdatedAt, &lastLoginAt,
		); err != nil {
			continue
		}

		// Don't include password hash in list

		if googleID.Valid {
			user.GoogleID = &googleID.String
		}
		if pictureURL.Valid {
			user.PictureURL = &pictureURL.String
		}
		if lastLoginAt.Valid {
			user.LastLoginAt = &lastLoginAt.Time
		}

		// Load roles and groups
		if err := s.loadUserRoles(ctx, &user); err != nil {
			return nil, err
		}
		if err := s.loadUserGroups(ctx, &user); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *UserService) ListRoles(ctx context.Context) ([]models.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		ORDER BY name
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]models.Role, 0)
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
			continue
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (s *UserService) ListGroups(ctx context.Context) ([]models.Group, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM groups
		ORDER BY name
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]models.Group, 0)
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatedAt, &group.UpdatedAt); err != nil {
			continue
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (s *UserService) UpdateUserActive(ctx context.Context, userID uuid.UUID, active bool) error {
	query := `UPDATE users SET active = $1, updated_at = $2 WHERE id = $3`
	_, err := s.db.Exec(ctx, query, active, time.Now(), userID)
	return err
}
