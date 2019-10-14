package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Endpoint interface {
	PerformRequest(r *gin.Engine)
}

type exchangeServer struct {
	endpoints map[string]Endpoint
	srv       *http.Server
}

// CreateNewServer - Creates a new server that will
//					 respond to requests.
func CreateNewServer() *exchangeServer {
	return &exchangeServer{make(map[string]Endpoint), nil}
}

// Register - Register the end points so the server
//			  can perform the requests against the
//			  registered requests.
func (s *exchangeServer) Register(name string, e Endpoint) {
	s.endpoints[name] = e
}

func (s *exchangeServer) Start() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	for name, e := range s.endpoints {
		fmt.Println("Key:", name, "Value:", e)
		e.PerformRequest(router)
	}

	s.srv = &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.Stop(5 * time.Second)
}

func (s *exchangeServer) Stop(timeout time.Duration) {
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}
