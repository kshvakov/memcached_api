Memcached API (Golang implementation)
=============


* Memcached ASCII protocol
* Get request with any params (key = method:base64(json(params)))
* Support get/gets/set/incr/decr/delete/stats commands
* JSON response for "get/gets" commands


**Example**

simple api server

```go

type User struct {
	UserId    int    `json:"user_id,omitempty"`
	UserName  string `json:"user_name"`
	UserLogin string `json:"user_login,omitempty"`
	UserToken string `json:"user_token,omitempty"`
}

type Users struct {
}

func (users *Users) GetUserById(userId int) (interface{}, error) {

	return &User{UserId: userId, UserName: "Test User"}, nil
}

func (users *Users) SetUser(user *User) error {

	fmt.Println(user)

	return nil
}

func main() {

	users := &Users{}

	api := memcached_api.New("127.0.0.1:3000")

	api.Get("GetUserById", users.GetUserById)
	api.Set("SetUser", users.SetUser)

	api.Stats(func() (map[string]uint, error) {

		return map[string]uint{"Total_User": 42, "Total...": 100000}, nil
	})

	api.Run()
}

```

Client Go

```go
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

func (api *MemcachedApi) SetUser() error {

	value, _ := json.Marshal(&User{UserId: 42, UserName: "New User", UserLogin: "new_login"})

	return api.memcache.Set(&memcache.Item{Key: "SetUser", Value: value})
}

func getCommand(method string, params ...interface{}) string {

	jsonParams, _ := json.Marshal(params)

	return fmt.Sprintf("%s:%s", method, base64.StdEncoding.EncodeToString(jsonParams))
}

func main() {

	api := NewMemcachedApi()

	item, _ := api.GetUserById(42)

	fmt.Println(string(item.Value))
}
```

Client PHP

```php
class MemcachedApi
{
	protected $_memcache;

	public function __construct()
	{
		$this->_memcache = new \Memcache();
		$this->_memcache->connect('127.0.0.1', 3000);
	}

	public function getUserById($userId)
	{
		return $this->_memcache->get($this->_getCommand("GetUserById", (int) $userId));
	}
	
	protected function _getCommand($method, ...$params)
	{
		return sprintf("%s:%s", $method, base64_encode(json_encode($params)));
	}

	public function setUser() 
	{
		return $this->_memcache->set(
			'SetUser', json_encode(
				[
					'user_id'    => 42, 
					'user_name'  => 'New User', 
					'user_login' => 'new_login'
				]
			)
		);
	}

	public function getStats()
	{
		return $this->_memcache->getStats();
	}
}

$Api = new MemcachedApi;

var_dump($Api->getUserById(42));

```


STAT

```
echo "stats" | nc 127.0.0.1 3000
STAT cmd_get 180
STAT cmd_set 20
STAT cmd_incr 100
STAT cmd_decr 100
STAT cmd_delete 20
STAT handler_GetAuthUser 40
STAT handler_SetUser 20
STAT handler_Increment 100
STAT handler_Delete 20
STAT handler_GetUserById 40
STAT handler_Cast 20
STAT handler_ReturnError 20
STAT handler_GetUserWhereIdIn 20
STAT handler_Decrement 100
STAT handler_GetUserByTwoParams 20
STAT Total_User 42
STAT Total... 100000
END
```