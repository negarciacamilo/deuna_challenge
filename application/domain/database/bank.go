package database

import (
	"context"
	"github.com/uptrace/bun"
)

type Bank struct {
	Base
	Name string `json:"name"`
}

func (*Bank) AfterCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
	banks := []Bank{{Base: Base{ID: 1}, Name: "Santander"}, {Base: Base{ID: 2}, Name: "BBVA"}, {Base: Base{ID: 3}, Name: "HSBC"}}
	for _, bank := range banks {
		query.DB().NewInsert().Model(&bank).On("CONFLICT (id) DO UPDATE").Exec(ctx)
	}
	return nil
}
