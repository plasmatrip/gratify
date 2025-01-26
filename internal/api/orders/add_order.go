package orders

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/plasmatrip/gratify/internal/api/auth"
	"github.com/plasmatrip/gratify/internal/apperr"
	"github.com/plasmatrip/gratify/internal/models"
	"github.com/rgurov/pgerrors"
)

func (o *Orders) AddOrder(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		o.deps.Logger.Sugar.Infow("failed to read request body", "error: ", err)
		http.Error(w, "failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Преобразуем тело запроса в строку и обрезаем лишние пробелы
	orderIDStr := string(body)
	orderIDStr = strings.TrimSpace(orderIDStr)

	err = goluhn.Validate(orderIDStr)
	if err != nil {
		o.deps.Logger.Sugar.Infow("invalid order number format", "error: ", err)
		http.Error(w, "invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	// Преобразуем строку в число
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		o.deps.Logger.Sugar.Infow("invalid order ID format", "error: ", err)
		http.Error(w, "invalid order ID format", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(auth.ValidLogin{}).(*auth.Claims).UserdID

	order := models.Order{
		Number: orderID,
		UserID: userID,
		Status: models.StatusNew,
		Date:   time.Now(),
	}

	err = o.deps.Repo.AddOrder(r.Context(), order)
	if err != nil {
		if errors.Is(err, apperr.ErrOrderAlreadyUploadedAnotherUser) {
			o.deps.Logger.Sugar.Infow("error adding order", "error: ", err)
			http.Error(w, "error adding order", http.StatusConflict)
			return
		}

		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrors.UniqueViolation {
				o.deps.Logger.Sugar.Infow("error adding order", "error: ", err)
				http.Error(w, "error adding order", http.StatusOK)
				return
			}

			o.deps.Logger.Sugar.Infow("", "error: ", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		o.deps.Logger.Sugar.Infow("error adding order", "error: ", err)
		http.Error(w, "error adding order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
