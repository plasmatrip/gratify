package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/plasmatrip/gratify/internal/apperr"
	"github.com/plasmatrip/gratify/internal/models"
)

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var lr models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		a.deps.Logger.Sugar.Infow("error in request handler", "error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(lr.Login) == 0 || len(lr.Password) == 0 {
		a.deps.Logger.Sugar.Infow("error in authentication data", "error: ", errors.New("empty login or password"))
		http.Error(w, "empty login or password", http.StatusUnauthorized)
		return
	}

	if err := a.deps.Repo.CheckLogin(r.Context(), lr); err != nil {
		if errors.Is(err, apperr.ErrBadLogin) {
			a.deps.Logger.Sugar.Infow("authentication error", "error: ", err)
			http.Error(w, "authentication error", http.StatusUnauthorized)
			return
		}

		a.deps.Logger.Sugar.Infow("internal error", "error: ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	token, err := a.LoginToken(lr)
	if err != nil {
		a.deps.Logger.Sugar.Infow("error generating JWT", "error: ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}
