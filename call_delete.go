package memcached_api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strings"
)

func (api *Api) callDelete(line []byte, connect net.Conn) {

	part := strings.SplitN(strings.Split(string(line), " ")[1], ":", 2)

	method := part[0]

	if handler, found := api.deleteHandlers[method]; found {

		data, _ := base64.StdEncoding.DecodeString(part[1])

		var tmp []interface{}

		if err := json.Unmarshal(data, &tmp); err == nil {

			params := make([]reflect.Value, len(tmp))

			for i, _ := range tmp {

				params[i] = reflect.ValueOf(tmp[i]).Convert(handler.typeIn[i])
			}

			result := handler.method.Call(params)

			if result[0].IsNil() {

				connect.Write([]byte("DELETED\r\n"))

			} else {

				connect.Write([]byte(fmt.Sprintf("SERVER_ERROR %s\r\n", result[0].Interface())))
			}

		} else {

			connect.Write([]byte(fmt.Sprintf("CLIENT_ERROR %s", err.Error())))
		}

	} else {

		connect.Write([]byte("ERROR\r\n"))
	}
}
