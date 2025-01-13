package orders

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/plasmatrip/gratify/internal/api/auth"
	"github.com/plasmatrip/gratify/internal/models"
)

func (o *Orders) GetOrders(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.ValidLogin{}).(*auth.Claims).UserdID

	orders, err := o.deps.Repo.GetOrders(r.Context(), userId)
	if err != nil {
		o.deps.Logger.Sugar.Infow("internal server error", "error: ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		o.deps.Logger.Sugar.Infow("no data", "error: ", err)
		http.Error(w, "no data", http.StatusNoContent)
		return
	}

	prepareOrder := make([]models.ResponseOrder, len(orders))
	for i, order := range orders {
		prepareOrder[i] = models.ResponseOrder{
			Number:  order.Number,
			Status:  order.Status.String(),
			Accrual: order.Accrual,
			Date:    order.Date.Format(time.RFC3339),
		}
	}

	foundOrders, err := json.Marshal(prepareOrder)
	if err != nil {
		o.deps.Logger.Sugar.Infow("serialization error", "error: ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(foundOrders)
	if err != nil {
		o.deps.Logger.Sugar.Infow("data write error", "error: ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
