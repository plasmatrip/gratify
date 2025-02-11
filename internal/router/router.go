package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/plasmatrip/gratify/internal/api/auth"
	"github.com/plasmatrip/gratify/internal/api/balance"
	"github.com/plasmatrip/gratify/internal/api/info"
	"github.com/plasmatrip/gratify/internal/api/middleware/compress"
	"github.com/plasmatrip/gratify/internal/api/orders"
	"github.com/plasmatrip/gratify/internal/controller"
	"github.com/plasmatrip/gratify/internal/deps"
)

func NewRouter(deps deps.Dependencies, controller *controller.Controller) *chi.Mux {

	r := chi.NewRouter()

	auth := auth.NewAuthService(deps)
	balance := balance.NewBalanceService(deps)
	orders := orders.NewOrdersService(deps)
	info := info.NewInfoService(deps)

	r.Use(deps.Logger.WithLogging)
	r.Use(compress.WithCompressed)

	r.Route("/api/user/register", func(r chi.Router) {
		r.Post("/", auth.Register)
	})

	r.Route("/api/user/login", func(r chi.Router) {
		r.Post("/", auth.Login)
	})

	r.Route("/api/user/orders", func(r chi.Router) {
		r.Use(auth.Validate)
		r.Post("/", orders.AddOrder)
		r.Get("/", orders.GetOrders)
	})

	r.Route("/api/user/balance", func(r chi.Router) {
		r.Use(auth.Validate)
		r.Get("/", balance.GetBalance)
	})

	r.Route("/api/user/balance/withdraw", func(r chi.Router) {
		r.Use(auth.Validate)
		r.Post("/", balance.Withdraw)
	})

	r.Route("/api/user/withdrawals", func(r chi.Router) {
		r.Use(auth.Validate)
		r.Get("/", balance.Withdrawals)
	})

	r.Route("/api/info", func(r chi.Router) {
		r.Get("/", info.Ping)
	})

	return r
}
