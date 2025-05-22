package v1

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/krijebr/printer-shop/internal/entity"
	"github.com/krijebr/printer-shop/internal/usecase"
	"github.com/labstack/echo/v4"
)

type ProductHandlers struct {
	usecase usecase.Product
}

func NewProductHandlers(u usecase.Product) *ProductHandlers {
	return &ProductHandlers{usecase: u}
}

func (p *ProductHandlers) getAllProducts() echo.HandlerFunc {
	return func(c echo.Context) error {
		var filter *entity.ProductFilter
		filter = nil
		if c.QueryParam("producer_id") != "" {
			filter = new(entity.ProductFilter)
			producerId, err := uuid.Parse(c.QueryParam("producer_id"))
			if err != nil {
				log.Println("Ошибка преобразования uuid", err)
				return c.NoContent(http.StatusBadRequest)
			}
			filter.ProducerId = &producerId
		}
		products, err := p.usecase.GetAll(c.Request().Context(), filter)
		if err != nil {
			log.Println("Ошибка получения товаров", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		log.Println("Получение всех товаров")
		c.Response().Header().Set(echo.HeaderContentType, "application/json")
		return c.JSON(http.StatusOK, products)
	}
}

func (p *ProductHandlers) createProduct() echo.HandlerFunc {
	type request struct {
		Name       string               `json:"name" validate:"required"`
		Price      float32              `json:"price" validate:"required"`
		ProducerId uuid.UUID            `json:"producer_id" validate:"required"`
		Status     entity.ProductStatus `json:"status" validate:"oneof=published hidden"`
	}
	return func(c echo.Context) error {
		var requestData request
		err := c.Bind(&requestData)
		if err != nil {
			log.Println("Ошибка чтения тела запроса ", err)
			return c.NoContent(http.StatusBadRequest)
		}
		validate := validator.New()
		err = validate.Struct(requestData)
		if err != nil {
			log.Println("Невалидные данные ", err)
			return c.String(http.StatusBadRequest, "")
		}
		product := entity.Product{
			Name:  requestData.Name,
			Price: requestData.Price,
			Producer: &entity.Producer{
				Id: requestData.ProducerId,
			},
			Status: requestData.Status,
		}
		newProduct, err := p.usecase.Create(c.Request().Context(), product)
		if err != nil {
			switch {
			case err == usecase.ErrProducerNotFound:
				log.Println("Производителя с таким id не найдено", err)
				return c.String(http.StatusBadRequest, "")
			default:
				log.Println("Ошибка создания продукта", err)
				return c.NoContent(http.StatusInternalServerError)
			}
		}
		log.Println("Продукт создан")
		c.Response().Header().Set(echo.HeaderContentType, "application/json")
		return c.JSON(http.StatusOK, newProduct)
	}
}

func (p *ProductHandlers) getProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func (p *ProductHandlers) updateProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func (p *ProductHandlers) deleteProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Not Implemented")
	}
}
func RegisterProductRoutes(u usecase.Product, g *echo.Group) {
	a := NewProductHandlers(u)
	g.GET("", a.getAllProducts())
	g.POST("", a.createProduct())
	g.GET("/:id", a.getProductById())
	g.PUT("/:id", a.updateProductById())
	g.DELETE("/:id", a.deleteProductById())
}
