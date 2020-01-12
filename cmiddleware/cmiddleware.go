package cmiddleware

import "github.com/labstack/echo"

// Usecase ...
type Usecase interface {
	CheckAuthHeader(next echo.HandlerFunc) echo.HandlerFunc
}
