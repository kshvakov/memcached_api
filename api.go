package memcached_api

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"reflect"
	"time"
)

type Api struct {
	address          string
	statHandler      informer
	getHandlers      map[string]*getHandler
	setHandlers      map[string]*setHandler
	deleteHandlers   map[string]*deleteHandler
	incrDecrHandlers map[string]incrDecr
	cmdStat          map[string]uint
	handlerStats     map[string]uint
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

func (api *Api) Delete(key string, handler interface{}) {

	method := reflect.ValueOf(handler)

	if method.Type().NumOut() != 1 || method.Type().Out(0).String() != "error" {

		log.Fatal("Invalid Delete handler")
	}

	typeIn := make([]reflect.Type, method.Type().NumIn())

	for i := 0; i < method.Type().NumIn(); i++ {

		typeIn[i] = method.Type().In(i)
	}

	api.deleteHandlers[key] = &deleteHandler{
		method: method,
		typeIn: typeIn,
	}
}

func (api *Api) Increment(key string, handler incrDecr) {

	if _, found := api.incrDecrHandlers[key]; found {

		log.Fatal("handler '%s' is already registered.", key)
	}

	api.incrDecrHandlers[key] = handler

}

func (api *Api) Decrement(key string, handler incrDecr) {

	if _, found := api.incrDecrHandlers[key]; found {

		log.Fatal("handler '%s' is already registered.", key)
	}

	api.incrDecrHandlers[key] = handler
}

func (api *Api) Stats(handler informer) {

	api.statHandler = handler
}

func (api *Api) Run() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	if listener, err := net.Listen("tcp", api.address); err == nil {

		for {

			if connect, err := listener.Accept(); err == nil {

				connect.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))

				go api.handle(connect)

			} else {

				log.Print(err.Error())
			}
		}

	} else {

		log.Fatal(err.Error())
	}
}

func (api *Api) handle(connect netConnector) {

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

			//	log.Print("Close connect")

			break
		}

		switch true {

		case bytes.HasPrefix(line, []byte("get")):

			api.callGet(line, connect)

		case bytes.HasPrefix(line, []byte("set")):

			api.callSet(line, reader, connect)

		case bytes.HasPrefix(line, []byte("delete")):

			api.callDelete(line, connect)

		case bytes.HasPrefix(line, []byte("incr")):

			api.callIncrementDecrement("incr", line, connect)

		case bytes.HasPrefix(line, []byte("decr")):

			api.callIncrementDecrement("decr", line, connect)

		case bytes.HasPrefix(line, []byte("stats")):

			for cmd, reqs := range api.cmdStat {

				connect.Write([]byte(fmt.Sprintf("STAT cmd_%s %d\r\n", cmd, reqs)))
			}

			for handler, reqs := range api.handlerStats {

				connect.Write([]byte(fmt.Sprintf("STAT handler_%s %d\r\n", handler, reqs)))
			}

			if stat, err := api.statHandler(); err == nil {

				for k, v := range stat {

					connect.Write([]byte(fmt.Sprintf("STAT %s %d\r\n", k, v)))
				}
			}

			connect.Write([]byte("END\r\n"))
		}
	}
}
