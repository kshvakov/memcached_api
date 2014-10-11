package memcached_api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strings"
)

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

		responseMessage, err := json.Marshal(response)

		if err != nil {

			responseMessage, _ = json.Marshal(map[string]string{"error": err.Error()})
		}

		connect.Write([]byte(fmt.Sprintf("VALUE %s 0 %d\r\n", method, len(responseMessage))))
		connect.Write(responseMessage)
		connect.Write([]byte("\r\n"))
	}

	connect.Write([]byte("END\r\n"))
}