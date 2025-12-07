package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var logger *slog.Logger

func guardAcceptHeaderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accept := c.Request().Header.Get("Accept")
		if accept == "" || accept == "*/*" || strings.Contains(accept, "application/json") {
			return next(c)
		}

		return c.JSON(http.StatusNotAcceptable, map[string]string{
			"message": "not acceptable",
			"detail":  fmt.Sprintf("accept header may be empty, */*, or application/json, got: '%s'", accept),
		})
	}
}

func guardContentTypeHeaderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		contentType := c.Request().Header.Get("Content-Type")
		if contentType == "" || contentType == "*/*" || strings.Contains(contentType, "application/json") {
			return next(c)
		}

		return c.JSON(http.StatusUnsupportedMediaType, map[string]string{
			"message": "invalid media type",
			"detail":  fmt.Sprintf("Content-Type header may be empty, */*, or application/json, got: '%s'", contentType),
		})
	}
}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	key, found := os.LookupEnv("API_KEY")

	if !found {
		return func(c echo.Context) error {
			return next(c)
		}
	}

	return func(c echo.Context) error {
		supplied := c.Request().Header.Get("X-Auth-Key")

		if supplied == key {
			return next(c)
		}

		return c.JSON(http.StatusUnauthorized, map[string]any{
			"message": "unauthorized",
			"detail":  "X-Auth-Key header is required",
		})
	}
}

func handler(c echo.Context) error {
	req := c.Request()

	var body any

	// Only decode if there is a body
	if req.ContentLength > 0 {
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			// If body is not valid JSON, fallback to string
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(req.Body)
			body = buf.String()
		}
	}

	// Remove the auth key from the logs...
	sanitisedHeaders := req.Header
	sanitisedHeaders.Del("X-Auth-Key")

	// Build response
	resp := map[string]any{
		"method":     req.Method,
		"path":       req.URL.Path,
		"query":      req.URL.Query(),
		"headers":    sanitisedHeaders,
		"body":       body,
		"request-id": c.Response().Header().Get(echo.HeaderXRequestID),
	}

	logger.Info("request received", "request", resp, "response-code", http.StatusOK)
	return c.JSON(http.StatusOK, resp)
}

func initLogger() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func getApiPort() (uint16, error) {
	var apiPort uint16 = 8080
	apiPortEnv, found := os.LookupEnv("API_PORT")

	if !found {
		return apiPort, nil
	}

	converted, err := strconv.ParseUint(apiPortEnv, 10, 16)

	if err != nil {
		return 0, err
	}

	return uint16(converted), nil
}

func main() {
	initLogger()

	apiPort, err := getApiPort()

	if err != nil {
		logger.Error("invalid port specified", "error", err)
		os.Exit(1)
	}

	e := echo.New()
	e.HideBanner = true

	// Middleware: Recover panics
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())

	// Middleware: Enforce Accept header
	e.Use(authMiddleware)
	e.Use(guardAcceptHeaderMiddleware)
	e.Use(guardContentTypeHeaderMiddleware)

	// Handler: just returns JSON
	e.Any("/*", handler)

	// Start server in goroutine for graceful shutdown
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", apiPort)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
