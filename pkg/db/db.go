package db

import (
	"database/sql"
	"fmt"
	"time"

	"log"

	"github.com/cenk/backoff"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Product struct {
	Id    int
	Name  string
	Price float64
}

func Connect(DbHost, DbPort, DbUser, DbPass, DbName string) (*sql.DB, error) {
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
					products(id INT NOT NULL AUTO_INCREMENT,name VARCHAR(30) NOT NULL,price FLOAT NOT NULL,PRIMARY KEY(id))`)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func GetProducts(conn *sql.DB) ([]Product, error) {
	rows, err := conn.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	var Products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.Id, &p.Name, &p.Price)
		if err != nil {
			return nil, err
		}
		Products = append(Products, p)
	}
	return Products, nil
}

func GetProduct(conn *sql.DB, id int) (Product, error) {
	var p Product
	rows := conn.QueryRow(fmt.Sprintf("SELECT * FROM products WHERE id = %d", id))
	err := rows.Scan(&p.Id, &p.Name, &p.Price)
	if err != nil {
		return p, err
	}
	return p, nil
}

func CreateProduct(conn *sql.DB, p *Product) error {
	query := fmt.Sprintf("INSERT INTO products (name, price) VALUES ('%s','%f')", p.Name, p.Price)
	result, err := conn.Exec(query)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.Id = int(id)
	return nil
}

func UpdateProduct(conn *sql.DB, p Product) error {
	query := fmt.Sprintf("UPDATE products SET name = '%s', price = '%f' WHERE id = %d", p.Name, p.Price, p.Id)
	_, err := conn.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func DeleteProduct(conn *sql.DB, id int) error {
	query := fmt.Sprintf("DELETE FROM products WHERE id = %d", id)
	_, err := conn.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
