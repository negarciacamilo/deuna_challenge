package database

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/uptrace/bun"
)

type Customer struct {
	Base
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
}

func (*Customer) AfterCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
	for i := 1; i < 11; i++ {
		customer := Customer{
			Base: Base{
				ID: uint64(i),
			},
			Name:     gofakeit.Name(),
			LastName: gofakeit.LastName(),
			Email:    gofakeit.Email(),
		}
		query.DB().NewInsert().Model(&customer).On("CONFLICT (id) DO UPDATE").Exec(ctx)
	}
	return nil
}
