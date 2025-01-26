package balance

import (
	"encoding/json"
	"net/http"

	"github.com/plasmatrip/gratify/internal/api/auth"
)

func (b *Balance) Withdrawals(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.ValidLogin{}).(*auth.Claims).UserdID

	foundWithdrawals, err := b.deps.Repo.Withdrawals(r.Context(), userID)
	if err != nil {
		b.deps.Logger.Sugar.Infow("internal server error", "error: ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	b.deps.Logger.Sugar.Infow("response received from db", "withdrawals", foundWithdrawals)

	if len(foundWithdrawals) == 0 {
		b.deps.Logger.Sugar.Infow("no data", "error: ", err)
		http.Error(w, "no data", http.StatusNoContent)
		return
	}

	preparedOrders, err := json.Marshal(foundWithdrawals)
	if err != nil {
		b.deps.Logger.Sugar.Infow("serialization error", "error: ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(preparedOrders)
	if err != nil {
		b.deps.Logger.Sugar.Infow("data write error", "error: ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
