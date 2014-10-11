package memcached_api

func New() *Api {

	return &Api{
		getHandlers: make(map[string]interface{}),
		setHandlers: make(map[string]*setHandler),
	}
}
