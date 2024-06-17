package main

import (
	"log"
	"net/http"
	"os"
	"rate-limiter/limiter"
	"rate-limiter/limiter/middleware"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	rateLimiter := limiter.NewRedisLimiter(rdb)
	http.Handle("/", middleware.RateLimitMiddleware(rateLimiter)(http.HandlerFunc(handler)))

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
