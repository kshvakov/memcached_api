package memcached_api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type connect struct {
}

func (c *connect) Read(p []byte) (n int, err error) {

	return len(p), nil
}

func (c *connect) Write(p []byte) (n int, err error) {

	return len(p), nil
}

func (c *connect) Close() error {

	return nil
}

func BenchmarkGet(b *testing.B) {

	b.StopTimer()

	handler := func(paramInt int, paramsString string) (interface{}, error) {

		return map[string]int{paramsString: paramInt}, nil
	}

	method := reflect.ValueOf(handler)

	typeIn := make([]reflect.Type, method.Type().NumIn())

	for i := 0; i < method.Type().NumIn(); i++ {

		typeIn[i] = method.Type().In(i)
	}

	api := &Api{
		getHandlers: map[string]*getHandler{"Handler": &getHandler{
			method: method,
			typeIn: typeIn,
		}},
		cmdStat:      make(map[string]uint),
		handlerStats: make(map[string]uint),
	}

	connect := &connect{}

	line := []byte(fmt.Sprintf("get %s", command("Handler", 42, "string")))

	b.StartTimer()

	for i := 0; i < b.N; i++ {

		api.callGet(line, connect)
	}
}

func command(method string, params ...interface{}) string {

	jsonParams, _ := json.Marshal(params)

	return fmt.Sprintf("%s:%s", method, base64.StdEncoding.EncodeToString(jsonParams))
}
