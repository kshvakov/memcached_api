package memcached_api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

func (api *Api) callSet(line []byte, reader *bufio.Reader, connect netConnector) {

	api.cmdStat["set"]++

	data, err := reader.ReadBytes('\n')

	if err != nil {

		connect.Write([]byte("NOT_STORED\r\n"))

		return
	}

	part := bytes.Split(line, []byte(" "))

	method := string(part[1])

	if handler, found := api.setHandlers[method]; found {

		api.handlerStats[method]++

		params := handler.params.Interface()

		if err := json.Unmarshal(data, &params); err == nil {

			result := handler.method.Call([]reflect.Value{reflect.ValueOf(params)})

			if result[0].IsNil() {

				connect.Write([]byte("STORED\r\n"))

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
