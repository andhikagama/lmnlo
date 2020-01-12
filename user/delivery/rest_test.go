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
		mockUCase.On("Store", mock.AnythingOfType(`*entity.User`)).Return(nil).Once()

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
		mockUCase.On("Store", mock.AnythingOfType(`*entity.User`)).Return(errors.New(`error`)).Once()

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
