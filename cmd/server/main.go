package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	ServerAddress = "werelord.lan" // take advantage of routers dns
	ServerPort    = ":9090"
)

func main() {

	srv := &http.Server{
		Addr:    ServerAddress + ServerPort,
		Handler: http.FileServer(http.Dir("assets")),
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		err := srv.ListenAndServe()
		// err := http.ListenAndServe(ServerPort, )
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("server closed")
		} else if err != nil {
			fmt.Println("Failed to start server", err)
			// force close gracefully
			done <- syscall.SIGINT
		}
	}()
	fmt.Printf("server started on %v", srv.Addr)
	<-done
	fmt.Println("server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() { cancel() }()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("server shutdown failed: %+v", err)
	}

	fmt.Println("server exited properly")

}
