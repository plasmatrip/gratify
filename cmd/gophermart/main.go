package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/plasmatrip/gratify/internal/api/middleware/logger"
	"github.com/plasmatrip/gratify/internal/config"
	"github.com/plasmatrip/gratify/internal/controller"
	"github.com/plasmatrip/gratify/internal/deps"
	"github.com/plasmatrip/gratify/internal/repository"
	"github.com/plasmatrip/gratify/internal/router"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	c, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	l, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer l.Close()

	db, err := repository.NewRepository(ctx, c.Database, *l)
	if err != nil {
		l.Sugar.Infow("database connection error: ", err)
		os.Exit(1)
	}
	defer db.Close()

	deps := &deps.Dependencies{
		Config: *c,
		Logger: *l,
		Repo:   *db,
	}

	ctrl := controller.NewController(c.ClientTimeout, *deps)
	ctrl.StartWorkers(ctx)
	ctrl.StartOrdersProcessor(ctx)

	server := http.Server{
		Addr: c.Host,
		Handler: func(next http.Handler) http.Handler {
			l.Sugar.Infow("The loyalty system \"Gophermart\" server is running. ", "Server address", c.Host)
			l.Sugar.Infow("Server config", "Accrual system address", c.Accrual)
			return next
		}(router.NewRouter(*deps, ctrl)),
	}

	go server.ListenAndServe()

	<-ctx.Done()

	server.Shutdown(context.Background())

	l.Sugar.Infow("The server has been shut down gracefully")

	os.Exit(0)
}
