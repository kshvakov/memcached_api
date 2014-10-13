package memcached_api

func New() *Api {

	return &Api{
		getHandlers:      make(map[string]*getHandler),
		setHandlers:      make(map[string]*setHandler),
		incrDecrHandlers: make(map[string]incrDecrHandler),
	}
}
