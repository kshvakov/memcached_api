package main

import (
	"fmt"
	"github.com/kshvakov/memcached_api"
)

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

func (users *Users) GetAuthUser(token string) (interface{}, error) {

	return &User{UserName: "Test User", UserToken: token}, nil
}

func (users *Users) Cast(intParam int, floatParam float64, stringParam string) (interface{}, error) {

	return map[string]interface{}{"Int": intParam, "Float": floatParam, "String": stringParam}, nil
}

func (users *Users) ReturnError() (interface{}, error) {

	return nil, fmt.Errorf("Error message")
}

func (users *Users) GetUserWhereIdIn(userIds []interface{}) (interface{}, error) {

	var result []*User

	for _, userId := range userIds {

		result = append(result, &User{UserId: int(userId.(float64)), UserName: "Test User"})
	}

	return result, nil
}

func (users *Users) SetUser(user *User) error {

	fmt.Println(user)

	return nil
}

func main() {

	users := &Users{}

	api := memcached_api.New("127.0.0.1:3000")

	api.Get("GetUserById", users.GetUserById)
	api.Get("GetAuthUser", users.GetAuthUser)
	api.Get("GetUserByTwoParams", users.GetUserByTwoParams)
	api.Get("Cast", users.Cast)
	api.Get("ReturnError", users.ReturnError)
	api.Get("GetUserWhereIdIn", users.GetUserWhereIdIn)
	api.Set("SetUser", users.SetUser)

	api.Increment("Increment", func(delta int64) (int64, error) {

		fmt.Printf("delta: %d\n", delta)

		return delta + 42, nil
	})

	api.Decrement("Decrement", func(delta int64) (int64, error) {

		fmt.Printf("delta: %d\n", delta)

		return 42 - delta, nil
	})

	api.Run()
}
