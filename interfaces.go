package memcached_api

type getHandler interface{}

type setHandler func(params interface{}) ([]byte, error)

type deleteHandler func(params interface{}) ([]byte, error)

type incrementHandler func(value int) ([]byte, error)

type decrementHandler func(value int) ([]byte, error)

type statHandler func() (map[string]interface{}, error)
