package http

import (
	"net/http"

	"github.com/andhikagama/lmnlo/models/entity"
	"github.com/andhikagama/lmnlo/models/response"
	"github.com/andhikagama/lmnlo/user"
	"github.com/labstack/echo"
)

// UserHTTPHandler ...
type UserHTTPHandler struct {
	Usecase user.Usecase
}

// NewUserHTTPHandler ...
func NewUserHTTPHandler(g *echo.Group, u user.Usecase) {
	handler := &UserHTTPHandler{
		Usecase: u,
	}

	g.POST(`/register`, handler.Register)
}

// Register ...
func (h *UserHTTPHandler) Register(c echo.Context) error {
	usr := new(entity.User)
	c.Bind(usr)

	err := h.Usecase.Register(usr)

	if err != nil {
		if err == response.ErrAlreadyExist {
			return c.JSON(http.StatusConflict, &response.Wrapper{
				Message: response.ErrAlreadyExist.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, &response.Wrapper{
			Message: response.ErrServer.Error(),
		})
	}

	return c.JSON(http.StatusOK, usr)
}
