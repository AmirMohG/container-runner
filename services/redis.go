package services

type RedisClient interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}
