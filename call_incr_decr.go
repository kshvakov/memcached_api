package memcached_api

import (
	"fmt"
	"net"
)

func (api *Api) callIncrementDecrement(command string, line []byte, connect net.Conn) {

	var (
		method string
		delta  int64
	)

	if command == "incr" {

		fmt.Sscanf(string(line), "incr %s %d", &method, &delta)

	} else {

		fmt.Sscanf(string(line), "decr %s %d", &method, &delta)
	}

	if handler, found := api.incrDecrHandlers[method]; found {

		if value, err := handler(delta); err == nil {

			connect.Write([]byte(fmt.Sprintf("%d\r\n", value)))

		} else {

			connect.Write([]byte(fmt.Sprintf("SERVER_ERROR %s", err.Error())))
		}

	} else {

		connect.Write([]byte("NOT_FOUND\r\n"))
	}
}
