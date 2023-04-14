package httpgin

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"homework8/internal/app"
)

type Server struct {
	port string
	app  *gin.Engine
}

func LogMiddleWare(c *gin.Context) {
	t := time.Now()

	log.Printf("-- received request -- | method: %s | path: %s\n", c.Request.Method, c.Request.URL.Path)

	c.Next()

	latency := time.Since(t)
	status := c.Writer.Status()

	log.Printf("-- handled request -- | status: %d | latency: %+v | method: %s | path: %s\n", status, latency, c.Request.Method, c.Request.URL.Path)
}

func NewHTTPServer(port string, a app.App) Server {
	gin.SetMode(gin.ReleaseMode)
	s := Server{port: port, app: gin.New()}

	// todo: add your own logic
	api := s.app.Group("/api/v1")

	// MiddleWare для логирования и паник
	api.Use(gin.Logger())
	api.Use(gin.Recovery())

	api.Use(LogMiddleWare)

	AppRouter(api, a)
	return s
}

func (s *Server) Listen() error {
	return s.app.Run(s.port)
}

func (s *Server) Handler() http.Handler {
	return s.app
}
