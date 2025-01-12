package repository

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/plasmatrip/gratify/internal/apperr"
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

func (r Repository) CheckLogin(ctx context.Context, userLogin models.LoginRequest) error {
	var user models.LoginRequest

	row := r.db.QueryRow(ctx, schema.SelectUser, pgx.NamedArgs{"login": userLogin.Login})

	err := row.Scan(&user.Login, &user.Password)
	if err != nil {
		return err
	}

	savedHash, err := hex.DecodeString(user.Password)
	if err != nil {
		return err
	}

	h := sha256.New()
	h.Write([]byte([]byte(userLogin.Password)))
	hash := h.Sum(nil)

	if user.Login != userLogin.Login || !bytes.Equal(hash, savedHash) {
		return apperr.ErrBadLogin
	}

	return nil
}

func (r Repository) RegisterUser(ctx context.Context, userLogin models.LoginRequest) error {
	var user models.LoginRequest
	row := r.db.QueryRow(ctx, schema.SelectUser, pgx.NamedArgs{"login": userLogin.Login})

	err := row.Scan(&user.Login, &user.Password)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}

	if len(user.Login) > 0 {
		return apperr.ErrLoginAlreadyTaken
	}

	h := sha256.New()
	h.Write([]byte([]byte(userLogin.Password)))
	hash := hex.EncodeToString(h.Sum(nil))

	res, err := r.db.Exec(ctx, "INSERT INTO users (login, password) VALUES (@login, @password)",
		pgx.NamedArgs{
			"login":    userLogin.Login,
			"password": hash,
		})
	if err != nil {
		return err
	}

	rows := res.RowsAffected()
	if rows == 0 {
		return apperr.ErrZeroRowInsert
	}

	return nil
}
