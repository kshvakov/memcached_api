package memcached_api

import (
	"fmt"
)

func (api *Api) callIncrementDecrement(command string, line []byte, connect netConnector) {

	var (
		method string
		delta  int64
	)

	if command == "incr" {

		api.cmdStat["incr"]++

		fmt.Sscanf(string(line), "incr %s %d", &method, &delta)

	} else {

		api.cmdStat["decr"]++

		fmt.Sscanf(string(line), "decr %s %d", &method, &delta)
	}

	if handler, found := api.incrDecrHandlers[method]; found {

		api.handlerStats[method]++

		if value, err := handler(delta); err == nil {

			connect.Write([]byte(fmt.Sprintf("%d\r\n", value)))

		} else {

			connect.Write([]byte(fmt.Sprintf("SERVER_ERROR %s", err.Error())))
		}

	} else {

		connect.Write([]byte("NOT_FOUND\r\n"))
	}
}
