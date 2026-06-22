package cache

type Cache interface {
	TryConsume(key string) (accepted bool, tokensBefore, tokensAfter float64)
}
