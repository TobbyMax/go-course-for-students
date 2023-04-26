package httpgin

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"homework9/internal/app"
)

func LoggerMiddleWare(c *gin.Context) {
	start := time.Now()

	log.Printf("-- received request -- | protocol: HTTP | method: %s | path: %s\n", c.Request.Method, c.Request.URL.Path)

	c.Next()

	latency := time.Since(start)
	status := c.Writer.Status()

	log.Printf("-- handled request -- | protocol: HTTP | status: %d | latency: %+v | method: %s | path: %s\n", status, latency, c.Request.Method, c.Request.URL.Path)
}

func NewHTTPServer(port string, a app.App) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	s := &http.Server{Addr: port, Handler: handler}

	// todo: add your own logic

	api := handler.Group("/api/v1")

	// MiddleWare для логирования и паник
	api.Use(gin.Logger())
	api.Use(gin.Recovery())

	api.Use(LoggerMiddleWare)

	AppRouter(api, a)
	return s
}
