package http

import (
	"io/ioutil"
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
	g.GET(`/user/:id`, handler.GetByID)
	g.DELETE(`/user/:id`, handler.Delete)
	g.PATCH(`/user/:id`, handler.PartialUpdate)
	g.POST(`/login`, handler.Login)
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

// GetByID ...
func (h *UserHTTPHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param(`id`))
	if err != nil || id == 0 {
		return c.JSON(http.StatusNotFound, &response.Wrapper{
			Message: response.ErrNotFound.Error(),
		})
	}

	res, err := h.Usecase.GetByID(int64(id))

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

	return c.JSON(http.StatusOK, res)
}

// Delete ...
func (h *UserHTTPHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param(`id`))
	if err != nil || id == 0 {
		return c.JSON(http.StatusNotFound, &response.Wrapper{
			Message: response.ErrNotFound.Error(),
		})
	}

	err = h.Usecase.Delete(int64(id))

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

	return c.NoContent(http.StatusNoContent)
}

// PartialUpdate ...
func (h *UserHTTPHandler) PartialUpdate(c echo.Context) error {
	var id int64

	if c.Param(`id`) != `` {
		intID, err := strconv.Atoi(c.Param(`id`))
		if err != nil {
			return c.JSON(http.StatusNotFound, &response.Wrapper{
				Message: response.ErrNotFound.Error(),
			})
		}

		id = int64(intID)
	}

	jsonPatch, _ := ioutil.ReadAll(c.Request().Body)
	res, err := h.Usecase.PartialUpdate(id, jsonPatch)

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

	return c.JSON(http.StatusOK, res)
}

// Login ...
func (h *UserHTTPHandler) Login(c echo.Context) error {
	auth := new(entity.User)
	c.Bind(auth)

	res, err := h.Usecase.Login(auth)
	if err != nil {
		if err == response.ErrLogin {
			return c.JSON(http.StatusNotFound, &response.Wrapper{
				Message: err.Error(),
			})
		}

		if err == response.ErrForbidden {
			return c.JSON(http.StatusForbidden, &response.Wrapper{
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, &response.Wrapper{
			Message: response.ErrServer.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}
