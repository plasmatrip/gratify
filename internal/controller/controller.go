package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/plasmatrip/gratify/internal/deps"
	"github.com/plasmatrip/gratify/internal/models"
)

type Result struct {
	order models.Order
	err   error
}

type Controler struct {
	works   chan models.Order
	results chan Result
	client  http.Client
	deps    deps.Dependencies
	wg      sync.WaitGroup
}

func NewConntroller(timeout time.Duration, deps deps.Dependencies) *Controler {
	return &Controler{
		works:   make(chan models.Order),
		results: make(chan Result),
		client:  http.Client{Timeout: deps.Config.ClientTimeout},
		deps:    deps,
		wg:      sync.WaitGroup{},
	}
}

func (c *Controler) StartWorkers(ctx context.Context) {
	for i := 0; i < c.deps.Config.Workers; i++ {
		go c.Worker(ctx, i)

	}

	go func() {

		for {
			select {
			case result := <-c.results:
				if result.err != nil {
					c.deps.Logger.Sugar.Infow("error interacting with the accrual service", "error", result.err)
				}

				err := c.deps.Repo.UpdateOrder(ctx, result.order)
				if err != nil {
					c.deps.Logger.Sugar.Infow("loyalty accumulation update error", "error", err)
				}
			case <-ctx.Done():
				c.wg.Wait()
				return
			}
		}

	}()
}

func (c *Controler) Worker(ctx context.Context, idx int) {
	c.wg.Add(1)
	defer c.wg.Done()
	for {
		select {
		case work := <-c.works:
			c.deps.Logger.Sugar.Infow("start work", "worker", idx)
			order, err := c.AccrualProcess(work)
			result := Result{
				order: order,
				err:   err,
			}
			c.results <- result
			c.deps.Logger.Sugar.Infow("stop work", "worker", idx)
		case <-ctx.Done():
			c.deps.Logger.Sugar.Infow("done work", "worker", idx)
			return
		}
	}

}

func (c *Controler) AddWork(order models.Order) {
	c.deps.Logger.Sugar.Infoln("add work")
	c.works <- order
}

func (c *Controler) AccrualProcess(order models.Order) (models.Order, error) {
	result := models.Order{}

	req, err := http.NewRequest(http.MethodGet, c.deps.Config.Accrual+"/api/orders/"+strconv.FormatInt(order.Number, 10), nil)
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
	c.deps.Logger.Sugar.Infow("result", "status", result.Status)

	return result, nil
}
