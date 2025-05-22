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
		Name       string               `json:"name" validate:"required,max=100,min=3"`
		Price      float32              `json:"price" validate:"required"`
		ProducerId uuid.UUID            `json:"producer_id" validate:"required,uuid"`
		Status     entity.ProductStatus `json:"status" validate:"required,oneof=published hidden"`
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
		productId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Println("Невалидный id", err)
			return c.NoContent(http.StatusBadRequest)
		}
		product, err := p.usecase.GetById(c.Request().Context(), productId)
		if err != nil {
			switch {
			case err == usecase.ErrProductNotFound:
				log.Println("Товар с таким id не найден", err)
				return c.NoContent(http.StatusBadRequest)
			default:
				log.Println("Ошибка получения товара", err)
				return c.NoContent(http.StatusInternalServerError)
			}
		}
		return c.JSON(http.StatusOK, product)
	}
}
func (p *ProductHandlers) updateProductById() echo.HandlerFunc {
	type request struct {
		Name       string               `json:"name" validate:"omitempty,max=100,min=3"`
		Price      float32              `json:"price" validate:"omitempty"`
		ProducerId uuid.UUID            `json:"producer_id" validate:"omitempty,uuid"`
		Status     entity.ProductStatus `json:"status" validate:"omitempty,oneof=published hidden"`
	}
	return func(c echo.Context) error {
		productId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Println("Невалидный id", err)
			return c.NoContent(http.StatusBadRequest)
		}
		var requestData request
		err = c.Bind(&requestData)
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
			Id:    productId,
			Name:  requestData.Name,
			Price: requestData.Price,
			Producer: &entity.Producer{
				Id: requestData.ProducerId,
			},
			Status: requestData.Status,
		}
		updatedProduct, err := p.usecase.Update(c.Request().Context(), product)
		if err != nil {
			switch {
			case err == usecase.ErrProductNotFound:
				log.Println("Товар с таким id не найдено", err)
				return c.NoContent(http.StatusBadRequest)
			case err == usecase.ErrProducerNotFound:
				log.Println("Производителя с таким id не найдено", err)
				return c.NoContent(http.StatusBadRequest)
			default:
				log.Println("Ошибка обновления товара", err)
				return c.NoContent(http.StatusInternalServerError)
			}
		}
		log.Println("Товар обновлен")
		c.Response().Header().Set(echo.HeaderContentType, "application/json")
		return c.JSON(http.StatusOK, updatedProduct)
	}
}
func (p *ProductHandlers) deleteProductById() echo.HandlerFunc {
	return func(c echo.Context) error {
		productId, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Println("Невалидный id", err)
			return c.NoContent(http.StatusBadRequest)
		}
		err = p.usecase.DeleteById(c.Request().Context(), productId)
		if err != nil {
			switch {
			case err == usecase.ErrProductNotFound:
				log.Println("Товар с таким id не найден", err)
				return c.NoContent(http.StatusBadRequest)
			default:
				log.Println("Ошибка удаления товара", err)
				return c.NoContent(http.StatusInternalServerError)
			}
		}
		log.Println("Товар удален")
		return c.NoContent(http.StatusOK)
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
