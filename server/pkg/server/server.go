package server

import (
	"context"
	"my-collection/server/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var logger = utils.MustGetLogger("server")

type Server struct {
	addr     string
	router   *gin.Engine
	apiRoute *gin.RouterGroup
}

func New(addr string) *Server {
	gin.SetMode("release")

	router := gin.New()
	router.Use(contextInitializer)
	router.Use(cors.Default())
	router.Use(httpLogger)

	server := &Server{
		addr:     addr,
		router:   router,
		apiRoute: router.Group("/api"),
	}

	server.registerStaticHandlers()
	return server
}

func (s *Server) RegisterHandler(h Handler) {
	h.RegisterRoutes(s.apiRoute)
}

func (s *Server) registerStaticHandlers() {
	s.router.Static("/ui", "ui/")
	s.router.StaticFile("/", "ui/index.html")
	s.router.GET("/spa/*route", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "ui/index.html")
	})
}

func (s *Server) Run(ctx context.Context) error {
	logger.Infof(ctx, "Starting server at address %s", s.addr)

	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}

	serverErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Infof(ctx, "Context cancelled, shutting down server gracefully...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Errorf(ctx, "Server forced to shutdown: %v", err)
			return err
		}

		logger.Infof(ctx, "Server exited gracefully")
		return nil
	case err := <-serverErr:
		return err
	}
}
