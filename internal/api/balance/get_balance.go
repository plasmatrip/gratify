package balance

import (
	"encoding/json"
	"net/http"

	"github.com/plasmatrip/gratify/internal/api/auth"
)

func (b *Balance) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.ValidLogin{}).(*auth.Claims).UserdID

	balance, err := b.deps.Repo.GetBalanceWithdrawn(r.Context(), userID)
	if err != nil {
		b.deps.Logger.Sugar.Infow("error receiving data", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	preparedBalance, err := json.Marshal(balance)
	if err != nil {
		b.deps.Logger.Sugar.Infow("serialization error", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(preparedBalance)
	if err != nil {
		b.deps.Logger.Sugar.Infow("data write error", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	b.deps.Logger.Sugar.Infow("response", "balance", balance)
}
