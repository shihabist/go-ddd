package api

import (
	"orderContext/application/behaviour"
	"orderContext/application/query"

	"orderContext/application/command"
	"orderContext/core/mediator"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type orderCommandHandler struct {
	mediator mediator.Mediator
}

type orderQueryHandler struct {
	orderservice query.OrderQueryService
}

func newOrderCommandHandler() orderCommandHandler {
	m := mediator.NewMediator().
		UseBehaviour(behaviour.NewLogger()).
		UseBehaviour(behaviour.NewValidator()).
		RegisterHandler(command.NewCreateOrderCommandHandler()).
		RegisterHandler(command.NewPayOrderCommandHandler())

	return orderCommandHandler{
		mediator: m,
	}
}

func newOrderQueryHandler() orderQueryHandler {
	return orderQueryHandler{
		orderservice: query.NewOrderQueryService(),
	}
}

// CreateOrder godoc
// @Summary Create a order
// @Description Create a new order
// @Tags order
// @Accept json
// @Produce json
// @Success 201 {object} string
// @Router /order [post]
func (o *orderCommandHandler) create(c echo.Context) error {
	return create(c, func() { o.mediator.Send(c.Request().Context(), command.CreateOrderCommand{Id: uuid.New().String()}) })
}

// PayOrder godoc
// @Summary Pay order
// @Description Pay the order
// @Tags order
// @Accept json
// @Produce json
// @Success 202 {object} string
// @Param id path string true "id"
// @Router /order/pay/{id} [put]
func (o *orderCommandHandler) pay(c echo.Context) error {
	return updateErr(c, func(id string) error {
		return o.mediator.Send(c.Request().Context(), command.PayOrderCommand{OrderId: id})
	})
}

// CancelOrder godoc
// @Summary Cancel order
// @Description Cancel the order
// @Tags order
// @Accept json
// @Produce json
// @Success 202 {object} string
// @Param id path string true "id"
// @Router /order/cancel/{id} [put]
func (o *orderCommandHandler) cancel(c echo.Context) error {
	return updateErr(c, func(id string) error {
		return o.mediator.Send(c.Request().Context(), command.CancelOrderCommand{OrderId: id})
	})
}

// func (o *orderHandler) ship(c echo.Context) error {
// 	return updateErr(c, func(id string) error { return o.orderservice.Ship(id) })
// }

// GetOrder godoc
// @Summary Get orders
// @Description Get all orders
// @Tags order
// @Accept json
// @Produce json
// @Success 200 {object} order.Order
// @Router /order [get]
func (o *orderQueryHandler) getOrders(c echo.Context) error {
	return get(c, o.orderservice.GetOrders())
}

// GetOrder godoc
// @Summary Get order
// @Description Get order
// @Tags order
// @Accept json
// @Produce json
// @Success 200 {object} order.Order
// @Router /order/:id [get]
func (o *orderQueryHandler) getOrder(c echo.Context) error {
	return get(c, func(id string) interface{} { return o.orderservice.GetOrder(id) })
}