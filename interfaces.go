package memcached_api

type statHandler func() (map[string]interface{}, error)
