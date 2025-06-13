package v1

import (
	"errors"

	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	. "github.com/krijebr/printer-shop/internal/delivery/http/common"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type OrderHandlers struct {
	usecase usecase.Order
}

func NewOrderHandlers(u usecase.Order) *OrderHandlers {
	return &OrderHandlers{usecase: u}
}

func (o *OrderHandlers) getAllOrders() echo.HandlerFunc {
	return func(c echo.Context) error {
		var filter *entity.OrderFilter
		filter = nil
		if c.QueryParam("user_id") != "" || c.QueryParam("order_status") != "" {
			filter = new(entity.OrderFilter)
		}
		if c.QueryParam("user_id") != "" {
			userId, err := uuid.Parse(c.QueryParam("user_id"))
			if err != nil {
				slog.Error("invalid user id", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrValidationErrorCode,
					Message: ErrValidationErrorMessage,
				})
			}
			userRole, ok := c.Get(UserRoleContextKey).(entity.UserRole)
			if !ok {
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
			if userRole != entity.UserRoleAdmin {
				userIdCtx, ok := c.Get(UserIdContextKey).(uuid.UUID)
				if !ok {
					return c.JSON(http.StatusInternalServerError, ErrResponse{
						Error:   ErrInternalErrorCode,
						Message: ErrInternalErrorMessage,
					})
				}
				if userIdCtx != userId {
					return c.JSON(http.StatusForbidden, ErrResponse{
						Error:   ErrForbiddenCode,
						Message: ErrForbiddenMessage,
					})
				}
			}
			filter.UserId = &userId
		}
		if c.QueryParam("order_status") != "" {
			validate := validator.New()
			err := validate.Var(c.QueryParam("order_status"), "oneof=new in_progress done")
			if err != nil {
				slog.Error("validation error", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrValidationErrorCode,
					Message: ErrValidationErrorMessage,
				})
			}
			orderStatus := entity.OrderStatus(c.QueryParam("order_status"))
			filter.Status = &orderStatus
		}

		orders, err := o.usecase.GetAll(c.Request().Context(), filter)
		if err != nil {
			slog.Error("orders receiving error", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		slog.Info("all orders received")
		return c.JSON(http.StatusOK, orders)
	}
}

func (o *OrderHandlers) createOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, ok := c.Get(UserIdContextKey).(uuid.UUID)
		if !ok {
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		order, err := o.usecase.Create(c.Request().Context(), userId)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrCartIsEmpty):
				slog.Error("producer not found", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrCartIsEmptyCode,
					Message: ErrCartIsEmptyMessage,
				})
			default:
				slog.Error("order creation error", slog.Any("error", err))
				return c.NoContent(http.StatusInternalServerError)
			}
		}
		return c.JSON(http.StatusOK, order)
	}
}

func (o *OrderHandlers) getOrderById() echo.HandlerFunc {
	return func(c echo.Context) error {
		orderId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			slog.Error("invalid order id", slog.Any("error", err))
			return c.JSON(http.StatusNotFound, ErrResponse{
				Error:   ErrResourceNotFoundCode,
				Message: ErrResourceNotFoundMessage,
			})
		}
		order, err := o.usecase.GetById(c.Request().Context(), orderId)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrOrderNotFound):
				slog.Error("order not found", slog.Any("error", err))
				return c.JSON(http.StatusNotFound, ErrResponse{
					Error:   ErrResourceNotFoundCode,
					Message: ErrResourceNotFoundMessage,
				})
			default:
				slog.Error("order receiving error", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		userRole, ok := c.Get(UserRoleContextKey).(entity.UserRole)
		if !ok {
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		if userRole != entity.UserRoleAdmin {
			userId, ok := c.Get(UserIdContextKey).(uuid.UUID)
			if !ok {
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
			if userId != order.UserId {
				return c.JSON(http.StatusForbidden, ErrResponse{
					Error:   ErrForbiddenCode,
					Message: ErrForbiddenMessage,
				})
			}
		}
		if !ok {
			return c.JSON(http.StatusInternalServerError, ErrResponse{
				Error:   ErrInternalErrorCode,
				Message: ErrInternalErrorMessage,
			})
		}
		slog.Info("order received")
		return c.JSON(http.StatusOK, order)
	}
}
func (o *OrderHandlers) updateOrderById() echo.HandlerFunc {
	type (
		Product struct {
			Id    uuid.UUID `json:"id"`
			Count int       `json:"count" validate:"gt=0"`
		}
		request struct {
			Status   entity.OrderStatus `json:"status" validate:"omitempty,oneof=new in_progress done"`
			Products []*Product         `json:"products"`
		}
	)
	return func(c echo.Context) error {
		orderId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			slog.Error("invalid product id", slog.Any("error", err))
			return c.JSON(http.StatusNotFound, ErrResponse{
				Error:   ErrResourceNotFoundCode,
				Message: ErrResourceNotFoundMessage,
			})
		}
		var requestData request
		err = c.Bind(&requestData)
		if err != nil {
			slog.Error("invalid request", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Error:   ErrInvalidRequestCode,
				Message: ErrInvalidRequestMessage,
			})
		}
		validate := validator.New()
		err = validate.Struct(requestData)
		if err != nil {
			slog.Error("validation error", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, ErrResponse{
				Error:   ErrValidationErrorCode,
				Message: ErrValidationErrorMessage,
			})
		}
		order := &entity.Order{}
		order.Id = orderId
		if requestData.Status != "" {
			order.Status = requestData.Status
		}
		if len(requestData.Products) == 0 {
			requestData.Products = nil
		}
		if requestData.Products != nil {
			productsMap := make(map[uuid.UUID]int)
			for _, productInRequest := range requestData.Products {
				err = validate.Struct(productInRequest)
				if err != nil {
					slog.Error("validation error", slog.Any("error", err))
					return c.JSON(http.StatusBadRequest, ErrResponse{
						Error:   ErrValidationErrorCode,
						Message: ErrValidationErrorMessage,
					})
				}
				if count, inMap := productsMap[productInRequest.Id]; inMap {
					productsMap[productInRequest.Id] = count + productInRequest.Count
				} else {
					productsMap[productInRequest.Id] = productInRequest.Count
				}
			}
			for id, count := range productsMap {
				product := &entity.ProductInCart{
					Product: &entity.Product{
						Id: id,
					},
					Count: count,
				}
				order.Products = append(order.Products, product)
			}
		}

		updatedOrder, err := o.usecase.UpdateById(c.Request().Context(), order)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrOrderNotFound):
				slog.Error("order not found", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrOrderNotExistCode,
					Message: ErrOrderNotExistMessage,
				})
			case errors.Is(err, usecase.ErrProductNotFound):
				slog.Error("product not found", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrProductNotExistCode,
					Message: ErrProductNotExistMessage,
				})
			case errors.Is(err, usecase.ErrOrderCantBeUpdated):
				slog.Error("order can't be updated", slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, ErrResponse{
					Error:   ErrOrderCantBeUpdatedCode,
					Message: ErrOrderCantBeUpdatedMessage,
				})
			default:
				slog.Error("order updating error", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		return c.JSON(http.StatusOK, updatedOrder)
	}
}
func (o *OrderHandlers) deleteOrderById() echo.HandlerFunc {
	return func(c echo.Context) error {
		orderId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			slog.Error("invalid order id", slog.Any("error", err))
			return c.JSON(http.StatusNotFound, ErrResponse{
				Error:   ErrResourceNotFoundCode,
				Message: ErrResourceNotFoundMessage,
			})
		}
		err = o.usecase.DeleteById(c.Request().Context(), orderId)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrOrderNotFound):
				slog.Error("order not found", slog.Any("error", err))
				return c.JSON(http.StatusNotFound, ErrResponse{
					Error:   ErrResourceNotFoundCode,
					Message: ErrResourceNotFoundMessage,
				})
			default:
				slog.Error("order receiving error", slog.Any("error", err))
				return c.JSON(http.StatusInternalServerError, ErrResponse{
					Error:   ErrInternalErrorCode,
					Message: ErrInternalErrorMessage,
				})
			}
		}
		slog.Info("product received")
		return c.NoContent(http.StatusOK)
	}
}
func RegisterOrderRoutes(u usecase.Order, g *echo.Group) {
	a := NewOrderHandlers(u)
	g.GET("", a.getAllOrders())
	g.POST("", a.createOrder())
	g.GET("/:id", a.getOrderById())
	g.PUT("/:id", a.updateOrderById())
	g.DELETE("/:id", a.deleteOrderById())
}
