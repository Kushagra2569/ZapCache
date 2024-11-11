package main

import "sync"

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
	"HSET": hset,
	"HGET": hget,
}

var SETValues = map[string]string{}
var SETMut = sync.RWMutex{}

var HSETValues = map[string]map[string]string{}
var HSETMut = sync.RWMutex{}

func ping(args []Value) Value {
	return Value{typ: TYP_STRING, str: "PONG"}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: TYP_ERROR, str: "Error: Wrong number of commands"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETMut.Lock()
	SETValues[key] = value
	SETMut.Unlock()

	return Value{typ: TYP_STRING, str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: TYP_ERROR, str: "Error: Wrong number of commands"}
	}

	key := args[0].bulk

	SETMut.RLock()
	value, ok := SETValues[key]
	SETMut.RUnlock()

	if !ok {
		return Value{typ: TYP_NULL}
	}

	return Value{typ: TYP_BULK, bulk: value}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: TYP_ERROR, str: "Error: Wrong number of commands"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	SETMut.Lock()
	if _, ok := HSETValues[hash]; !ok {
		HSETValues[hash] = map[string]string{}
	}
	HSETValues[hash][key] = value
	SETMut.Unlock()

	return Value{typ: TYP_STRING, str: "OK"}
}
func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: TYP_ERROR, str: "Error: Wrong number of commands"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	SETMut.RLock()
	value, ok := HSETValues[hash][key]
	SETMut.RUnlock()

	if !ok {
		return Value{typ: TYP_NULL}
	}

	return Value{typ: TYP_BULK, bulk: value}
}
