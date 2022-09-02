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
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
	//RefreshToken string `json:"refresh_token"`
}

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

type RepositoryUser interface {
	FindAll(ctx context.Context) (u []User, err error)
	Create(ctx context.Context, u *User) error
	Get(ctx context.Context, email string) (User, error)
	UpdateRefreshToken(ctx context.Context, userId string, rt string) error
	GetUserByRt(ctx context.Context, rt string) (WithRtUser, error)
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

func (r *repository) UpdateRefreshToken(ctx context.Context, userId string, rt string) error {
	q := `
		UPDATE users SET refresh_token = $1
  		WHERE id = $2;	
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if _, err := r.client.Query(ctx, q, rt, userId); err != nil {
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

func (r *repository) FindAll(ctx context.Context) (u []User, err error) {
	q := `
		SELECT id, name FROM users;
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

func (r *repository) Get(ctx context.Context, email string) (User, error) {
	q := `
		SELECT id, name, email, password_hash FROM public.users WHERE email = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var user User
	err := r.client.QueryRow(ctx, q, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

type WithRtUser struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	RefreshToken string `json:"refresh_token"`
}

func (r *repository) GetUserByRt(ctx context.Context, rt string) (WithRtUser, error) {
	q := `
		SELECT id, name, email, refresh_token FROM public.users WHERE refresh_token = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var user WithRtUser
	err := r.client.QueryRow(ctx, q, rt).Scan(&user.ID, &user.Name, &user.Email, &user.RefreshToken)
	if err != nil {
		return WithRtUser{}, err
	}

	return user, nil
}
