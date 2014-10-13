package memcached_api

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"reflect"
)

type setHandler struct {
	method reflect.Value
	params reflect.Value
}

type getHandler struct {
	method reflect.Value
	typeIn []reflect.Type
}

type Api struct {
	address          string
	statHandler      statHandler
	getHandlers      map[string]*getHandler
	setHandlers      map[string]*setHandler
	incrDecrHandlers map[string]func(delta int64) (int64, error)
}

func (api *Api) Get(key string, handler interface{}) {

	method := reflect.ValueOf(handler)

	if method.Type().NumOut() != 2 {

		log.Fatal("Invalid Get handler")
	}

	if method.Type().Out(0).Kind() != reflect.Interface || method.Type().Out(1).String() != "error" {

		log.Fatal("Invalid Get handler")
	}

	typeIn := make([]reflect.Type, method.Type().NumIn())

	for i := 0; i < method.Type().NumIn(); i++ {

		typeIn[i] = method.Type().In(i)
	}

	api.getHandlers[key] = &getHandler{
		method: method,
		typeIn: typeIn,
	}
}

func (api *Api) Set(key string, handler interface{}) {

	method := reflect.ValueOf(handler)

	if method.Type().NumIn() == 1 && method.Type().In(0).Kind() == reflect.Ptr {

		api.setHandlers[key] = &setHandler{
			method: method,
			params: reflect.New(method.Type().In(0).Elem()),
		}

		return
	}

	log.Fatal("Invalid Set handler")
}

/*
func (api *Api) Delete(key string, handler deleteHandler) {

}
*/
func (api *Api) Increment(key string, handler func(delta int64) (int64, error)) {

	if _, found := api.incrDecrHandlers[key]; found {

		log.Fatal("handler '%s' is already registered.", key)
	}

	api.incrDecrHandlers[key] = handler

}

func (api *Api) Decrement(key string, handler func(delta int64) (int64, error)) {

	if _, found := api.incrDecrHandlers[key]; found {

		log.Fatal("handler '%s' is already registered.", key)
	}

	api.incrDecrHandlers[key] = handler
}

func (api *Api) Stats(handler statHandler) {

	api.statHandler = handler
}

func (api *Api) Run() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	if listener, err := net.Listen("tcp", api.address); err == nil {

		for {

			if connect, err := listener.Accept(); err == nil {

				go api.handle(connect)

			} else {

				log.Print(err.Error())
			}
		}

	} else {

		log.Fatal(err.Error())
	}
}

func (api *Api) handle(connect net.Conn) {

	defer func() {

		if message := recover(); message != nil {

			connect.Write([]byte(fmt.Sprintf("SERVER_ERROR %s\r\nEND\r\n", message)))

			log.Printf("Server error: %s", message)
		}
	}()

	reader := bufio.NewReader(connect)

	for {

		line, err := reader.ReadBytes('\n')

		if err != nil {

			connect.Close()

			log.Print("Close connect")

			break
		}

		switch true {

		case bytes.HasPrefix(line, []byte("get")):

			api.callGet(line, connect)

		case bytes.HasPrefix(line, []byte("set")):

			api.callSet(line, reader, connect)

		case bytes.HasPrefix(line, []byte("incr")):

			api.callIncrementDecrement("incr", line, connect)

		case bytes.HasPrefix(line, []byte("decr")):

			api.callIncrementDecrement("decr", line, connect)

		case bytes.HasPrefix(line, []byte("stats")):

			log.Print("STAT")

			connect.Write([]byte("STAT hh 42\r\n"))
			connect.Write([]byte("END\r\n"))
		}
	}
}
