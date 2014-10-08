Memcached API (Golang implementation)
=============


* Memcached ASCII protocol
* Support multiget request
* Get request with any params (key = method:base64(json(params)))
* JSON response


@todo: finalize it )

simple

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

func (users *Users) GetUserByTwoParams(login string, userId int) (interface{}, error) {

	return &User{UserId: userId, UserName: "Test User", UserLogin: login}, nil
}

func (users *Users) Cast(intParam int, floatParam float64, stringParam string) (interface{}, error) {

	return map[string]interface{}{"Int": intParam, "Float": floatParam, "String": stringParam}, nil
}

func main() {

	users := &Users{}

	api := memcached_api.New()

	api.Get("GetUserById", users.GetUserById)
	api.Get("GetUserByTwoParams", users.GetUserByTwoParams)
	api.Get("Cast", users.Cast)

	api.Run()
}

```