package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/plasmatrip/gratify/internal/config"
	"github.com/plasmatrip/gratify/internal/logger"
	"github.com/plasmatrip/gratify/internal/models"
	"github.com/plasmatrip/gratify/internal/repository"
)

type Result struct {
	order models.Order
	err   error
}

type Controler struct {
	works   chan models.Order
	results chan Result
	client  http.Client
	cfg     config.Config
	log     logger.Logger
	db      repository.Repository
}

func NewConntroller(timeout time.Duration, cfg config.Config, log logger.Logger, db repository.Repository) *Controler {
	return &Controler{
		works:   make(chan models.Order),
		results: make(chan Result),
		client:  http.Client{Timeout: cfg.ClientTimeout},
		cfg:     cfg,
		log:     log,
		db:      db,
	}
}

func (c Controler) StartWorkers(ctx context.Context) {
	for i := 1; i < c.cfg.Workers; i++ {
		go c.Worker(ctx, i)
	}

	go func() {
		for {
			select {
			case result := <-c.results:
				if result.err != nil {
					c.log.Sugar.Infow("error interacting with the accrual service", "error", result.err)
				}

				err := c.db.UpdateOrder(ctx, result.order)
				if err != nil {
					c.log.Sugar.Infow("loyalty accumulation update error", "error", result.err)
				}

			case <-ctx.Done():
				return
			}
		}
	}()

}

func (c Controler) Worker(ctx context.Context, idx int) {
	select {
	case work := <-c.works:
		// fmt.Printf("worker %d start\n", idx)
		order, err := c.AccrualProcess(work)
		result := Result{
			order: order,
			err:   err,
		}
		c.results <- result
		// fmt.Printf("worker %d end\n", idx)
	case <-ctx.Done():
		// fmt.Printf("worker %d stop\n", idx)
		return
	}
}

func (c Controler) AddWork(order models.Order) {
	c.works <- order
}

func (c Controler) AccrualProcess(order models.Order) (models.Order, error) {
	result := models.Order{}

	req, err := http.NewRequest(http.MethodGet, c.cfg.Accrual+"/api/orders/"+strconv.FormatInt(order.Number, 10), nil)
	if err != nil {
		return result, err
	}

	req.Header.Set("Content-Lenght", "0")

	resp, err := c.client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-type")

	if strings.Contains(contentType, "application/json") {
		var accrualResponse models.AccrualResponse
		if err := json.NewDecoder(resp.Body).Decode(&accrualResponse); err != nil {
			return result, err
		}

		orderID, err := strconv.ParseInt(accrualResponse.Order, 10, 64)
		if err != nil {
			return result, err
		}

		result.UserID = order.UserID
		result.Number = orderID
		result.Accrual = accrualResponse.Accrual
		result.Status.Scan(accrualResponse.Status)
	}

	if strings.Contains(contentType, "text/plain") {
		retryAfter := resp.Header.Get("Retry-After")
		print(retryAfter)
	}

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }

	return result, nil
}
