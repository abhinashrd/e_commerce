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

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

var db *pgxpool.Pool

func main() {
	// connect DB
	url := os.Getenv("PG_URL")
	if url == "" {
		url = "postgres://postgres:postgres@localhost:5432/usersvc?sslmode=disable"
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

	r.POST("/users", createUser)
	r.GET("/users/:id", getUser)
	r.GET("/users", listUsers)

	srv := &http.Server{Addr: ":8080", Handler: r}
	go srv.ListenAndServe()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctxShut, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctxShut)
}

// Create user
func createUser(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.QueryRow(context.Background(),
		"INSERT INTO users(name, email) VALUES($1,$2) RETURNING id, created_at",
		u.Name, u.Email).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, u)
}

// Get user by ID
func getUser(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var u User
	err := db.QueryRow(context.Background(),
		"SELECT id, name, email, created_at FROM users WHERE id=$1", id).
		Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

// List all users
func listUsers(c *gin.Context) {
	rows, err := db.Query(context.Background(),
		"SELECT id, name, email, created_at FROM users ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, u)
	}
	c.JSON(http.StatusOK, users)
}
