package memcached_api

import (
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

type deleteHandler struct {
	method reflect.Value
	typeIn []reflect.Type
}
