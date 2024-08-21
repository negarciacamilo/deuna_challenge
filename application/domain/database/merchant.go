package database

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/uptrace/bun"
	"math/rand"
)

type Merchant struct {
	Base
	Name              string `json:"name"`
	BankAccountNumber uint64 `json:"bank_account_number"`
	Email             string `json:"email"`
}

func (*Merchant) AfterCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
	for i := 1; i < 11; i++ {
		merchant := Merchant{
			Base: Base{
				ID: uint64(i),
			},
			Name:              gofakeit.Company(),
			Email:             gofakeit.Email(),
			BankAccountNumber: uint64(rand.Int63n(9223372036854775807)),
		}
		query.DB().NewInsert().Model(&merchant).On("CONFLICT (id) DO UPDATE").Exec(ctx)
	}
	return nil
}
