package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_customMiddleware "github.com/andhikagama/lmnlo/cmiddleware/usecase"
	cfg "github.com/andhikagama/lmnlo/config"
	userHandler "github.com/andhikagama/lmnlo/user/delivery"
	_userRepository "github.com/andhikagama/lmnlo/user/repository"
	_userUsecase "github.com/andhikagama/lmnlo/user/usecase"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
)

var config cfg.Config

func init() {
	config = cfg.NewViperConfig()
	log.SetFormatter(&log.JSONFormatter{})
	if config.GetBool(`debug`) {
		log.Warn(`Lmnlo is running in debug mode`)
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	//Setup Database Connection
	dbHost := config.GetString(`database.host`)
	dbPort := config.GetString(`database.port`)
	dbUser := config.GetString(`database.user`)
	dbPass := config.GetString(`database.pass`)
	dbName := config.GetString(`database.name`)

	dsn := dbUser + `:` + dbPass + `@tcp(` + dbHost + `:` + dbPort + `)/` + dbName + `?parseTime=1`
	log.Info("connecting to database")
	db, err := sql.Open(`mysql`, dsn)
	if err != nil {
		log.Error(fmt.Sprintf("database connection failed. Err: %v", err.Error()))
		os.Exit(1)
	}

	defer db.Close()

	e := echo.New()

	// For Health Check
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong!")
	})

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.PATCH},
		ExposeHeaders: []string{`X-Cursor`},
	}))

	gv1 := e.Group(`/v1`)
	gv1.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.PATCH},
		ExposeHeaders: []string{`X-Cursor`},
	}))

	//Initiate Repository for each entity
	userRepository := _userRepository.NewUserRepository(db)

	// Initiate Custom Middleware
	customMiddleware := _customMiddleware.NewMiddlewareUsecase(userRepository)
	gv1.Use(customMiddleware.CheckAuthHeader)

	//Initiate Usecase for each entity
	userUsecase := _userUsecase.NewUserUsecase(userRepository)

	//Initiate Handler for each entity
	userHandler.NewUserHTTPHandler(gv1, userUsecase)

	log.Infof(`Connected to database : %v on %v`, config.GetString(`database.name`), config.GetString(`database.host`))
	log.Infof(`Lmnlo server running at address : %v`, config.GetString(`server.address`))
	e.Start(config.GetString("server.address"))

}
