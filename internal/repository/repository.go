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

	// // создаем таблицу, при ошибке прокидываем ее наверх
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

func (r Repository) CheckLogin(ctx context.Context, userLogin *models.LoginRequest) error {
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

	userLogin.ID = user.ID

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

func (r Repository) GetOrders(ctx context.Context, userID int32) ([]models.Order, error) {
	orders := []models.Order{}

	rows, err := r.db.Query(ctx, schema.SelectOrders, pgx.NamedArgs{
		"user_id": userID,
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

func (r Repository) GetBalanceWithdrawn(ctx context.Context, userID int32) (models.Balance, error) {
	var balance models.Balance

	row := r.db.QueryRow(ctx, schema.SelectUserBalanceWithdrawn, pgx.NamedArgs{
		"user_id": userID,
	})
	err := row.Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return models.Balance{}, err
		}
	}

	return balance, nil
}

func (r Repository) GetBalance(ctx context.Context, userID int32) (float32, error) {
	var balance float32

	row := r.db.QueryRow(ctx, schema.SelectUserBalance, pgx.NamedArgs{
		"user_id": userID,
	})
	err := row.Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (r Repository) Withdraw(ctx context.Context, order models.Order) error {
	// начинаем транзакцию
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})
	if err != nil {
		return err
	}

	// при ошибке коммита откатываем назад
	defer func() error {
		return tx.Rollback(ctx)
	}()

	_, err = tx.Exec(ctx, schema.InsertOrderWithdraw, pgx.NamedArgs{
		"id":      order.Number,
		"user_id": order.UserID,
		"sum":     order.Sum,
		"date":    order.Date,
	})
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, schema.UpdateBalanceWithdraw, pgx.NamedArgs{
		"user_id": order.UserID,
		"sum":     order.Sum,
	})
	if err != nil {
		return err
	}

	// запускаем коммит
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) Withdrawals(ctx context.Context, userID int32) ([]models.Withdraw, error) {
	withdrawals := []models.Withdraw{}

	rows, err := r.db.Query(ctx, schema.SelectWithdrawals, pgx.NamedArgs{
		"user_id": userID,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		withdraw := models.Withdraw{}

		err := rows.Scan(
			&withdraw.Order,
			&withdraw.Sum,
			&withdraw.Date,
		)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, withdraw)
	}

	return withdrawals, nil
}

func (r Repository) UpdateOrder(ctx context.Context, order models.Order) error {
	// обновляем статус заказа, если пришел Processed
	if order.Status == models.StatusProcessed {
		// начинаем транзакцию
		tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})
		if err != nil {
			return err
		}

		// при ошибке коммита откатываем назад
		defer func() error {
			return tx.Rollback(ctx)
		}()

		_, err = tx.Exec(ctx, schema.UpdateOrderAccrual, pgx.NamedArgs{
			"user_id": order.UserID,
			"accrual": order.Accrual,
			"status":  order.Status,
		})
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, schema.UpsertBalanceAccrual, pgx.NamedArgs{
			"user_id": order.UserID,
			"accrual": order.Accrual,
		})
		if err != nil {
			return err
		}

		// запускаем коммит
		err = tx.Commit(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	// обновляем статус заказа, если пришел не Processed
	_, err := r.db.Exec(ctx, schema.UpdateOrderStatus, pgx.NamedArgs{
		"user_id": order.UserID,
		"status":  order.Status,
	})
	if err != nil {
		return err
	}

	return nil
}
