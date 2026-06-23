package main

import (
	"log"
	"net/http/httputil"
	"net/url"
	"rate_limiter/middleware"
	"rate_limiter/services"

	"github.com/gin-gonic/gin"
)

func main() {
	backendUrl, err := url.Parse("http://localhost:9001")
	if err != nil {
		log.Default().Fatalf("failed to parse backend URL: %v", err)
		return
	}

	serviceList := services.InitServices()
	if serviceList == nil {
		log.Fatal("Failed to initialize services")
		return
	}

	server := gin.Default()

	proxy := httputil.NewSingleHostReverseProxy(backendUrl)

	// Set this when cloudflare is deployed
	server.SetTrustedProxies(nil)

	server.Use(gin.Recovery(), gin.Logger(), middleware.RateLimitMiddleware(serviceList))

	server.NoRoute(func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	if err := server.Run(":8080"); err != nil {
		log.Fatalf("server run failed: %v", err)
		return
	}
	log.Printf("running server on :8080")
}
