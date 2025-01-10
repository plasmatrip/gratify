package repository

import (
	"bytes"
	"context"
	"crypto/sha256"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/plasmatrip/gratify/internal/api/errors"
	"github.com/plasmatrip/gratify/internal/logger"
	"github.com/plasmatrip/gratify/internal/models"
	"github.com/plasmatrip/gratify/internal/repository/schema"
)

type Repository struct {
	db *pgxpool.Pool
	l  logger.Logger
}

func NewRepository(ctx context.Context, dsn string, l logger.Logger) (*Repository, error) {
	// открываем БД
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	r := &Repository{
		db: db,
		l:  l,
	}

	// создаем таблицу, при ошибке прокидываем ее наверх
	// err = r.createTables(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	return r, nil
}

func (r Repository) createTables(ctx context.Context) error {
	_, err := r.db.Exec(ctx, schema.DBSchema)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) Ping(ctx context.Context) error {
	return r.db.Ping(ctx)
}

func (r Repository) Close() {
	r.db.Close()
}

func (r Repository) CheckLogin(ctx context.Context, userLogin models.UserLogin) error {
	var user models.UserLogin

	row := r.db.QueryRow(ctx, "SELECT * FROM users WHERE login = @login", pgx.NamedArgs{"login": userLogin.Login})

	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return err
	}

	h := sha256.New()
	h.Write([]byte([]byte(user.Password)))
	hash := h.Sum(nil)

	if user.Login != userLogin.Login || bytes.Equal(hash, []byte(user.Password)) {
		return errors.ErrBadLogin
	}

	return nil
}
