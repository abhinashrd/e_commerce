package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

var db *pgxpool.Pool

func main() {
	// connect DB
	url := os.Getenv("PG_URL")
	if url == "" {
		url = "postgres://postgres:postgres@localhost:5432/productsvc?sslmode=disable"
	}
	ctx := context.Background()
	var err error
	db, err = pgxpool.New(ctx, url)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// routes
	r.POST("/products", createProduct)
	r.GET("/products", listProducts)
	r.GET("/products/:id", getProduct)

	srv := &http.Server{Addr: ":8082", Handler: r}
	go srv.ListenAndServe()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctxShut, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctxShut)
}

// Create product
func createProduct(c *gin.Context) {
	var p Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.QueryRow(context.Background(),
		"INSERT INTO products(name, price) VALUES($1,$2) RETURNING id",
		p.Name, p.Price).Scan(&p.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

// List all products
func listProducts(c *gin.Context) {
	rows, err := db.Query(context.Background(),
		"SELECT id, name, price FROM products ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err == nil {
			products = append(products, p)
		}
	}
	c.JSON(http.StatusOK, products)
}

// Get product by ID
func getProduct(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var p Product
	err := db.QueryRow(context.Background(),
		"SELECT id, name, price FROM products WHERE id=$1", id).
		Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}
