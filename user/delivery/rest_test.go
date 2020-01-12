package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andhikagama/lmnlo/models/entity"
	handler "github.com/andhikagama/lmnlo/user/delivery"
	"github.com/andhikagama/lmnlo/user/mocks"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockUser = entity.User{
	ID:       1,
	Email:    `andhika.gama@outlook.com`,
	Password: `aiueo`,
	Address:  `Menteng`,
}

var mockUsers = []*entity.User{
	&mockUser,
}

func TestStore(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockUCase := new(mocks.Usecase)
		mockUCase.On("Register", mock.AnythingOfType(`*entity.User`)).Return(nil).Once()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("user")
		handler := handler.UserHTTPHandler{
			Usecase: mockUCase,
		}
		handler.Register(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockUCase.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockUCase := new(mocks.Usecase)
		mockUCase.On("Register", mock.AnythingOfType(`*entity.User`)).Return(errors.New(`error`)).Once()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("user")

		handler := handler.UserHTTPHandler{
			Usecase: mockUCase,
		}

		handler.Register(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockUCase.AssertExpectations(t)
	})
}

func TestFetch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockUCase := new(mocks.Usecase)
		mockUCase.On("Fetch", mock.AnythingOfType(`*filter.User`)).Return(mockUsers, nil).Once()

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("user")

		handler := handler.UserHTTPHandler{
			Usecase: mockUCase,
		}
		handler.Fetch(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockUCase.AssertExpectations(t)
	})

	t.Run("success-with-param", func(t *testing.T) {
		mockUCase := new(mocks.Usecase)
		mockUCase.On("Fetch", mock.AnythingOfType(`*filter.User`)).Return(mockUsers, nil).Once()

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("user")
		c.QueryParams().Add(`email`, `andhika.gama@outlook.com`)
		c.QueryParams().Add(`num`, `100`)
		c.QueryParams().Add(`cursor`, `0`)
		c.QueryParams().Add(`address`, `men`)

		handler := handler.UserHTTPHandler{
			Usecase: mockUCase,
		}
		handler.Fetch(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockUCase.AssertExpectations(t)
	})

	t.Run("error-bad-param", func(t *testing.T) {
		mockUCase := new(mocks.Usecase)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("user")
		c.QueryParams().Add(`num`, `aaa`)

		handler := handler.UserHTTPHandler{
			Usecase: mockUCase,
		}
		handler.Fetch(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockUCase.AssertExpectations(t)
	})
	t.Run("error-bad-param", func(t *testing.T) {
		mockUCase := new(mocks.Usecase)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath("user")
		c.QueryParams().Add(`cursor`, `aaa`)

		handler := handler.UserHTTPHandler{
			Usecase: mockUCase,
		}
		handler.Fetch(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockUCase.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {

		mockUCase := new(mocks.Usecase)
		mockUCase.On("Fetch", mock.AnythingOfType(`*filter.User`)).Return(nil, errors.New(`Error`)).Once()

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/user", strings.NewReader(""))

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("user")
		handler := handler.UserHTTPHandler{
			Usecase: mockUCase,
		}
		handler.Fetch(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockUCase.AssertExpectations(t)
	})
}
