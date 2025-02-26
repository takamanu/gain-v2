package main

import (
	"context"
	"fmt"
	"gain-v2/configs"
	"gain-v2/helper"
	email "gain-v2/helper/email"
	encrypt "gain-v2/helper/encrypt"
	"gain-v2/routes"
	"gain-v2/utils/database"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	dataUser "gain-v2/features/users/data"
	handlerUser "gain-v2/features/users/handler"
	serviceUser "gain-v2/features/users/service"

	dataLogging "gain-v2/features/logging/data"
	handlerLogging "gain-v2/features/logging/handler"
	serviceLogging "gain-v2/features/logging/service"
)

func main() {
	// Initialize Echo framework
	e := echo.New()

	// Load configuration settings
	config := configs.InitConfig()

	// Create a background context
	ctx := context.Background()

	// Initialize PostgreSQL database connection
	db, err := database.InitDBPostgres(*config)
	if err != nil {
		e.Logger.Fatal("Database initialization failed: ", err.Error())
	}

	// Initialize Redis database connection
	redis, _ := database.InitRedis(*config)
	if err != nil {
		e.Logger.Fatal("Redis initialization failed: ", err.Error())
	}

	// Initialize Elasticsearch database connection
	elastic, _ := database.InitElasticSearch(*config)
	if err != nil {
		e.Logger.Fatal("Elasticsearch initialization failed: ", err.Error())
	}

	e.IPExtractor = echo.ExtractIPFromRealIPHeader()

	// Create helper instances
	encryptHelper := encrypt.New()
	emailHelper := email.New(*config)
	jwtInterface := helper.New(config.Secret, config.RefSecret)

	// Initialize feature instances
	userModel := dataUser.NewData(db, redis, ctx)
	loggingModel := dataLogging.NewData(redis, elastic)

	userServices := serviceUser.NewService(userModel, jwtInterface, emailHelper, encryptHelper)
	loggingServices := serviceLogging.NewService(loggingModel)

	userController := handlerUser.NewHandler(userServices, jwtInterface)
	loggingController := handlerLogging.NewHandler(loggingServices, jwtInterface)

	// Set up API routes
	group := e.Group("/api/v1")

	routes.RouteUser(group, userController, *config)
	routes.RouteLogging(group, loggingController, *config)

	// Handle "not found" errors for specific endpoints
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Endpoint not found", nil))
	})

	e.GET("/api", func(c echo.Context) error {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Endpoint not found", nil))
	})

	e.GET("/api/v1", func(c echo.Context) error {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Endpoint not found", nil))
	})

	// Configure middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}, time=${time_rfc3339}\n",
		}))

	// Debug logging for database connections
	e.Logger.Debug(db)
	e.Logger.Debug(redis)
	e.Logger.Debug(elastic)

	// Start the server
	e.Logger.Info(fmt.Sprintf("Server is listening on port :%d", config.ServerPort))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.ServerPort)).Error())
}
