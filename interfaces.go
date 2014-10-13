package memcached_api

type incrDecrHandler func(delta int64) (int64, error)

type statHandler func() (map[string]interface{}, error)
