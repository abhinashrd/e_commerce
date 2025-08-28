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

type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userId"`
	ProductID int64     `json:"productId"`
	Quantity  int       `json:"quantity"`
	Total     int64     `json:"total"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

var db *pgxpool.Pool

func main() {
	// connect DB
	url := os.Getenv("PG_URL")
	if url == "" {
		url = "postgres://postgres:postgres@localhost:5432/ordersvc?sslmode=disable"
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
	r.POST("/orders", createOrder)
	r.GET("/orders/:id", getOrder)

	srv := &http.Server{Addr: ":8083", Handler: r}
	go srv.ListenAndServe()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctxShut, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctxShut)
}

// Create order handler
func createOrder(c *gin.Context) {
	var o Order
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Placeholder: you could call productsvc here to fetch price
	// For now assume total = quantity * 75000 (hardcoded for Laptop)
	o.Total = int64(o.Quantity) * 75000

	err := db.QueryRow(context.Background(),
		"INSERT INTO orders(user_id, product_id, quantity, total) VALUES($1,$2,$3,$4) RETURNING id, created_at",
		o.UserID, o.ProductID, o.Quantity, o.Total).Scan(&o.ID, &o.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, o)
}

// Get order by ID
func getOrder(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var o Order
	err := db.QueryRow(context.Background(),
		"SELECT id, user_id, product_id, quantity, total, created_at FROM orders WHERE id=$1", id).
		Scan(&o.ID, &o.UserID, &o.ProductID, &o.Quantity, &o.Total, &o.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, o)
}
