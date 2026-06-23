package services

import (
	"log"
	"rate_limiter/cache"
	"rate_limiter/config"
)

type ServiceList struct {
	Config       *config.Config
	CacheService cache.Cache
}

func InitServices(config *config.Config) *ServiceList {
	serviceList := ServiceList{
		Config: config,
		CacheService: cache.NewMemoryCache(cache.CacheConfig{
			MaxToken:   config.MaxToken,
			RefillRate: config.RefillRate,
		}),
	}
	if serviceList.CacheService == nil {
		log.Fatal("Failed to create cache service")
		return nil
	}
	return &serviceList
}
