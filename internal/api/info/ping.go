package info

import (
	"net/http"
)

func (i *Info) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	err := i.deps.Repo.Ping(r.Context())
	if err != nil {
		i.deps.Logger.Sugar.Infow(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
