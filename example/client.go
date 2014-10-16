package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
)

type User struct {
	UserId    int    `json:"user_id,omitempty"`
	UserName  string `json:"user_name"`
	UserLogin string `json:"user_login,omitempty"`
	UserToken string `json:"user_token,omitempty"`
}

func NewMemcachedApi() *MemcachedApi {

	return &MemcachedApi{
		memcache: memcache.New("127.0.0.1:3000"),
	}
}

type MemcachedApi struct {
	memcache *memcache.Client
}

func (api *MemcachedApi) GetUserById(userId int) (*memcache.Item, error) {

	return api.memcache.Get(command("GetUserById", userId))
}

func (api *MemcachedApi) GetUserByTwoParams(login string, userId int) (*memcache.Item, error) {

	return api.memcache.Get(command("GetUserByTwoParams", login, userId))
}

func (api *MemcachedApi) GetAuthUser(token string) (*memcache.Item, error) {

	return api.memcache.Get(command("GetAuthUser", token))
}

func (api *MemcachedApi) Cast() (*memcache.Item, error) {

	return api.memcache.Get(command("Cast", 42, 3.14159265359, "Hello"))
}

func (api *MemcachedApi) MultiGet() (map[string]*memcache.Item, error) {

	return api.memcache.GetMulti([]string{
		command("GetUserById", 42),
		command("GetAuthUser", "token"),
	},
	)
}

func (api *MemcachedApi) ReturnError() (*memcache.Item, error) {

	return api.memcache.Get(command("ReturnError"))
}

func (api *MemcachedApi) NotFoundmethod() (*memcache.Item, error) {

	return api.memcache.Get(command("notFoundmethod"))
}

func (api *MemcachedApi) GetUserWhereIdIn() (*memcache.Item, error) {

	return api.memcache.Get(command("GetUserWhereIdIn", []int{1, 2, 3, 4}))
}

func (api *MemcachedApi) SetUser() error {

	value, _ := json.Marshal(&User{UserId: 42, UserName: "New User", UserLogin: "new_login"})

	return api.memcache.Set(&memcache.Item{Key: "SetUser", Value: value})
}

func (api *MemcachedApi) Increment(delta uint64) (uint64, error) {

	return api.memcache.Increment("Increment", delta)
}

func (api *MemcachedApi) Decrement(delta uint64) (uint64, error) {

	return api.memcache.Decrement("Decrement", delta)
}

func (api *MemcachedApi) Delete() error {

	return api.memcache.Delete(command("Delete", 42))
}

func command(method string, params ...interface{}) string {

	jsonParams, _ := json.Marshal(params)

	return fmt.Sprintf("%s:%s", method, base64.StdEncoding.EncodeToString(jsonParams))
}

func main() {

	api := NewMemcachedApi()

	item, _ := api.GetUserById(42)

	fmt.Println(string(item.Value))

	item, _ = api.GetUserByTwoParams("Login", 42)

	fmt.Println(string(item.Value))

	item, _ = api.GetAuthUser("token")

	fmt.Println(string(item.Value))

	item, _ = api.Cast()

	fmt.Println(string(item.Value))

	item, _ = api.NotFoundmethod()

	fmt.Println(string(item.Value))

	items, _ := api.MultiGet()

	for k, item := range items {

		fmt.Println(k, string(item.Value))
	}

	item, _ = api.ReturnError()

	fmt.Println(string(item.Value))

	item, _ = api.GetUserWhereIdIn()

	fmt.Println(string(item.Value))

	fmt.Println(api.SetUser())

	for i := 1; i <= 5; i++ {

		newValue, _ := api.Increment(uint64(i))

		fmt.Println(newValue)
	}

	for i := 1; i <= 5; i++ {

		newValue, _ := api.Decrement(uint64(i))

		fmt.Println(newValue)
	}

	fmt.Println(api.Delete())
}
