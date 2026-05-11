package store

import (
	"context"
	"homework/internal/model"
	"log"

	"github.com/jackc/pgx/v5"
)

type Store struct {
	db *pgx.Conn
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(ctx context.Context, user model.User) error {
	query := `
	INSERT INTO users (name, email)
	VALUES ($1, $2)
	`

	_, err := s.db.Exec(ctx, query, user.Name, user.Email)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Store) DeleteUser(ctx context.Context, id int) error {
	query := `
	DELETE FROM users
	WHERE id = $1
	`

	_, err := s.db.Exec(ctx, query, id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Store) GetUser(ctx context.Context, id int) (model.User, error) {
	query := `
	SELECT id, name, email
	FROM users
	WHERE id = $1
	`

	var user model.User
	err := s.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		log.Println(err)
		return model.User{}, err
	}

	return user, nil
}

func (s *Store) UpdateUser(ctx context.Context, user model.User) error {
	query := `
	UPDATE users
	SET name = $1, email = $2
	WHERE id = $3
	`

	_, err := s.db.Exec(ctx, query, user.Name, user.Email, user.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Store) GetUsersList(ctx context.Context) ([]model.User, error) {
	query := `
	SELECT id, name, email
	FROM users
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
