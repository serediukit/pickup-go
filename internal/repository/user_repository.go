package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"pickup-srv/internal/cache"
	"pickup-srv/internal/models"
	"pickup-srv/proto"
)

type UserRepository struct {
	db    *sql.DB
	cache cache.Cacher
}

func NewUserRepository(db *sql.DB, cacheClient cache.Cacher) *UserRepository {
	return &UserRepository{db: db, cache: cacheClient}
}

func (r *UserRepository) GetUsers(params *proto.UserSearchParams, limit int32) ([]*proto.User, error) {
	cacheKey := strconv.Itoa(int(params.Id))
	lastQueryResult, err := r.cache.Get(context.Background(), cacheKey)
	if err != nil {
		log.Printf("No cache exists for user %d\n", params.Id)
	}
	log.Printf("Last query length for user %d - %s\n", params.Id, lastQueryResult)

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

	var users []*proto.User
	for rows.Next() {
		user := &proto.User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.City, &user.RegDt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	err = r.cache.Set(context.Background(), cacheKey, len(users), 6*time.Hour)
	log.Printf("Set last query length for user %d - %d\n", params.Id, len(users))

	return users, nil
}

func (r *UserRepository) CreateUser(user *models.UserRegistrationEvent) error {
	query := `
INSERT INTO users (name, age, city, gender, search_gender, search_age_from, search_age_to, location) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(query, user.Name, user.Age, user.City, user.Gender, user.SearchGender, user.SearchAgeFrom, user.SearchAgeTo, user.Location)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
