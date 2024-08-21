package database

import (
	"github.com/uptrace/bun"
	"time"
)

type Base struct {
	ID uint64 `json:"id" bun:"id,pk,autoincrement"`
	BaseDates
}

type BaseDates struct {
	CreatedAt time.Time    `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt bun.NullTime `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt *time.Time   `json:"deleted_at,omitempty" bun:",soft_delete,nullzero"`
}
