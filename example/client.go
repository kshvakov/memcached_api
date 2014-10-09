package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
)

func NewMemcachedApi() *MemcachedApi {

	return &MemcachedApi{
		memcache: memcache.New("127.0.0.1:3000"),
	}
}

type MemcachedApi struct {
	memcache *memcache.Client
}

func (api *MemcachedApi) GetUserById(userId int) (*memcache.Item, error) {

	return api.memcache.Get(getCommand("GetUserById", userId))
}

func (api *MemcachedApi) GetUserByTwoParams(login string, userId int) (*memcache.Item, error) {

	return api.memcache.Get(getCommand("GetUserByTwoParams", login, userId))
}

func (api *MemcachedApi) MultiGet() (map[string]*memcache.Item, error) {

	return api.memcache.GetMulti([]string{
		getCommand("GetUserById", 42),
		getCommand("GetAuthUser", "token"),
	},
	)
}

func getCommand(method string, params ...interface{}) string {

	jsonParams, _ := json.Marshal(params)

	return fmt.Sprintf("%s:%s", method, base64.StdEncoding.EncodeToString(jsonParams))
}

func main() {

	api := NewMemcachedApi()

	item, _ := api.GetUserById(42)

	fmt.Println(string(item.Value))

	item, _ = api.GetUserByTwoParams("Login", 42)

	fmt.Println(string(item.Value))

	items, _ := api.MultiGet()

	for k, item := range items {

		fmt.Println(k, string(item.Value))
	}
}
