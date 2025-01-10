package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/plasmatrip/gratify/internal/api"
	"github.com/plasmatrip/gratify/internal/api/auth"
	"github.com/plasmatrip/gratify/internal/api/balance"
	"github.com/plasmatrip/gratify/internal/api/info"
	"github.com/plasmatrip/gratify/internal/api/orders"
)

func NewRouter(deps api.Dependencies) *chi.Mux {

	r := chi.NewRouter()

	auth := auth.NewAuthService(deps)
	balance := balance.NewBalanceService(deps)
	orders := orders.NewOrdersService(deps)
	info := info.NewInfoService(deps)

	r.Use(deps.Logger.WithLogging)

	r.Route("/api/user/register", func(r chi.Router) {
		r.Post("/", auth.Register)
	})

	r.Route("/api/user/login", func(r chi.Router) {
		r.Post("/", auth.Login)
	})

	r.Route("/api/user/orders", func(r chi.Router) {
		r.Post("/", orders.AddOrders)
		r.Get("/", orders.GetOrders)
	})

	r.Route("/api/user/balance", func(r chi.Router) {
		r.Get("/", balance.GetBalance)
	})

	r.Route("/api/user/withdraw", func(r chi.Router) {
		r.Post("/", balance.Withdraw)
	})

	r.Route("/api/user/withdrawals", func(r chi.Router) {
		r.Get("/", balance.Withdrawals)
	})

	r.Route("/api/info", func(r chi.Router) {
		r.Get("/", info.Ping)
	})

	return r
}
