package memcached_api

func New(address string) *Api {

	return &Api{
		address:          address,
		getHandlers:      make(map[string]*getHandler),
		setHandlers:      make(map[string]*setHandler),
		deleteHandlers:   make(map[string]*deleteHandler),
		incrDecrHandlers: make(map[string]func(delta int64) (int64, error)),
	}
}
