package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/plasmatrip/gratify/internal/apperr"
	"github.com/plasmatrip/gratify/internal/deps"
	"github.com/plasmatrip/gratify/internal/models"
)

type Result struct {
	order models.Order
	err   error
}

type Controller struct {
	works   chan models.Order
	results chan Result
	wait    chan struct{}
	client  http.Client
	deps    deps.Dependencies
	wg      sync.WaitGroup
	mu      sync.Mutex
	cond    *sync.Cond
}

// NewController returns a new instance of Controller.
//
// The Controller is responsible for processing orders in parallel. The number
// of workers is determined by the Workers field in the deps.Config.
//
// The timeout is the maximum time the controller will wait for a response from
// the accrual service when processing orders.
//
// The deps is a struct containing the dependencies required by the controller,
// such as the repository and logger.
func NewController(timeout time.Duration, deps deps.Dependencies) *Controller {
	ctrl := &Controller{
		works:   make(chan models.Order, deps.Config.WorkBuffer),
		results: make(chan Result),
		wait:    make(chan struct{}, deps.Config.Workers),
		client:  http.Client{Timeout: deps.Config.ClientTimeout},
		deps:    deps,
		wg:      sync.WaitGroup{},
		mu:      sync.Mutex{},
	}
	ctrl.cond = sync.NewCond(&ctrl.mu)
	return ctrl
}

// StartOrdersProcessor starts a goroutine that periodically gets unprocessed orders
// from the database and adds them to the works channel. When the context is
// canceled, the goroutine stops and the works channel is closed.
func (c *Controller) StartOrdersProcessor(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(c.deps.Config.ProcessorInterval) * time.Second)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				orders, err := c.deps.Repo.GetUnprocessedOrders(ctx)
				if err != nil {
					c.deps.Logger.Sugar.Infow("error receiving unprocessed orders from database", "error", err)
					continue
				}
				for _, order := range orders {
					c.works <- order
				}
			case <-ctx.Done():
				close(c.works)
				return
			}
		}
	}()
}

// StartWorkers initializes and starts a number of worker goroutines as defined
// by the Workers field in the configuration. Each worker listens for tasks on
// the works channel and processes them using the AccrualProcessor function.
// The results of processing are sent to the results channel. This method also
// starts a goroutine to listen for results, logging any errors and updating
// the order status in the repository. The function blocks until all workers
// have been started, and it waits for the context to be canceled before
// shutting down gracefully.
func (c *Controller) StartWorkers(ctx context.Context) {
	c.wg.Add(c.deps.Config.Workers)

	wg := &sync.WaitGroup{}
	wg.Add(c.deps.Config.Workers)

	for i := 0; i < c.deps.Config.Workers; i++ {
		go c.Worker(ctx, i, wg)
	}

	wg.Wait()

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

// Worker is a goroutine that processes tasks from the works channel. Each worker
// listens for tasks and processes them using the AccrualProcessor function. The
// results are sent to the results channel. If the worker receives a signal on the
// wait channel, it pauses processing until it is signaled to continue. The worker
// stops processing tasks and exits when the context is canceled.
func (c *Controller) Worker(ctx context.Context, idx int, wg *sync.WaitGroup) {
	c.deps.Logger.Sugar.Infow("worker started", "worker index", idx)

	wg.Done()

	defer c.wg.Done()

	for {
		select {
		case work := <-c.works:
			select {
			case <-c.wait:
				c.cond.L.Lock()
				c.deps.Logger.Sugar.Infoln("the worker paused", "worker index", idx)
				c.cond.Wait()
				c.deps.Logger.Sugar.Infoln("the worker unpaused", "worker index", idx)
				c.cond.L.Unlock()
			default:
			}

			c.deps.Logger.Sugar.Infow("the worker started performing the task", "worker index", idx, "task", work)

			order, err := c.AccrualProcessor(work)
			c.results <- Result{
				order: order,
				err:   err,
			}

			c.deps.Logger.Sugar.Infow("the worker completed the task", "worker index", idx)
		case <-ctx.Done():
			c.deps.Logger.Sugar.Infow("the worker is stopped", "worker index", idx)
			return
		}
	}
}

// AccrualProcessor is a function that sends a request to the accrual service with the order number and
// processes the response. The function returns the processed order and error. If the order is not
// registered in the accrual system, the function returns ErrOrderIsNotRegisteredInAccrual. If the
// accrual service returns an error, the function returns ErrInternalServerAccrualError. If the
// accrual service returns a retry-after header, the function stops all workers until the time
// specified in the header.
func (c *Controller) AccrualProcessor(order models.Order) (models.Order, error) {
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

	switch resp.StatusCode {
	case http.StatusOK:
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

	case http.StatusNoContent:
		return result, apperr.ErrOrderIsNotRegisteredInAccrual

	case http.StatusTooManyRequests:
		retryAfter := resp.Header.Get("Retry-After")

		retry, err := strconv.Atoi(retryAfter)
		if err != nil {
			c.deps.Logger.Sugar.Infow("invalid timeout format", "error: ", err)
			return result, err
		}

		go func() {
			for i := 0; i < c.deps.Config.Workers; i++ {
				c.wait <- struct{}{}
			}
			timer := time.NewTimer(time.Duration(retry) * time.Second)
			<-timer.C
			c.cond.Broadcast()
		}()

	case http.StatusInternalServerError:
		return result, apperr.ErrInternalServerAccrualError
	}

	c.deps.Logger.Sugar.Infow("response received from accrual service", "result", result, "status", result.Status.String())

	return result, nil
}
