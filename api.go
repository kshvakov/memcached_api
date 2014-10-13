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
	statHandler      statHandler
	getHandlers      map[string]*getHandler
	setHandlers      map[string]*setHandler
	incrDecrHandlers map[string]incrDecrHandler
}

func (api *Api) Get(key string, handler interface{}) error {

	method := reflect.ValueOf(handler)

	if method.Type().NumOut() != 2 || method.Type().Out(0).Kind() != reflect.Interface || method.Type().Out(1).String() != "error" {

		log.Print("Invalid Get handler")

		return fmt.Errorf("Invalid Get handler")
	}

	typeIn := make([]reflect.Type, method.Type().NumIn())

	for i := 0; i < method.Type().NumIn(); i++ {

		typeIn[i] = method.Type().In(i)
	}

	api.getHandlers[key] = &getHandler{
		method: method,
		typeIn: typeIn,
	}

	return nil
}

func (api *Api) Set(key string, handler interface{}) error {

	method := reflect.ValueOf(handler)

	if method.Type().NumIn() == 1 && method.Type().In(0).Kind() == reflect.Ptr {

		api.setHandlers[key] = &setHandler{
			method: method,
			params: reflect.New(method.Type().In(0).Elem()),
		}

		return nil
	}

	log.Print("Invalid Set handler")

	return fmt.Errorf("Invalid Set handler")
}

/*
func (api *Api) Delete(key string, handler deleteHandler) {

}
*/
func (api *Api) Increment(key string, handler incrDecrHandler) error {

	if _, found := api.incrDecrHandlers[key]; found {

		log.Printf("handler '%s' is already registered.", key)

		return fmt.Errorf("handler '%s' is already registered.", key)
	}

	api.incrDecrHandlers[key] = handler

	return nil
}

func (api *Api) Decrement(key string, handler incrDecrHandler) error {

	if _, found := api.incrDecrHandlers[key]; found {

		log.Printf("handler '%s' is already registered.", key)

		return fmt.Errorf("handler '%s' is already registered.", key)
	}

	api.incrDecrHandlers[key] = handler

	return nil
}

func (api *Api) Stats(handler statHandler) {

	api.statHandler = handler
}

func (api *Api) Run() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	if listener, err := net.Listen("tcp", "127.0.0.1:3000"); err == nil {

		for {

			if connect, err := listener.Accept(); err == nil {

				go api.handle(connect)

			} else {

				fmt.Println(err.Error())
			}
		}

	} else {

		fmt.Println(err.Error())
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

		log.Print(string(line))

		log.Print("--command--")
	}
}
