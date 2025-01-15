package orders

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/plasmatrip/gratify/internal/api/auth"
	"github.com/plasmatrip/gratify/internal/models"
)

func (o *Orders) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.ValidLogin{}).(*auth.Claims).UserdID

	foundOrders, err := o.deps.Repo.GetOrders(r.Context(), userID)
	if err != nil {
		o.deps.Logger.Sugar.Infow("internal server error", "error: ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if len(foundOrders) == 0 {
		o.deps.Logger.Sugar.Infow("no data", "error: ", err)
		http.Error(w, "no data", http.StatusNoContent)
		return
	}

	orders := make([]models.ResponseOrder, len(foundOrders))
	for i, order := range foundOrders {
		orders[i] = models.ResponseOrder{
			Number:  strconv.Itoa(int(order.Number)),
			Status:  order.Status.String(),
			Accrual: order.Accrual,
			Date:    order.Date.Format(time.RFC3339),
		}
	}

	preparedOrders, err := json.Marshal(orders)
	if err != nil {
		o.deps.Logger.Sugar.Infow("serialization error", "error: ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(preparedOrders)
	if err != nil {
		o.deps.Logger.Sugar.Infow("data write error", "error: ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
