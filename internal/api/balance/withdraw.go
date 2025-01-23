package balance

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/plasmatrip/gratify/internal/api/auth"
	"github.com/plasmatrip/gratify/internal/models"
)

func (b *Balance) Withdraw(w http.ResponseWriter, r *http.Request) {
	var withdraw models.Withdraw

	if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
		b.deps.Logger.Sugar.Infow("error in request handler", "error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := goluhn.Validate(withdraw.Order)
	if err != nil {
		b.deps.Logger.Sugar.Infow("invalid order number format", "error: ", err)
		http.Error(w, "invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	userID := r.Context().Value(auth.ValidLogin{}).(*auth.Claims).UserdID

	balance, err := b.deps.Repo.GetBalance(r.Context(), userID)
	if err != nil {
		b.deps.Logger.Sugar.Infow("error receiving data", "error: ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if balance < withdraw.Sum {
		b.deps.Logger.Sugar.Infow("user account has insufficient funds", "user: ", userID, "sum", withdraw.Sum)
		http.Error(w, "there are insufficient funds in your account to complete this transaction", http.StatusPaymentRequired)
		return
	}

	orderID, err := strconv.ParseInt(withdraw.Order, 10, 64)
	if err != nil {
		b.deps.Logger.Sugar.Infow("invalid order ID format", "error: ", err)
		http.Error(w, "invalid order ID format", http.StatusBadRequest)
		return
	}

	order := models.Order{
		Number: orderID,
		UserID: userID,
		Sum:    withdraw.Sum,
	}

	if err := b.deps.Repo.Withdraw(r.Context(), order); err != nil {
		b.deps.Logger.Sugar.Infow("error add ", "error: ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
