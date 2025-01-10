package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/plasmatrip/gratify/internal/models"
)

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var userLogin models.UserLogin

	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		a.deps.Logger.Sugar.Infow("error in request handler", "error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(userLogin.Login) == 0 || len(userLogin.Password) == 0 {
		a.deps.Logger.Sugar.Infow("error in authentication data", "error: ", errors.New("empty login or password"))
		http.Error(w, "empty login or password", http.StatusUnauthorized)
		return
	}

	if err := a.deps.Repo.CheckLogin(r.Context(), userLogin); err != nil {
		if errors.Is(err, errors.ErrUnsupported) {
			a.deps.Logger.Sugar.Infow("authentication error", "error: ", err)
			http.Error(w, "empty login or password", http.StatusUnauthorized)
			return
		}

		a.deps.Logger.Sugar.Infow("internal error", "error: ", err)
		http.Error(w, "empty login or password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
