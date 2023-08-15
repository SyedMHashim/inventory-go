package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/apiserver/pkg/config"
	"example.com/apiserver/pkg/db"
	"example.com/apiserver/pkg/router"

	_ "github.com/go-sql-driver/mysql"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	cfg, err := config.LoadConfig()
	checkError(err)

	conn, err := db.Connect(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName)
	checkError(err)
	defer func() {
		log.Println("Closing DB connection")
		conn.Close()
	}()
	log.Println("Connected to DB")

	router := router.Initialise(conn)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: router,
	}

	go func() {
		// serve connections
		if srvErr := server.ListenAndServe(); srvErr != nil && !errors.Is(srvErr, http.ErrServerClosed) {
			log.Fatal("Server Error", srvErr)
		}
	}()

	// Listen for Shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutdown HTTP server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown Error")
	}

	close(quit)
	cancel()
}
