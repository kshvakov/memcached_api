package memcached_api

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"
)

type Api struct {
	statHandler statHandler
	getHandlers map[string]getHandler
}

func (api *Api) Get(key string, handler getHandler) error {

	if isValidGetHandler(handler) {

		api.getHandlers[key] = handler

		return nil
	}

	log.Print("Invalid Get handler")

	return fmt.Errorf("Invalid Get handler")
}

func (api *Api) Set(key string, handler setHandler) {

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

	if listener, err := net.Listen("tcp", "127.0.0.1:3000"); err == nil {

		for {

			if connect, err := listener.Accept(); err == nil {

				go api.Handle(connect)

			} else {

				fmt.Println(err.Error())
			}
		}

	} else {

		fmt.Println(err.Error())
	}
}

func (api *Api) Handle(connect net.Conn) {

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

			fmt.Println("set: ", string(request))

			if data, err := reader.ReadBytes('\n'); err == nil {

				log.Printf("data: %s", string(data))

			} else {

				log.Printf("data err: %s", err.Error())
			}

			connect.Write([]byte("STORED\r\n"))

		case bytes.HasPrefix(request, []byte("stats")):

			log.Print("STAT")

			connect.Write([]byte("STAT hh 42\r\n"))
			connect.Write([]byte("END\r\n"))
		}

		log.Print(string(request))

		log.Print("--command--")
	}
}

func (api *Api) callGet(request []byte, connect net.Conn) {

	var response interface{}

	commands := bytes.Split(request, []byte(" "))

	for _, command := range commands[1:] {

		part := strings.SplitN(string(command), ":", 2)

		method := part[0]

		if handler, found := api.getHandlers[method]; found {

			data, _ := base64.StdEncoding.DecodeString(part[1])

			var tmp []interface{}

			if err := json.Unmarshal(data, &tmp); err == nil {

				reflectHandler := reflect.ValueOf(handler)

				params := make([]reflect.Value, len(tmp))

				for i, _ := range tmp {

					params[i] = reflect.ValueOf(tmp[i]).Convert(reflectHandler.Type().In(i))
				}

				result := reflectHandler.Call(params)

				if result[1].IsNil() {

					response = result[0].Interface()

				} else {

					response = map[string]string{"error": fmt.Sprint(result[1].Interface())}
				}

			} else {

				response = map[string]string{"error": fmt.Sprintf("Invalid params (%s)", err.Error())}
			}

		} else {

			response = map[string]string{"error": "Method not found"}
		}

		if responseMessage, err := json.Marshal(response); err == nil {

			connect.Write([]byte(fmt.Sprintf("VALUE %s 0 %d\r\n", method, len(responseMessage))))
			connect.Write(responseMessage)
			connect.Write([]byte("\r\n"))

		} else {

			errorMessage, _ := json.Marshal(map[string]string{"error": err.Error()})

			connect.Write([]byte(fmt.Sprintf("VALUE %s 0 %d\r\n", method, len(errorMessage))))
			connect.Write(errorMessage)
			connect.Write([]byte("\r\n"))
		}
	}

	connect.Write([]byte("END\r\n"))
}

func isValidGetHandler(handler getHandler) bool {

	methodType := reflect.ValueOf(handler).Type()

	if methodType.NumOut() != 2 {

		return false
	}

	return methodType.Out(0).Kind() == reflect.Interface && methodType.Out(1).String() == "error"
}
