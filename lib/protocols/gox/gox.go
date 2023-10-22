package gox

import (
	"fmt"
	"reflect"

	"slack/lib/proto"
)

var funcMap = map[string]any{}

func Request(target, data string, variableMap map[string]any) error {
	err := callFunction(data, []any{target, variableMap}, funcMap)
	if err != nil {
		return err.(error)
	}
	return nil
}

func callFunction(name string, args []interface{}, funcMap map[string]interface{}) interface{} {
	f, ok := funcMap[name]
	if !ok {
		fmt.Printf("function %s not found\n", name)
		return nil
	}

	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		fmt.Printf("%s is not a function\n", name)
		return nil
	}
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}
	out := v.Call(in)
	if len(out) == 0 {
		return nil
	}
	return out[0].Interface()
}

func setRequest(data string, vmap map[string]any) {
	vmap["request"] = &proto.Request{
		Raw: []byte(data),
	}
}

func setResponse(data string, vmap map[string]any) {
	vmap["response"] = &proto.Response{
		Raw: []byte(data),
	}
}

func setFullTarget(data string, vmap map[string]any) {
	vmap["fulltarget"] = data
}

func setTarget(data string, vmap map[string]any) {
	vmap["target"] = data
}
