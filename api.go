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

type Api struct {
	statHandler statHandler
	getHandlers map[string]interface{}
	setHandlers map[string]*setHandler
}

func (api *Api) Get(key string, handler interface{}) error {

	if isValidGetHandler(handler) {

		api.getHandlers[key] = handler

		return nil
	}

	log.Print("Invalid Get handler")

	return fmt.Errorf("Invalid Get handler")
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

func (api *Api) Delete(key string, handler deleteHandler) {

}

func (api *Api) Increment(key string, handler incrementHandler) {

}

func (api *Api) Decrement(key string, handler decrementHandler) {

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

		request, err := reader.ReadBytes('\n')

		if err != nil {

			connect.Close()

			log.Print("Close connect")

			break
		}

		switch true {

		case bytes.HasPrefix(request, []byte("get")):

			api.callGet(request, connect)

		case bytes.HasPrefix(request, []byte("set")):

			if data, err := reader.ReadBytes('\n'); err == nil {

				api.callSet(request, data, connect)

			} else {

				log.Printf("data err: %s", err.Error())

				connect.Write([]byte("NOT_STORED\r\n"))
			}

		case bytes.HasPrefix(request, []byte("stats")):

			log.Print("STAT")

			connect.Write([]byte("STAT hh 42\r\n"))
			connect.Write([]byte("END\r\n"))
		}

		log.Print(string(request))

		log.Print("--command--")
	}
}

func isValidGetHandler(handler interface{}) bool {

	methodType := reflect.ValueOf(handler).Type()

	if methodType.NumOut() != 2 {

		return false
	}

	return methodType.Out(0).Kind() == reflect.Interface && methodType.Out(1).String() == "error"
}
