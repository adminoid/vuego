package user

import (
	"context"
	"fmt"
	"github.com/adminoid/vuego/pkg/clients/postgresql"
	"github.com/adminoid/vuego/pkg/logging"
	"strings"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

type RepositoryUser interface {
	FindAll(ctx context.Context) (u []User, err error)
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
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

		err = rows.Scan(&u.ID, &u.Name)
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

func NewRepository(client postgresql.Client, logger *logging.Logger) RepositoryUser {
	return &repository{
		client: client,
		logger: logger,
	}
}
