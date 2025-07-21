package repository

import (
	"database/sql"
	"fmt"
	"pickup-srv/internal/models"
	"pickup-srv/proto"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUsers(params *proto.UserSearchParams, limit int32) ([]*models.User, error) {
	query := `
SELECT id, name, age, city, reg_dt 
FROM users 
WHERE gender = $1
  AND search_gender = $2
  AND age BETWEEN $3 AND $4
  AND search_age_from <= $5
  AND search_age_to >= $6
  AND ABS(location - $7) <= 100
  AND id != $8
ORDER BY ABS(location - $9), reg_dt DESC
LIMIT $10`

	rows, err := r.db.Query(
		query,
		params.SearchGender,
		params.Gender,
		params.SearchAgeFrom,
		params.SearchAgeTo,
		params.Age,
		params.Age,
		params.Location,
		params.Id,
		params.Location,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.City, &user.RegDt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return users, nil
}
