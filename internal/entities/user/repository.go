package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/adminoid/vuego/pkg/clients/postgresql"
	"github.com/adminoid/vuego/pkg/logging"
	"github.com/jackc/pgconn"
	"strings"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	//Password     string `json:"password"`
	PasswordHash []byte `json:"password_hash"`
	RefreshToken string `json:"refresh_token"`
}

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

type RepositoryUser interface {
	FindAll(ctx context.Context) (u []User, err error)
	Create(ctx context.Context, u *User) error
}

func NewRepository(client postgresql.Client, logger *logging.Logger) RepositoryUser {
	return &repository{
		client: client,
		logger: logger,
	}
}

func formatQuery(q string) string {
	replacer := strings.NewReplacer("\t", "", "\n", "")
	tmp := replacer.Replace(q)
	fmt.Println(tmp)
	return tmp
	//return replacer.Replace(q)
}

func (r *repository) FindAll(ctx context.Context) (u []User, err error) {
	q := `
		SELECT id, name FROM public.users;
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)

	for rows.Next() {
		var u User

		err = rows.Scan(&u.ID, &u.Email)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) Create(ctx context.Context, u *User) error {
	q := `
		INSERT INTO users
		    (name, email, password_hash) 
		VALUES 
		       ($1, $2, $3) 
		RETURNING id
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, u.Name, u.Email, u.PasswordHash).Scan(&u.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}

	return nil
}
