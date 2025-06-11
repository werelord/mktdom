package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		err := http.ListenAndServe(":9090", http.FileServer(http.Dir("assets")))
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("server closed")
		} else if err != nil {
			fmt.Println("Failed to start server", err)
			// force close gracefully
			done <- syscall.SIGINT
		}
	}()
	<-done
	return
}
