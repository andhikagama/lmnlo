package usecase

import (
	"net/http"
	"strings"

	cmware "github.com/andhikagama/lmnlo/cmiddleware"
	"github.com/andhikagama/lmnlo/helper"
	"github.com/andhikagama/lmnlo/models/response"
	"github.com/andhikagama/lmnlo/user"
	"github.com/labstack/echo"
)

const (
	token = `Authorization`
)

type cmwareUsecase struct {
	userRepo user.Repository
}

// NewMiddlewareUsecase ...
func NewMiddlewareUsecase(
	ar user.Repository,
) cmware.Usecase {
	return &cmwareUsecase{
		ar,
	}
}

// CheckAuthHeader ...
func (cm *cmwareUsecase) CheckAuthHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := ``
		if c.Request().Header.Get(`Authorization`) != `` {
			temp := strings.Split(c.Request().Header.Get(`Authorization`), ` `)
			token = temp[1]
		}

		if skipper(c) {
			return next(c)
		}

		ok, err := cm.userRepo.ValidateToken(token)
		if err != nil || !ok {
			return c.JSON(http.StatusUnauthorized, &response.Wrapper{
				Message: response.ErrUnAuthorized.Error(),
			})
		}

		if ok {
			usr, err := helper.ClaimTokenString(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, &response.Wrapper{
					Message: response.ErrUnAuthorized.Error(),
				})
			}
			c.Set(`user`, usr)
			return next(c)
		}

		return c.JSON(http.StatusUnauthorized, &response.Wrapper{
			Message: response.ErrUnAuthorized.Error(),
		})
	}

}

func skipper(c echo.Context) bool {
	path := c.Request().URL.Path
	ver := `v1/`
	realPath := strings.TrimLeft(path, ver)

	switch realPath {
	case `ping`, `login`, `register`:
		return true
	}
	return false
}
