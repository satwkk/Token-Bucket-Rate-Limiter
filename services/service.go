package services

import (
	"log"
	"rate_limiter/cache"
)

type ServiceList struct {
	CacheService cache.Cache
}

func InitServices() *ServiceList {
	serviceList := ServiceList{
		CacheService: cache.NewMemoryCache(),
	}
	if serviceList.CacheService == nil {
		log.Fatal("Failed to create cache service")
		return nil
	}
	return &serviceList
}
