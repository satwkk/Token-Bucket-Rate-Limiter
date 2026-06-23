package main

import (
	"log"
	"net/http/httputil"
	"net/url"
	"rate_limiter/config"
	"rate_limiter/middleware"
	"rate_limiter/services"

	"github.com/gin-gonic/gin"
)

func main() {
	if ok := config.LoadEnv(); !ok {
		log.Fatal("Error loading .env file")
		return
	}

	proxyConfig := config.LoadProxyConfig()
	if proxyConfig == nil {
		log.Fatal("Error loading proxy config")
		return
	}

	backendUrl, err := url.Parse(proxyConfig.BackendUrl)
	if err != nil {
		log.Fatalf("failed to parse backend URL: %v", err)
		return
	}

	serviceList := services.InitServices(proxyConfig)
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
