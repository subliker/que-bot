package dispatcher

type QueueConfig struct {
	CacheSize int `yaml:"cache_size" env:"CACHE_SIZE" env-default:"1024" env-description:"CacheSize is how many queue can be stored in cache"`
	CacheTTL  int `yaml:"cache_ttl" env:"CACHE_TTL" env-default:"86400" env-description:"CacheTTL is how long(in seconds) queue can be stored in cache"`
}
