package server

import (
	"context"
	//"errors"
	"log"
	"net/http"
	//"strings"
	"strconv"
	//"sync"
	"time"

	"github.com/dchote/survey-engine/config"
	"github.com/dchote/survey-engine/server/handlers"

	"github.com/jinzhu/gorm"

	"github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	e *echo.Echo
)

// StartServer starts the HTTP server
func StartServer(cfg config.SurveyConfig, db *gorm.DB, assets *rice.Box) {
	if e != nil {
		return
	}

	// instantiate echo instance
	e = echo.New()
	e.HideBanner = true
	e.Server.ReadTimeout = 60 * time.Second
	e.Server.WriteTimeout = 240 * time.Second

	// JWT config
	jwtConfig := middleware.JWTConfig{
		Claims:     &handlers.JWTCustomClaims{},
		SigningKey: []byte(cfg.Server.JWTSigningKey),
	}

	// setup middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// prevent caching by client (e.g. Safari)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			return next(c)
		}
	})

	if assets != nil {
		assetHandler := http.FileServer(assets.HTTPBox())
		e.GET("/", echo.WrapHandler(assetHandler))
		e.GET("/favicon.ico", echo.WrapHandler(assetHandler))
		e.GET("/css/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
		e.GET("/js/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
		e.GET("/fonts/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
		e.GET("/img/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
	}

	// Base API
	e.GET("/health", handlers.Health())
	e.GET("/config", handlers.Config())

	// User API
	handlers.InitUserHandler(db)
	e.POST("/v1/login", handlers.Login())
	e.POST("/v1/register", handlers.Register())

	userAPI := e.Group("/v1/user")
	userAPI.Use(middleware.JWTWithConfig(jwtConfig))
	userAPI.GET("/info", handlers.UserInfo())
	userAPI.GET("/list", handlers.ListUsers())
	userAPI.GET("/:name", handlers.FetchUser())
	userAPI.POST("/:name", handlers.UpdateUser())
	userAPI.DELETE("/:name", handlers.DeleteUser())

	// Survey API
	handlers.InitSurveyHandler(db)
	e.GET("/v1/survey/:name", handlers.FetchSurvey())
	e.GET("/v1/survey/settings", handlers.SurveySettings())
	e.GET("/v1/survey/list", handlers.ListSurveys())

	// survey session
	e.POST("/v1/survey/newSession", handlers.NewSurveySession())
	e.POST("/v1/survey/session/:id", handlers.UpdateSurveySession())
	e.GET("/v1/survey/session/:id", handlers.FetchSurveySession())
	e.GET("/v1/survey/sessionDump/:id", handlers.FetchAllSurveySessionData())

	// survey session data
	e.POST("/v1/survey/newSessionData", handlers.NewSurveySessionData())
	e.POST("/v1/survey/sessionData/:id", handlers.UpdateSurveySessionData())
	e.GET("/v1/survey/sessionData/:id", handlers.FetchSurveySessionData())
	e.DELETE("/v1/survey/sessionData/:id", handlers.DeleteSurveySessionData())

	// Survey API authenticated functions
	surveyAPI := e.Group("/v1/survey")
	surveyAPI.Use(middleware.JWTWithConfig(jwtConfig))
	surveyAPI.POST("/settings", handlers.UpdateSurveySettings())
	surveyAPI.POST("/:name", handlers.UpdateSurvey())
	surveyAPI.DELETE("/:name", handlers.DeleteSurvey())
	surveyAPI.GET("/sessionList", handlers.ListSurveySessions())
	surveyAPI.GET("/export", handlers.ExportData())
	surveyAPI.GET("/reset", handlers.ResetData())

	// Logic API
	handlers.InitLogicHandler(db)

	e.GET("/functions.js", handlers.BuildFunctionsJS())
	e.GET("/v1/logic/all", handlers.ListLogicFunctions())

	// Logic API Authenticated Functions
	logicAPI := e.Group("/v1/logic")
	logicAPI.Use(middleware.JWTWithConfig(jwtConfig))

	logicAPI.GET("/:functionName", handlers.FetchLogicFunction())
	logicAPI.DELETE("/:functionName", handlers.DeleteLogicFunction())
	logicAPI.POST("/:functionName", handlers.UpdateLogicFunction())

	// Instructions API and pass in our static assets
	handlers.InitInstructionHandler(db, assets)

	e.GET("/instructionHTML/:name", handlers.BuildInstructionHTML())
	e.GET("/instructionHTML/:name/:page", handlers.BuildInstructionHTML())

	e.GET("/v1/instructions/all", handlers.ListInstructions())
	e.GET("/v1/instructions/:name", handlers.FetchInstruction())

	// Instructions API authenticated functions
	instructionAPI := e.Group("/v1/instructions")
	instructionAPI.Use(middleware.JWTWithConfig(jwtConfig))

	instructionAPI.POST("/:name", handlers.UpdateInstruction())
	instructionAPI.DELETE("/:name", handlers.DeleteInstruction())

	// Images API
	handlers.InitImageHandler(db)
	e.GET("/v1/images/all", handlers.ListImages())
	e.GET("/v1/images/:name", handlers.FetchImage())
	e.GET("/images/:name", handlers.ViewImage())

	// Images API authenticated functions
	imagesAPI := e.Group("/v1/images")
	imagesAPI.Use(middleware.JWTWithConfig(jwtConfig))

	imagesAPI.POST("/:name", handlers.UploadImage())
	imagesAPI.DELETE("/:name", handlers.DeleteImage())

	// Email Notifications
	handlers.InitEmailHandler()
	e.POST("/v1/notify", handlers.SendNotification())

	// Start the server
	log.Println("starting server on http://" + cfg.Server.ListenAddress + ":" + strconv.Itoa(cfg.Server.ListenPort))
	e.Logger.Fatal(e.Start(cfg.Server.ListenAddress + ":" + strconv.Itoa(cfg.Server.ListenPort)))
}

// StopServer stop the HTTP server
func StopServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
