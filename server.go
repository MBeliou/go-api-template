package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/MBeliou/go-api-template/handler"
	"github.com/asdine/storm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {

	port := flag.Int("port", 1323, "The port the server will listen to")
	dbName := flag.String("database", "mydb.db", "which database file to use. Will be created if does not exist")
	flag.Parse()

	e := echo.New()

	//	Setup middleware
	e.Logger.SetLevel(log.ERROR)
	e.Use(middleware.Logger())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(handler.JwtKey),
		Skipper: func(c echo.Context) bool {
			// Skip authentication for and signup login requests
			if c.Path() == "/login" || c.Path() == "/signup" || c.Path() == "/users" || c.Path() == "/fetchAll" {
				return true
			}
			return false
		},
	}))

	e.Use(middleware.CORS())

	db, err := storm.Open(*dbName)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Database connection
	h := &handler.Handler{DB: db}

	// Routes
	e.POST("/signup", h.Signup)
	e.POST("/login", h.Login)
	e.POST("/follow/:id", h.Follow)

	e.GET("/feed", h.FetchPosts)
	e.POST("/post", h.CreatePost)

	e.GET("/users", h.GetAll) // TODO: only in dev
	e.GET("/fetchAll", h.FetchAllPosts)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", *port)); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
