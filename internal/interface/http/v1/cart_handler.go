package http

import (
	"github.com/labstack/echo/v4"
	"github.com/pkittipat/try-cart/internal/app/service"
)

type cartHandler struct {
	cartSrv *service.CartService
}

func RegisterCartHandler(
	router *echo.Group,
	cartSrv *service.CartService,
) {
	handler := cartHandler{
		cartSrv: cartSrv,
	}

	router.POST("/", handler.CreateCart)
}

// Added Product to cart
// If there is no cart then create and push product into cart.
func (h *cartHandler) CreateCart(e echo.Context) error {
	return nil
}
