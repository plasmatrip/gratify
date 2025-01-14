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
	err = r.createTables(ctx)
	if err != nil {
		return nil, err
	}

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

	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
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

func (r Repository) RegisterUser(ctx context.Context, userLogin models.LoginRequest) (int32, error) {
	h := sha256.New()
	h.Write([]byte([]byte(userLogin.Password)))
	hash := hex.EncodeToString(h.Sum(nil))

	var id int32

	err := r.db.QueryRow(ctx, schema.InsertUser, pgx.NamedArgs{
		"login":    userLogin.Login,
		"password": hash,
	}).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (r Repository) AddOrder(ctx context.Context, order models.Order) error {
	rows, err := r.db.Query(ctx, schema.SelectOrderFromAnotherUser, pgx.NamedArgs{
		"id":      order.Number,
		"user_id": order.UserID,
	})
	if err != nil {
		return err
	}
	rows.Close()

	if rows.CommandTag().RowsAffected() > 0 {
		return apperr.ErrOrderAlreadyUploadedAnotherUser
	}

	_, err = r.db.Exec(ctx, schema.InsertOrder, pgx.NamedArgs{
		"id":      order.Number,
		"user_id": order.UserID,
		"status":  order.Status,
		"date":    order.Date,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) GetOrders(ctx context.Context, userId int32) ([]models.Order, error) {
	orders := []models.Order{}

	rows, err := r.db.Query(ctx, schema.SelectOrders, pgx.NamedArgs{
		"user_id": userId,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := models.Order{}

		err := rows.Scan(
			&order.Number,
			&order.UserID,
			&order.Status,
			&order.Accrual,
			&order.Sum,
			&order.Date,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
