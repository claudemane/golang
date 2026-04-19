package repository

import (
	"database/sql"
	"fmt"
	"practice5/internal/models"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// allowed columns for order_by (whitelist against SQL injection)
var allowedOrderBy = map[string]string{
	"id":         "u.id",
	"name":       "u.name",
	"email":      "u.email",
	"gender":     "u.gender",
	"birth_date": "u.birth_date",
}

func (r *Repository) GetPaginatedUsers(f models.UserFilter) (models.PaginatedResponse, error) {
	args := []interface{}{}
	argIdx := 1
	conditions := []string{}

	// Dynamic filtering
	if f.ID != nil {
		conditions = append(conditions, fmt.Sprintf("u.id = $%d", argIdx))
		args = append(args, *f.ID)
		argIdx++
	}
	if f.Name != nil {
		conditions = append(conditions, fmt.Sprintf("u.name ILIKE $%d", argIdx))
		args = append(args, "%"+*f.Name+"%")
		argIdx++
	}
	if f.Email != nil {
		conditions = append(conditions, fmt.Sprintf("u.email ILIKE $%d", argIdx))
		args = append(args, "%"+*f.Email+"%")
		argIdx++
	}
	if f.Gender != nil {
		conditions = append(conditions, fmt.Sprintf("u.gender = $%d", argIdx))
		args = append(args, *f.Gender)
		argIdx++
	}
	if f.BirthDate != nil {
		conditions = append(conditions, fmt.Sprintf("u.birth_date::date = $%d", argIdx))
		args = append(args, *f.BirthDate)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Dynamic order_by with whitelist
	orderCol, ok := allowedOrderBy[f.OrderBy]
	if !ok {
		orderCol = "u.id" // default
	}

	// Count total
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM users u %s`, where)
	var totalCount int
	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return models.PaginatedResponse{}, err
	}

	// Paginated fetch
	offset := (f.Page - 1) * f.PageSize
	args = append(args, f.PageSize, offset)

	query := fmt.Sprintf(`
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM users u
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		where, orderCol, argIdx, argIdx+1,
	)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return models.PaginatedResponse{}, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return models.PaginatedResponse{}, err
		}
		users = append(users, u)
	}

	if users == nil {
		users = []models.User{}
	}

	return models.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       f.Page,
		PageSize:   f.PageSize,
	}, nil
}

// GetCommonFriends returns common friends of two users in a single JOIN query (no N+1)
func (r *Repository) GetCommonFriends(userID1, userID2 int) ([]models.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date
		FROM users u
		JOIN user_friends uf1 ON uf1.friend_id = u.id AND uf1.user_id = $1
		JOIN user_friends uf2 ON uf2.friend_id = u.id AND uf2.user_id = $2
		ORDER BY u.id
	`

	rows, err := r.db.Query(query, userID1, userID2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return nil, err
		}
		friends = append(friends, u)
	}

	if friends == nil {
		friends = []models.User{}
	}

	return friends, nil
}
