package memcached_api

import (
	"bytes"
	"encoding/json"
	"net"
	"reflect"
)

func (api *Api) callSet(request []byte, data []byte, connect net.Conn) {

	split := bytes.Split(request, []byte(" "))

	method := string(split[1])

	if handler, found := api.setHandlers[method]; found {

		params := handler.params.Interface()

		if err := json.Unmarshal(data, &params); err == nil {

			handler.method.Call([]reflect.Value{reflect.ValueOf(params)})

			connect.Write([]byte("STORED\r\n"))

		} else {

			connect.Write([]byte("NOT_STORED\r\n"))
		}

	} else {

		connect.Write([]byte("NOT_STORED\r\n"))
	}
}
