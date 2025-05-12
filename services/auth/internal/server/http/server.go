package http

import (
	"context"
	"log"
	"net"
	"net/http"

	_ "net/http/pprof"

	"github.com/labstack/echo"
)

type Server struct {
	echo *echo.Echo
	addr string
}

func NewHTTPServer(addr string, checker ReadyChecker) *Server {
	echoServer := echo.New()
	handler := NewHandler(checker)

	echoServer.GET("/health", handler.Health)
	echoServer.GET("/ready", handler.Ready)
	echoServer.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))

	return &Server{echo: echoServer, addr: addr}
}

func (s *Server) Serve(listener net.Listener) error {
	log.Printf("http server serving on %s", s.addr)

	return s.echo.Server.Serve(listener)
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("shutting down http server...")

	return s.echo.Shutdown(ctx)
}
