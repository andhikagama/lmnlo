package http

import (
	"net/http"
	"strconv"

	"github.com/andhikagama/lmnlo/models/entity"
	"github.com/andhikagama/lmnlo/models/filter"
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
	g.GET(`/user`, handler.Fetch)
	g.PUT(`/user/:id`, handler.Update)
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

// Fetch ...
func (h *UserHTTPHandler) Fetch(c echo.Context) error {
	f := new(filter.User)

	f.Num = int64(200)
	if c.QueryParam(`num`) != `` {
		intNum, err := strconv.Atoi(c.QueryParam(`num`))

		if err != nil {
			return c.JSON(http.StatusBadRequest, &response.Wrapper{
				Message: response.ErrBadRequest.Error(),
			})
		}

		f.Num = int64(intNum)
	}

	f.Cursor = int64(0)
	strNextCursor := `0`
	if c.QueryParam(`cursor`) != `` {
		intCursor, err := strconv.Atoi(c.QueryParam(`cursor`))

		if err != nil {
			return c.JSON(http.StatusBadRequest, &response.Wrapper{
				Message: response.ErrBadRequest.Error(),
			})
		}

		f.Cursor = int64(intCursor)
		strNextCursor = c.QueryParam(`cursor`)
	}

	if c.QueryParam(`email`) != `` {
		f.Email = c.QueryParam(`email`)
	}

	if c.QueryParam(`address`) != `` {
		f.Address = c.QueryParam(`address`)
	}

	res, err := h.Usecase.Fetch(f)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &response.Wrapper{
			Message: response.ErrServer.Error(),
		})
	}

	if len(res) > 0 {
		nextCursor := res[len(res)-1].ID
		strNextCursor = strconv.Itoa(int(nextCursor))
	}

	c.Response().Header().Set(`X-Cursor`, strNextCursor)

	return c.JSON(http.StatusOK, res)
}

// Update ...
func (h *UserHTTPHandler) Update(c echo.Context) error {
	usr := new(entity.User)
	c.Bind(usr)

	id, err := strconv.Atoi(c.Param(`id`))
	if err != nil || id == 0 {
		return c.JSON(http.StatusNotFound, &response.Wrapper{
			Message: response.ErrNotFound.Error(),
		})
	}

	usr.ID = int64(id)

	err = h.Usecase.Update(usr)

	if err != nil {
		if err == response.ErrNotFound {
			return c.JSON(http.StatusNotFound, &response.Wrapper{
				Message: response.ErrNotFound.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, &response.Wrapper{
			Message: response.ErrServer.Error(),
		})
	}

	return c.JSON(http.StatusOK, usr)
}
