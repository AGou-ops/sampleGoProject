package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Permissions []string

func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

type PermissionModel struct {
	DB *sql.DB
}

func (p PermissionModel) GetAllForUser(userId int64) (Permissions, error) {
	query := `
	SELECT permissions.code
	FROM permissions
	INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id INNER JOIN users ON users_permissions.user_id = users.id
	WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := p.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions

	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (p PermissionModel) AddForUser(userID int64, code ...string) error {
	query := `
	INSERT INTO users_permissions
	SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, userID, pq.Array(code))

	return err
}
