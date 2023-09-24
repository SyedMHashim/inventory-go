package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	web "example.com/apiserver/internal/adapter/in/http"
	"example.com/apiserver/internal/adapter/out/mysql"
	"example.com/apiserver/internal/application/domain"
	"example.com/apiserver/internal/application/service"
	"example.com/apiserver/internal/config"

	"github.com/cenk/backoff"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := ConnectToMySQL(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Println("Closing DB connection")
		db.Close()
	}()
	log.Println("Connected to DB")

	router := mux.NewRouter()

	productPersistenceAdapter := mysql.NewProductPersistenceAdapter(db)

	createProductService := service.NewCreateProductService(domain.GenerateOrderID, productPersistenceAdapter)
	deleteProductService := service.NewDeleteProductService(productPersistenceAdapter)
	getProductService := service.NewGetProductService(productPersistenceAdapter)
	updateProductService := service.NewUpdateProductService(productPersistenceAdapter)

	createProductHandler := web.NewCreateProductHandler(createProductService)
	deleteProductHandler := web.NewDeleteProductHandler(deleteProductService)
	getProductHandler := web.NewGetProductHandler(getProductService)
	updateProductHandler := web.NewUpdateProductHandler(updateProductService)

	router.HandleFunc("/product", createProductHandler.HandleCreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/product/{id}", deleteProductHandler.HandleDeleteProduct).Methods(http.MethodDelete)
	router.HandleFunc("/product/{id}", getProductHandler.HandleGetProduct).Methods(http.MethodGet)
	router.HandleFunc("/products", getProductHandler.HandleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/product/{id}", updateProductHandler.HandleUpdateProduct).Methods(http.MethodPut)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: router,
	}

	go func() {
		// serve connections
		log.Printf("Listening at :%s\n", cfg.ServerPort)
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

func ConnectToMySQL(DbHost, DbPort, DbUser, DbPass, DbName string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DbUser, DbPass, DbHost, DbPort, DbName)
	conn, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	notify := func(err error, duration time.Duration) {
		log.Println("Connecting to db...")
	}

	exponentialBackoff := backoff.NewExponentialBackOff()
	exponentialBackoff.InitialInterval = time.Second * 5
	retrier := backoff.WithMaxRetries(exponentialBackoff, 3)

	err = backoff.RetryNotify(func() error {
		return conn.Ping()
	}, retrier, notify)
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS 
					products(id VARCHAR(36) NOT NULL,name VARCHAR(30) NOT NULL,price FLOAT NOT NULL,PRIMARY KEY(id))`)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
