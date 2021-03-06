package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	API       *echo.Echo
	Db        *mongo.Database
	Client    *mongo.Client
	TLSConfig *tls.Config
	Debug     bool
	LogFile   *os.File
}

func (a *App) Initialize() {
	var err error
	a.API = echo.New()

	a.LogFile, err = os.OpenFile("rest-server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}

	// TODO: debug-timeout
	a.API.Validator = &CustomValidator{validator: validator.New()}
	a.API.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: a.LogFile}))
	a.API.Use(middleware.Recover())
	a.API.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 240 * time.Second,
	}))
	err = godotenv.Load()
	if err != nil {
		log.Print("Reading environment failed.")
	}
	a.Debug = GetDebug()
	origins, err := SplitOrigins()
	if err != nil {
		a.API.Logger.Fatalf("parsing origins failed: %s", err)
	}
	a.API.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
		AllowMethods: []string{http.MethodGet},
	}))
	if a.Debug {
		a.API.Logger.Printf(fmt.Sprintf("CORS: %v", origins))
	}

	a.Db, a.Client, err = a.getDbConnection()
	if err != nil {
		a.API.Logger.Fatal("initializing db connection failed: %s", err)
	}
	a.API.Logger.Printf("database connection succeed db: %s", a.Db.Name())
	a.API.GET("/", a.getHealthz)

	apiEndpoints := a.API.Group("/api/v1")

	route := apiEndpoints.GET("/tutkintoKoulutukset", a.getExaminations)
	route.Name = "get-examinations"
	a.API.Logger.Info("service initialized succesfully")
}

func (a *App) Run() {
	defer a.LogFile.Close()
	if os.Getenv("APP_SSL_PUBLIC") != "" && os.Getenv("APP_SSL_PRIVATE") != "" {
		a.API.Logger.Fatal(a.API.StartTLS(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), os.Getenv("APP_SSL_PUBLIC"), os.Getenv("APP_SSL_PRIVATE")))
	}
	a.API.Logger.Fatal("no ssl-files, nothing to do")
}
