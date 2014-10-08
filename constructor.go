package memcached_api

func New() *Api {

	return &Api{
		getHandlers: make(map[string]getHandler),
	}
}
