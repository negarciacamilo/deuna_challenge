package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/negarciacamilo/deuna_challenge/application/apierrors"
	"github.com/negarciacamilo/deuna_challenge/application/domain"
	dbd "github.com/negarciacamilo/deuna_challenge/application/domain/database"
	"github.com/negarciacamilo/deuna_challenge/application/logger"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

const (
	// Operations
	Fetching = "fetching"
	Creating = "creating"
	Updating = "updating"
	Deleting = "deleting"
)

type Database interface {
	GetDB() *bun.DB
	GetMock() sqlmock.Sqlmock
	HandleDBError(ctx *domain.ContextInformation, table, operation string, err error) apierrors.ApiError
}

type database struct {
	db   *bun.DB
	mock sqlmock.Sqlmock
}

func New(forTest ...bool) Database {
	if forTest != nil && forTest[0] {
		mockdb, mock, _ := sqlmock.New()
		db := bun.NewDB(mockdb, pgdialect.New())
		return &database{db: db, mock: mock}
	}

	sqlb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(viper.GetString("PAYMENTS_DSN"))))
	db := bun.NewDB(sqlb, pgdialect.New())

	if err := db.Ping(); err != nil {
		logger.Panic("can't ping DB", "new-db", err, nil)
	}

	createTables(db)

	return &database{db: db, mock: nil}
}

func (d *database) GetDB() *bun.DB {
	return d.db
}

func (d *database) GetMock() sqlmock.Sqlmock {
	return d.mock
}

func (d *database) HandleDBError(ctx *domain.ContextInformation, table, operation string, err error) apierrors.ApiError {
	if errors.Is(err, sql.ErrNoRows) {
		apierr := apierrors.NewNotFoundApiError(fmt.Sprintf("error, %s not found", table))
		logger.Error(apierr.Message(), logger.GetCallerFunctionName(), err, ctx)
		return apierr
	} else {
		apierr := apierrors.NewInternalServerApiError(fmt.Sprintf("error %s %s", operation, table), err)
		logger.Error(apierr.Message(), logger.GetCallerFunctionName(), err, ctx)
		return apierr
	}
}

func createTables(db *bun.DB) {
	models := []interface{}{
		(*dbd.Bank)(nil),
		(*dbd.Customer)(nil),
		(*dbd.Merchant)(nil),
		(*dbd.Payment)(nil),
	}

	for _, model := range models {
		if _, err := db.NewCreateTable().IfNotExists().Model(model).Exec(context.Background()); err != nil {
			logger.Panic(fmt.Sprintf("can't create table %s", model), "create-tables", err, nil)
		}
	}
}
