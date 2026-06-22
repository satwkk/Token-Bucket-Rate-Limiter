package middleware

import (
	"log"
	"net/http"
	"rate_limiter/services"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(serviceList *services.ServiceList) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.GetHeader("X-Rate-Limit-ID")
		if ip == "" {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		cacheSvc := serviceList.CacheService
		accepted, tokensBefore, tokensAfter := cacheSvc.TryConsume(ip)

		log.Printf("[BEFORE]: Tokens for %s: %f", ip, tokensBefore)
		log.Printf("[AFTER]: Tokens for %s: %f", ip, tokensAfter)

		if !accepted {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			return
		}

		ctx.Next()
	}
}
